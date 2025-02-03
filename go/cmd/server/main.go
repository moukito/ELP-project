package main

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
	bfsChan               chan worker.Task[image.Gray, []geometry.Contour]
	findQuadrilateralChan chan worker.Task[[]geometry.Contour, geometry.ContourWithArea]
}

// Server is a struct that encapsulates logic for handling TCP connections.
type Server struct {
	host       string
	port       string
	stopCtx    context.Context
	cancel     context.CancelFunc
	numWorkers int
}

// newServer create a new server instance
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

// listen sets up a server to listen for incoming connections.
func (server *Server) listen() net.Listener {
	listener, err := net.Listen(protocol, fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Printf("Server is listening on IP address %v and port %v...", server.host, server.port)

	return listener
}

// receiveImage handles receiving an image over a connection.
func (server *Server) receiveImage(conn net.Conn) (image.Image, string) {
	// Create a buffer to store the incoming data
	var dataBuffer bytes.Buffer
	buffer := make([]byte, bufferSize) // Temporary buffer size for chunks

	for {
		// Read incoming data into the temporary buffer
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(err.Error())
			if err.Error() == "EOF" || err == io.EOF {
				// End of data, break the loop
				break
			}
			// Handle unexpected errors
			log.Fatalf("Error reading from connection: %v", err)
		}

		// Write the received chunk to the data buffer
		dataBuffer.Write(buffer[:n])

		// Check for delimiter (e.g., "EOF")
		if bytes.Contains(dataBuffer.Bytes(), []byte("EOF")) {
			log.Println("End of data detected.")
			break
		}
	}
	// Remove the delimiter
	data := dataBuffer.Bytes()
	data = bytes.TrimSuffix(data, []byte("EOF"))

	// Decode the accumulated data into an image.Image
	img, format, err := image.Decode(&dataBuffer)
	if err != nil {
		log.Fatalf("Error decoding image: %v", err)
	}

	log.Printf("Image decoded successfully. Format: %s", format)
	return img, format
}

// imageToBuffer converts an image into a byte buffer based on the provided format.
func imageToBuffer(img image.Image, format string) (*bytes.Buffer, error) {
	// Create a new bytes.Buffer
	var buffer bytes.Buffer

	// Encode the image into the buffer
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

	// Return the resulting buffer
	return &buffer, nil
}

// sendImage handles sending an image based on its provided format over a connection.
func (server *Server) sendImage(conn net.Conn, img image.Image, format string) {
	// Use a buffer to encode the image
	buffer, err := imageToBuffer(img, format)
	if err != nil {
		log.Fatalf("Error encoding image: %v", err)
	}

	// Send the data from the buffer in chunks
	data := buffer.Bytes()
	dataLen := len(data)
	sent := 0

	// Send the entire buffer content in chunks
	for sent < dataLen {
		// Determine the number of bytes to send in this chunk
		chunkSize := bufferSize // e.g., 1024 bytes
		if dataLen-sent < bufferSize {
			chunkSize = dataLen - sent
		}

		// Write the chunk to the connection
		n, err := conn.Write(data[sent : sent+chunkSize])
		if err != nil {
			log.Fatalf("Error sending data: %v", err)
		}

		// Advance the cursor position
		sent += n
	}

	log.Printf("Image sent successfully. Total bytes: %d", dataLen)
}

// handleConnection limits the number of simultaneous socket workers and dispatches tasks.
func (server *Server) handleConnection(conn net.Conn, workerChannels workerChannels) {
	defer conn.Close()

	// Limit the number of active socket workers using the semaphore
	workerChannels.socketSemaphore <- conn
	defer func() { <-workerChannels.socketSemaphore }() // Release semaphore when done

	log.Printf("New connection from %s", conn.RemoteAddr())

	// Receive image over the connection
	log.Println("Receiving image...")
	img, format := server.receiveImage(conn)
	if img == nil {
		log.Printf("Failed to receive image from %s", conn.RemoteAddr())
		return
	}
	log.Println("Image received successfully!")
	// Create a dedicated result channel for this task
	resultGrayChan := make(chan worker.Task[image.Image, image.Image], 100)

	// Convert the received image to a concrete type that supports SubImage
	rgbaImg, ok := img.(*image.RGBA)
	if !ok {
		// Convert to RGBA if it's not already in that format
		bounds := img.Bounds()
		rgbaImg = image.NewRGBA(bounds)
		draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)
	}

	// Split the image into row chunks
	bounds := img.Bounds()
	totalRows := bounds.Max.Y - bounds.Min.Y
	chunkSize := (totalRows + server.numWorkers - 1) / server.numWorkers // Rows per worker

	// Dispatch grayscale tasks to treatment workers
	for i := 0; i < server.numWorkers; i++ {
		startY := bounds.Min.Y + i*chunkSize
		endY := startY + chunkSize + overlapSize

		if startY > overlapSize {
			startY -= overlapSize
		}

		if endY > bounds.Max.Y {
			endY = bounds.Max.Y
		}

		// Define the sub-rectangle for this chunk
		subBounds := image.Rect(bounds.Min.X, startY, bounds.Max.X, endY)

		// Define a sub-image for the range of rows
		subImage, ok := rgbaImg.SubImage(subBounds).(*image.RGBA)
		if !ok {
			log.Fatalf("SubImage cast failed: expected *image.RGBA")
		}

		// Create and submit a task for the subimage
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
		case <-server.stopCtx.Done(): // Handle shutdown gracefully
			log.Println("Server is shutting down, closing connection.")
			return
		}
	}
	close(resultGrayChan)

	// Wait for all results and assemble the full image
	results := make([]*image.Gray, server.numWorkers)

	for i := 0; i < server.numWorkers; i++ {
		select {
		case result := <-resultCannyChan:
			if result.Err != nil {
				log.Printf("Error processing image for %s: %v", conn.RemoteAddr(), result.Err)
				return
			}
			results[i] = result.Output.(*image.Gray)
		case <-server.stopCtx.Done(): // Handle shutdown gracefully
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

	resultBfsChan := make(chan worker.Task[image.Gray, []geometry.Contour], 100)

	task2 := worker.Task[image.Gray, []geometry.Contour]{
		Conn:       conn,
		Input:      *cannyImage,
		ResultChan: resultBfsChan,
		Function:   FindContoursBFSWrapper,
	}
	workerChannels.bfsChan <- task2

	bfsResult := make([]geometry.Contour, 0)
	for i := 0; i < 1; /*todo server.numWorkers*/ i++ {
		select {
		case result := <-resultBfsChan:
			if result.Err != nil {
				log.Printf("Error processing image for %s: %v", conn.RemoteAddr(), result.Err)
				return
			}
			bfsResult = result.Output
		case <-server.stopCtx.Done(): // Handle shutdown gracefully
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
		case <-server.stopCtx.Done(): // Handle shutdown gracefully
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

func FindContoursBFSWrapper(gray image.Gray) ([]geometry.Contour, error) {
	return utils.FindContoursBFS(&gray), nil
}

func ApplyCannyEdgeDetectionWrapper(img image.Image) (image.Image, error) {
	return utils.ApplyCannyEdgeDetection(img.(*image.Gray)), nil
}

func GrayscaleWrapper(img image.Image) (image.Image, error) {
	return imageUtils.Grayscale(img), nil
}

// todo : reformat code
// run executes the workflow of the server: listening a file, receiving the file over the connection, treating the image, and sending the image back to client.
func (server *Server) run() {
	listener := server.listen()
	defer func(listener net.Listener) {
		// Handle listener closing gracefully
		var opErr *net.OpError
		if err := listener.Close(); err != nil && !(errors.As(err, &opErr) && !opErr.Temporary()) {
			log.Fatalf("Unexpected error closing listener: %v", err)
		}
	}(listener)

	fmt.Println("The server is running... (Press Ctrl + C to stop)")

	// Channels for managing concurrent workers
	socketSemaphore := make(chan net.Conn, 5)
	imageChan := make(chan worker.Task[image.Image, image.Image], 100)
	bfsChan := make(chan worker.Task[image.Gray, []geometry.Contour], 100)
	findQuadrilateralChan := make(chan worker.Task[[]geometry.Contour, geometry.ContourWithArea], 100)

	channels := workerChannels{
		socketSemaphore:       socketSemaphore,
		imageChan:             imageChan,
		bfsChan:               bfsChan,
		findQuadrilateralChan: findQuadrilateralChan,
	}

	// Start worker pools
	go worker.StartWorkerPool("Image Worker", numWorkers, worker.TreatmentWorker, imageChan)
	go worker.StartWorkerPool("BFS worker", numWorkers, worker.TreatmentWorker, bfsChan)
	go worker.StartWorkerPool("FindQuadrilateral worker", numWorkers, worker.TreatmentWorker, findQuadrilateralChan)

	// Use a goroutine to listen for server stop signals
	go func() {
		<-server.stopCtx.Done() // Wait for cancellation
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
			// Special case: listener.Closed() produces `use of closed network connection` error
			var opErr *net.OpError
			if errors.As(err, &opErr) && !opErr.Temporary() {
				log.Println("Listener has been closed. Stopping server gracefully.")
				return
			}
			// Other errors
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		select {
		case <-server.stopCtx.Done(): // Stop if shutdown signal received
			log.Println("Server is shutting down, closing new connection.")
			conn.Close()
		default:
			go server.handleConnection(conn, channels)
		}
	}
}

// main initialize the functionality of a TCP server.
func main() {
	// Open a file for logging
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

	// Set the output of the default logger to the file
	log.SetOutput(logFile)

	log.Println("Starting server...")

	server := newServer(host, port, numWorkers)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan // Wait for signal
		log.Println("Interrupt signal received.")
		server.cancel() // Signal server to stop
	}()

	server.run()
	log.Println("Server shut down gracefully.")
	os.Exit(0)
}
