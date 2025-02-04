package main

/*
Package main implements a TCP server designed for distributed image processing using a worker pool architecture.
The server supports concurrent image processing tasks such as grayscale transformation, edge detection, and geometry computation.

---

### Features
1. **TCP Communication**:
   - Handles incoming connections from clients.
   - Receives image data over TCP.
   - Sends the processed image back to the client.

2. **Worker Pool**:
   - Utilizes a worker pool to process tasks concurrently.
   - Supports tasks like grayscale image transformation, edge detection, and contour finding.

3. **Image Processing Pipeline**:
   - Processes images in chunks for efficient parallelism.
   - Tasks include:
     - Grayscale conversion.
     - Canny edge detection.
     - Contour finding and quadrilateral detection.

---

### Constants
- `host` (string): Host address for the server (default: localhost).
- `port` (string): Port for the server (default: 14750).
- `protocol` (string): Communication protocol (default: TCP).
- `bufferSize` (int): Size of the buffer used for TCP communication.
- `overlapSize` (int): Overlap size between chunks of image processing.
- `numWorkers` (int): Number of workers in the worker pool (defaults to the number of CPU cores).

---

### Structures
#### `workerChannels`
Represents the channels used for communication between tasks and workers.
- Fields:
  - `socketSemaphore`: Used to limit simultaneous socket connections.
  - `imageChan`: Tasks for image transformation (e.g., grayscale, edge detection).
  - `bfsChan`: Tasks for finding contours using BFS.
  - `findQuadrilateralChan`: Tasks for detecting quadrilaterals from contours.

#### `Server`
Represents the TCP server.
- Fields:
  - `host`: Host address for the server.
  - `port`: Port for the server.
  - `stopCtx`: Context to signal server shutdown.
  - `cancel`: Callback function to trigger the context cancellation.
  - `numWorkers`: Number of concurrent workers.

- Methods:
  - `listen()`: Starts listening on the specified host and port.
  - `receiveImage(conn net.Conn)`: Receives and decodes an image from the connection.
  - `sendImage(conn net.Conn, img image.Image, format string)`: Encodes and sends an image to the client.
  - `handleConnection(conn net.Conn, workerChannels workerChannels)`: Manages the entire image processing pipeline for a TCP connection.
  - `run()`: Main loop for accepting and managing connections.
  - `newServer(host string, port string, numWorkers int) *Server`: Initializes a new server instance.

---

### Workflow
1. **Connection Handling**:
   - Begins by listening on the specified `host` and `port`.
   - Accepts incoming TCP connections.
   - Receives the image data from the client using `receiveImage`.

2. **Image Processing**:
   - Splits the image into chunks for parallel processing by workers.
   - Chunks are processed in stages:
     - Grayscale transformation.
     - Canny edge detection.
     - Contour and quadrilateral detection.

3. **Result Aggregation**:
   - Combines processed chunks into the final output image.
   - Sends the final processed image back to the client using `sendImage`.

4. **Worker Pool**:
   - Uses multiple worker pools for different computations (e.g., grayscale conversion, BFS for contours).
   - Tasks are distributed to workers via channels.

5. **Graceful Shutdown**:
   - Listens for an interrupt signal (e.g., CTRL + C).
   - Stops accepting new connections and gracefully shuts down.

---

### Key Image Processing Functions
#### `GrayscaleWrapper(img image.Image) (image.Image, error)`
Converts an image to grayscale using a utility function.

#### `ApplyCannyEdgeDetectionWrapper(img image.Image) (image.Image, error)`
Applies Canny edge detection to a grayscale image.

#### `FindQuadrilateralWrapper(contours []geometry.Contour) (geometry.ContourWithArea, error)`
Finds the largest quadrilateral from a set of contours.

---

### Logging
- Logs server events to `server.log`.
- Important logs include:
  - Server start and shutdown.
  - New connections.
  - Errors during image processing.
  - Task completion and results.

---

### Example Usage:
1. Start the server with:
   ```
   go run main.go
   ```

2. Connect to the server using a TCP client and send an image for processing.

3. Processed image is returned to the client.

---

### Dependencies
- **geometry**: Used for contour and geometric computations.
- **imageUtils**: Provides utility functions for image transformations.
- **utils**: Contains advanced image processing algorithms like edge detection and contour extraction.
- **worker**: Manages task distribution and the worker pool.
*/

import (
	"ELP-project/internal/geometry"
	"ELP-project/internal/imageUtils"
	"ELP-project/internal/utils"
	"ELP-project/internal/worker"
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime"
	"sort"
)

const (
	host        = "localhost"
	port        = "14750"
	protocol    = "tcp"
	bufferSize  = 1024
	overlapSize = 20
)

var numWorkers = runtime.NumCPU()

type workerChannels struct {
	socketSemaphore       chan net.Conn
	imageChan             chan worker.Task[image.Image, image.Image]
	bfsChan               chan worker.Task[image.Rectangle, []geometry.Contour]
	findQuadrilateralChan chan worker.Task[[]geometry.Contour, geometry.ContourWithArea]
}

type Server struct {
	host       string
	port       string
	stopCtx    context.Context
	cancel     context.CancelFunc
	numWorkers int
}

func newServer(host string, port string, numWorkers int) *Server {
	ctx, cancel := context.WithCancel(context.Background())
	return &Server{
		host:       host,
		port:       port,
		stopCtx:    ctx,
		cancel:     cancel,
		numWorkers: numWorkers,
	}
}

func (server *Server) listen() net.Listener {
	listener, err := net.Listen(protocol, fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Printf("Server is listening on IP address %v and port %v...", server.host, server.port)

	return listener
}

func (server *Server) receiveImage(conn net.Conn) (image.Image, string) {
	var dataBuffer bytes.Buffer
	buffer := make([]byte, bufferSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(err.Error())
			if err.Error() == "EOF" || err == io.EOF {
				break
			}
			log.Fatalf("Error reading from connection: %v", err)
		}

		dataBuffer.Write(buffer[:n])

		if bytes.Contains(dataBuffer.Bytes(), []byte("EOF")) {
			log.Println("End of data detected.")
			break
		}
	}
	data := dataBuffer.Bytes()
	data = bytes.TrimSuffix(data, []byte("EOF"))

	img, format, err := image.Decode(&dataBuffer)
	if err != nil {
		log.Fatalf("Error decoding image: %v", err)
	}

	log.Printf("Image decoded successfully. Format: %s", format)
	return img, format
}

func imageToBuffer(img image.Image, format string) (*bytes.Buffer, error) {
	var buffer bytes.Buffer

	switch format {
	case "jpeg":
		err := jpeg.Encode(&buffer, img, nil)
		if err != nil {
			log.Fatalf("failed to encode image to JPEG: %v", err)
		}
	case "png":
		err := png.Encode(&buffer, img)
		if err != nil {
			log.Fatalf("failed to encode image to PNG: %v", err)
		}
	default:
		log.Fatalf("unsupported format: %v", format)
	}

	return &buffer, nil
}

func (server *Server) sendImage(conn net.Conn, img image.Image, format string) {
	buffer, err := imageToBuffer(img, format)
	if err != nil {
		log.Fatalf("Error encoding image: %v", err)
	}

	data := buffer.Bytes()
	dataLen := len(data)
	sent := 0

	for sent < dataLen {
		chunkSize := bufferSize
		if dataLen-sent < bufferSize {
			chunkSize = dataLen - sent
		}

		n, err := conn.Write(data[sent : sent+chunkSize])
		if err != nil {
			log.Fatalf("Error sending data: %v", err)
		}

		sent += n
	}

	log.Printf("Image sent successfully. Total bytes: %d", dataLen)
}

func (server *Server) handleConnection(conn net.Conn, workerChannels workerChannels) {
	defer conn.Close()

	workerChannels.socketSemaphore <- conn
	defer func() { <-workerChannels.socketSemaphore }()

	log.Printf("New connection from %s", conn.RemoteAddr())

	log.Println("Receiving image...")
	img, format := server.receiveImage(conn)
	if img == nil {
		log.Printf("Failed to receive image from %s", conn.RemoteAddr())
		return
	}
	log.Println("Image received successfully!")
	resultGrayChan := make(chan worker.Task[image.Image, image.Image], 100)

	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		bounds := img.Bounds()
		rgbaImg = image.NewRGBA(bounds)
		draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)
	}

	bounds := img.Bounds()
	totalRows := bounds.Max.Y - bounds.Min.Y
	chunkSize := (totalRows + server.numWorkers - 1) / server.numWorkers

	for i := 0; i < server.numWorkers; i++ {
		startY := bounds.Min.Y + i*chunkSize
		endY := startY + chunkSize + overlapSize

		if startY > overlapSize {
			startY -= overlapSize
		}

		if endY > bounds.Max.Y {
			endY = bounds.Max.Y
		}

		subBounds := image.Rect(bounds.Min.X, startY, bounds.Max.X, endY)

		subImage, ok := rgbaImg.SubImage(subBounds).(*image.RGBA)
		if !ok {
			log.Fatalf("SubImage cast failed: expected *image.RGBA")
		}

		task := worker.Task[image.Image, image.Image]{
			Conn:       conn,
			Input:      subImage,
			ResultChan: resultGrayChan,
			Function:   GrayscaleWrapper,
		}
		workerChannels.imageChan <- task
	}

	resultCannyChan := make(chan worker.Task[image.Image, image.Image], 100)

	for i := 0; i < server.numWorkers; i++ {
		select {
		case result := <-resultGrayChan:
			if result.Err != nil {
				log.Printf("Error processing image for %s: %v", conn.RemoteAddr(), result.Err)
				return
			}
			task := worker.Task[image.Image, image.Image]{
				Conn:       conn,
				Input:      result.Output,
				ResultChan: resultCannyChan,
				Function:   ApplyCannyEdgeDetectionWrapper,
			}
			workerChannels.imageChan <- task
		case <-server.stopCtx.Done():
			log.Println("Server is shutting down, closing connection.")
			return
		}
	}
	close(resultGrayChan)

	results := make([]*image.Gray, server.numWorkers)

	for i := 0; i < server.numWorkers; i++ {
		select {
		case result := <-resultCannyChan:
			if result.Err != nil {
				log.Printf("Error processing image for %s: %v", conn.RemoteAddr(), result.Err)
				return
			}
			results[i] = result.Output.(*image.Gray)
		case <-server.stopCtx.Done():
			log.Println("Server is shutting down, closing connection.")
			return
		}
	}
	close(resultCannyChan)

	sort.Slice(results, func(i, j int) bool {
		return results[i].Rect.Min.Y < results[j].Rect.Min.Y
	})

	cannyImage := image.NewGray(bounds)
	for i, chunk := range results {
		startY := bounds.Min.Y + i*chunkSize
		chunkHeight := chunk.Rect.Dy() - overlapSize
		draw.Draw(cannyImage, image.Rect(bounds.Min.X, startY, bounds.Max.X, startY+chunkHeight), chunk, image.Point{X: bounds.Min.X, Y: startY}, draw.Src)
	}

	resultBfsChan := make(chan worker.Task[image.Rectangle, []geometry.Contour], 100)

	FindContoursBFSWrapper := func(rect image.Rectangle) ([]geometry.Contour, error) {
		return utils.FindContoursBFS(cannyImage, rect), nil
	}

	for i := 0; i < server.numWorkers; i++ {
		startY := bounds.Min.Y + i*chunkSize
		endY := startY + chunkSize

		if endY > bounds.Max.Y {
			endY = bounds.Max.Y
		}

		rect := image.Rect(bounds.Min.X, startY, bounds.Max.X, endY)

		task := worker.Task[image.Rectangle, []geometry.Contour]{
			Conn:       conn,
			Input:      rect,
			ResultChan: resultBfsChan,
			Function:   FindContoursBFSWrapper,
		}
		workerChannels.bfsChan <- task
	}

	bfsResult := make([]geometry.Contour, 0)
	for i := 0; i < server.numWorkers; i++ {
		select {
		case result := <-resultBfsChan:
			if result.Err != nil {
				log.Printf("Error processing image for %s: %v", conn.RemoteAddr(), result.Err)
				return
			}
			bfsResult = append(bfsResult, result.Output...)
		case <-server.stopCtx.Done():
			log.Println("Server is shutting down, closing connection.")
			return
		}
	}
	close(resultBfsChan)

	resultFindQuadrilateralChan := make(chan worker.Task[[]geometry.Contour, geometry.ContourWithArea], 100)
	for i := 0; i < numWorkers; i++ {
		start := i * (len(bfsResult) / numWorkers)
		end := (i + 1) * (len(bfsResult) / numWorkers)

		if i == numWorkers-1 {
			end = len(bfsResult)
		}

		task := worker.Task[[]geometry.Contour, geometry.ContourWithArea]{
			Conn:       conn,
			Input:      bfsResult[start:end],
			ResultChan: resultFindQuadrilateralChan,
			Function:   FindQuadrilateralWrapper,
		}
		workerChannels.findQuadrilateralChan <- task
	}

	findQuadrilateralResult := make([]geometry.ContourWithArea, 0)
	for i := 0; i < server.numWorkers; i++ {
		select {
		case result := <-resultFindQuadrilateralChan:
			if result.Err != nil {
				log.Printf("Error processing image for %s: %v", conn.RemoteAddr(), result.Err)
				return
			}
			findQuadrilateralResult = append(findQuadrilateralResult, result.Output)
		case <-server.stopCtx.Done():
			log.Println("Server is shutting down, closing connection.")
			return
		}
	}
	close(resultFindQuadrilateralChan)

	contourA4 := geometry.ContourWithArea{
		Area: 0,
	}
	for _, contour := range findQuadrilateralResult {
		if contour.Area > contourA4.Area {
			contourA4 = contour
		}
	}

	center := geometry.Point{
		X: img.Bounds().Dx() / 2,
		Y: img.Bounds().Dy() / 2,
	}
	contourA4.Contour = utils.FindCorner(contourA4.Contour, center)

	rect := image.Rect(contourA4.Contour[0].X, contourA4.Contour[0].Y, contourA4.Contour[1].X, contourA4.Contour[1].Y)
	finalImage := image.NewRGBA(rect)
	draw.Draw(finalImage, rect, img, image.Pt(contourA4.Contour[0].X, contourA4.Contour[0].Y), draw.Src)

	log.Printf("Sending processed image back to %s", conn.RemoteAddr())
	server.sendImage(conn, finalImage, format)
	log.Println("Connection finished:", conn.RemoteAddr())
}

func FindQuadrilateralWrapper(contours []geometry.Contour) (geometry.ContourWithArea, error) {
	return utils.FindQuadrilateral(contours), nil
}

func ApplyCannyEdgeDetectionWrapper(img image.Image) (image.Image, error) {
	return utils.ApplyCannyEdgeDetection(img.(*image.Gray)), nil
}

func GrayscaleWrapper(img image.Image) (image.Image, error) {
	return imageUtils.Grayscale(img), nil
}

func (server *Server) run() {
	listener := server.listen()
	defer func(listener net.Listener) {
		var opErr *net.OpError
		if err := listener.Close(); err != nil && !(errors.As(err, &opErr) && !opErr.Temporary()) {
			log.Fatalf("Unexpected error closing listener: %v", err)
		}
	}(listener)

	fmt.Println("The server is running... (Press Ctrl + C to stop)")

	socketSemaphore := make(chan net.Conn, 5)
	imageChan := make(chan worker.Task[image.Image, image.Image], 100)
	bfsChan := make(chan worker.Task[image.Rectangle, []geometry.Contour], 100)
	findQuadrilateralChan := make(chan worker.Task[[]geometry.Contour, geometry.ContourWithArea], 100)

	channels := workerChannels{
		socketSemaphore:       socketSemaphore,
		imageChan:             imageChan,
		bfsChan:               bfsChan,
		findQuadrilateralChan: findQuadrilateralChan,
	}

	go worker.StartWorkerPool("Image Worker", numWorkers, worker.TreatmentWorker, imageChan)
	go worker.StartWorkerPool("BFS worker", numWorkers, worker.TreatmentWorker, bfsChan)
	go worker.StartWorkerPool("FindQuadrilateral worker", numWorkers, worker.TreatmentWorker, findQuadrilateralChan)

	go func() {
		<-server.stopCtx.Done()
		log.Println("Shutting down server...")
		err := listener.Close()
		if err != nil {
			log.Printf("Error closing listener: %v", err)
			return
		}
		close(socketSemaphore)
		close(imageChan)
		close(bfsChan)
		close(findQuadrilateralChan)
		log.Println("All workers will stop after completing their tasks.")

	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			var opErr *net.OpError
			if errors.As(err, &opErr) && !opErr.Temporary() {
				log.Println("Listener has been closed. Stopping server gracefully.")
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		select {
		case <-server.stopCtx.Done():
			log.Println("Server is shutting down, closing new connection.")
			conn.Close()
		default:
			go server.handleConnection(conn, channels)
		}
	}
}

func main() {
	logFile, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			log.Fatalf("Error closing log file: %v", err)
		}
	}(logFile)

	log.SetOutput(logFile)

	log.Println("Starting server...")

	server := newServer(host, port, numWorkers)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		log.Println("Interrupt signal received.")
		server.cancel()
	}()

	server.run()
	log.Println("Server shut down gracefully.")
	os.Exit(0)
}
