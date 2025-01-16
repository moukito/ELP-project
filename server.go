package main

import (
	"fmt"
	"net"
	"os"
)

const bufferSize = 1024

func ListenAttempt(serverAddress string) net.Listener {
	listener, err := net.Listen("tcp", serverAddress)

	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}

	fmt.Println("Server is listening on %s", serverAddress)

	return listener
}

func ConnAttempt(listener net.Listener) net.Conn {
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		os.Exit(2)
	}

	fmt.Println("Data incoming...")

	return conn
}

func ReceiveImgAttempt(conn net.Conn) *os.File {
	file, err := os.Create("received_image.jpg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(3)
	}

	return file
}

func ReceiveImg(file *os.File, conn net.Conn, buffer []byte) {
	fmt.Println("Receiving image...")

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Image received sucessfully!")
			break
		}
		if n > 0 {
			_, err := file.Write(buffer[:n])
			if err != nil {
				fmt.Println("Error writing to file:", err)
				os.Exit(4)
			}
		}
	}

	fmt.Println("Image saved as 'received_image.jpg'")
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: ./server <server_address>")
		return
	}

	serverAddress := args[1] // format: host:port

	listener := ListenAttempt(serverAddress)

	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Println("Error closing server:", err)
			os.Exit(1)
		}
	}(listener)

	conn := ConnAttempt(listener)

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
			os.Exit(2)
		}
	}(conn)

	file := ReceiveImgAttempt(conn)

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing received file:", err)
		}
	}(file)

	buffer := make([]byte, bufferSize) // Buffer to hold chunks of the file

	ReceiveImg(file, conn, buffer)

	fmt.Println("Image saved as 'received_image.jpg'")
}
