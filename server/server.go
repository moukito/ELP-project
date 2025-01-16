package main

import (
	"fmt"
	"net"
	"os"
)

func main() {

	listener, err := net.Listen("tcp", ":6666")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 6666...")
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}
	defer conn.Close()

	file, err := os.Create("received_image.jpg")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, 1024)

	fmt.Println("Receiving image...")
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("File received successfully")
			break
		}
		if n > 0 {
			_, err := file.Write(buffer[:n])
			if err != nil {
				fmt.Println("Error writing to file:", err)
				return
			}
		}
	}

	fmt.Println("Image saved as 'received_image.jpg'")
}
