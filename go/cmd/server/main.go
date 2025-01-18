package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net"
	"os"
)

const (
	host       = "localhost"
	port       = "14750"
	protocol   = "tcp"
	bufferSize = 1024
)

type Server struct {
	host string
	port string
}

func NewServer(host string, port string) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (server *Server) listen() net.Listener {
	listener, err := net.Listen(protocol, fmt.Sprintf("%s:%s", server.host, server.port))
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Println("Server is listening on port 6666...")

	return listener
}

func (server *Server) receiveImage(conn net.Conn) (image.Image, string) {
	// Create a buffer to store the incoming data
	var dataBuffer bytes.Buffer
	buffer := make([]byte, bufferSize) // Temporary buffer size for chunks

	for {
		// Read incoming data into the temporary buffer
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println(err.Error())
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

		// Advance the sent position
		sent += n
	}

	log.Printf("Image sent successfully. Total bytes: %d", dataLen)
}

func (server *Server) run() {
	listener := server.listen()
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			log.Fatalf("Error closing listener: %v", err)
		}
	}(listener)

	for {
		conn, err := listener.Accept()
		// todo : workers
		if err != nil {
			log.Println("Error accepting connection:", err)
		}

		log.Println("Receiving image...")
		img, format := server.receiveImage(conn)
		log.Println("Image received successfully!")

		// todo : treat the image

		log.Println("Sending image...")
		server.sendImage(conn, img, format)

		conn.Close()
	}
}

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

	server := NewServer(host, port)
	server.run()
}
