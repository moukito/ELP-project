package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

const bufferSize = 1024

// Client is a struct that encapsulates the functionality for sending files over a TCP connection.
type Client struct {
	host string
	port string
}

// NewClient creates a new Client instance with the given image file path and server address.
func NewClient(host string, port string) *Client {
	return &Client{
		host: host,
		port: port,
	}
}

// OpenFile opens the given file and returns a pointer to os.File.
func OpenFile(imageFilePath string) *os.File {
	file, err := os.Open(imageFilePath)
	if err != nil {
		log.Fatalf("error opening image file: %v", err)
	}

	return file
}

// connect establishes a TCP connection to the server.
func (client *Client) connect() net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", client.host, client.port))
	if err != nil {
		log.Fatalf("error connecting to server: %v", err)
	}

	return conn
}

// sendImage sends the image file to the server using the given connection.
func (client *Client) sendImage(file *os.File, conn net.Conn) {
	buffer := make([]byte, bufferSize)

	for {
		n, err := file.Read(buffer)
		if n > 0 {
			_, writeErr := conn.Write(buffer[:n])
			if writeErr != nil {
				log.Fatalf("Error sending data: %v", writeErr)
			}
		}

		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatalf("Error reading file: %v", err)
		}
	}
}

// Run executes the workflow of the client: opening a file, connecting, and sending the file.
func (client *Client) Run(imageFilePath string) {
	// Open image file
	file := OpenFile(imageFilePath)
	log.Printf("Image file opened: %s", file.Name())
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(file)

	// Connect to server
	conn := client.connect()
	log.Printf("Connected to server: %s", conn.RemoteAddr().String())
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Error closing connection: %v", err)
		}
	}(conn)

	log.Println("Sending image...")
	client.sendImage(file, conn)
	log.Println("Image sent successfully!")
}

// main initializes the process of reading an image file and sending it over a TCP connection to a server.
func main() {
	// Open a file for logging
	logFile, err := os.OpenFile("client.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
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

	args := os.Args

	if len(args) > 3 || len(args) < 2 {
		fmt.Println("Usage: ./client <image_file_path> <server_address>")
		log.Fatal("Invalid number of arguments")
	}

	imageFilePath := args[1]
	log.Printf("Image file path: %s", imageFilePath)

	host := "localhost"
	port := "14750"
	if len(args) == 3 {
		tmpHost, tmpPort, err := net.SplitHostPort(args[2])
		if err != nil {
			log.Fatalf("Invalid server address format: %v", err)
		}
		host = tmpHost
		port = tmpPort
	}
	log.Printf("Server address: %s:%s", host, port)

	// Create and run the client
	client := NewClient(host, port)
	client.Run(imageFilePath)
}
