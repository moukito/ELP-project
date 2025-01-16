package main

import (
	"fmt"
	"net"
	"os"
)

const bufferSize = 1024

// todo: put this in an other file
func OpenFile(path string) *os.File {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening image file:", err)
		os.Exit(1)
	}

	return file
}

// todo: put this in an other file
func Connector(network string, address string) net.Conn {
	conn, err := net.Dial(network, address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(2)
	}

	return conn
}

func sendImage(conn net.Conn, file *os.File, buffer []byte) {
	for {
		n, err := file.Read(buffer)
		if err != nil {
			break
		}
		if n > 0 {
			_, err = conn.Write(buffer[:n])
			if err != nil {
				fmt.Println("Error sending data:", err)
				os.Exit(3)
			}
		}
	}
}

// main initializes the process of reading an image file and sending it over a TCP connection to a server.
func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("Usage: ./client <image_file_path> <server_address>")
		return
	}

	imageFilePath := args[1]
	serverAddress := args[2]

	// Open the image file
	imageFile := OpenFile(imageFilePath)
	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			fmt.Println("Error closing image file:", err)
			os.Exit(1)
		}
	}(imageFile)

	// Open the tcp connection
	conn := Connector("tcp", serverAddress)
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
			os.Exit(2)
		}
	}(conn)

	buffer := make([]byte, bufferSize) // Buffer to hold chunks of the file

	fmt.Println("Sending image...")
	sendImage(conn, imageFile, buffer)
	fmt.Println("Image sent successfully!")
}
