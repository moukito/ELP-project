package main

/*
Package main implements a TCP client for sending and receiving image files to/from a server.
The client connects to a specified server, sends an image file,
and then receives a processed image from the server, saving it locally.

---

### Key Features
- **Server Connection**:
  - Connects to a TCP server for communication.
  - Default server address is `localhost:14750`.
- **Image File Transmission**:
  - Sends an image file to the server using a buffered approach.
  - Receives the processed image file from the server and saves it locally.
- **Dynamic File Handling**:
  - If a file with the same output name exists, generates a new name to avoid overwriting.

---

### Constants

- `defaultHost`: The default hostname of the server (`"localhost"`).
- `defaultPort`: The default port of the server (`"14750"`).
- `bufferSize`: Buffer size (in bytes) used for reading/writing data (`1024`).

---

### Types

#### `Client`
Defines the TCP client for communication with the server.

- **Fields**:
  - `host string`: The server's hostname.
  - `port string`: The server's port.

- **Methods**:
  - `connect() net.Conn`: Establishes a connection to the server and returns the connection object.
  - `sendImage(file *os.File, conn net.Conn)`: Sends the specified image file to the server.
  - `receiveImage(conn net.Conn, file *os.File)`: Receives the processed image from the server and saves it locally.
  - `run(imageFilePath string)`: Coordinates the process of connecting, sending, and receiving.

---

### Functions

#### `newClient(host string, port string) *Client`
Creates and initializes a new instance of `Client`.

- **Parameters**:
  - `host string`: Hostname of the server.
  - `port string`: Port of the server.
- **Returns**:
  - A pointer to a new `Client` instance.

#### `Client.connect() net.Conn`
Connects to the specified server and returns the established connection.

- **Panics**:
  - If the connection fails.

#### `Client.sendImage(file *os.File, conn net.Conn)`
Sends the given image file to the server using the specified connection.

- **Parameters**:
  - `file *os.File`: The file object of the image to send.
  - `conn net.Conn`: The connection object.

#### `Client.receiveImage(conn net.Conn, file *os.File)`
Receives a file from the server and writes it to the specified file object.

- **Parameters**:
  - `conn net.Conn`: The connection object.
  - `file *os.File`: The output file object where data is written.

---

### Main Functionality

#### `main()`
The entry point of the application.

- **Behavior**:
  - Validates command-line arguments to ensure proper usage.
  - Parses the image file path and (optionally) the server address from arguments.
  - Creates a `Client` instance and manages the workflow:
    1. Opens the image file.
    2. Connects to the server.
    3. Sends the image to the server.
    4. Receives the processed image from the server and saves it with an appropriate name.
  - Logs all activities to the file `client.log`.

---

### Example Usage
```bash
# Run the client with the image file and optional server address
./client path/to/image.png localhost:14750
```

---

### Workflow Steps
1. **Initialization**:
   - The client accepts an image file path and an optional server address as command-line arguments.
   - If the server address is not provided, the default address (`localhost:14750`) is used.
2. **Connection**:
   - Establishes a TCP connection to the server.
3. **Data Transmission**:
   - Reads the image file in chunks of `bufferSize` bytes and sends it to the server.
   - A special "EOF" marker is sent to indicate the end of the file.
4. **Receiving Processed Image**:
   - Reads the processed image data from the server and writes it to a local file.
   - If the output file already exists, a new filename is generated to avoid overwriting.
5. Logs all activities (including errors) to a log file named `client.log`.

---

### File Handling
- The client ensures proper cleanup:
  - Opens files for reading or writing.
  - Closes files and network connections gracefully on completion or error.

---

### Error Handling
- Handles network errors (e.g., connection failures, data transmission errors) and file I/O errors.
- Ensures proper logging of all encountered errors.

---

### Example Workflow in Code
```go
func main() {
    // Parse arguments
    imageFilePath := "example.png"
    host := "localhost"
    port := "14750"

    // Create a new client
    client := newClient(host, port)
    client.run(imageFilePath)
}
```
*/

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
)

const (
	defaultHost = "localhost"
	defaultPort = "14750"
	bufferSize  = 1024
)

type Client struct {
	host string
	port string
}

func newClient(host string, port string) *Client {
	return &Client{
		host: host,
		port: port,
	}
}

func (client *Client) connect() net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", client.host, client.port))
	if err != nil {
		log.Fatalf("error connecting to server: %v", err)
	}

	return conn
}

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
	_, err := conn.Write([]byte("EOF"))
	if err != nil {
		log.Fatalf("Error sending EOF: %v", err)
	}
}

func (client *Client) receiveImage(conn net.Conn, file *os.File) {
	buffer := make([]byte, bufferSize)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err.Error() == "EOF" || err == io.EOF {
				break
			}
			log.Fatalf("Error reading from connection: %v", err)
		}

		_, writeErr := file.Write(buffer[:n])
		if writeErr != nil {
			log.Fatalf("Error writing to file: %v", writeErr)
		}
	}
}

func (client *Client) run(imageFilePath string) {
	file, err := os.Open(imageFilePath)
	if err != nil {
		log.Fatalf("error opening image file: %v", err)
	}
	log.Printf("Image file opened: %s", file.Name())
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(file)

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

	newFileName := "output_" + filepath.Base(file.Name())
	fileIndex := 1
	for {
		if _, err := os.Stat(newFileName); os.IsNotExist(err) {
			break
		} else {
			newFileName = fmt.Sprintf("output_%d_%s", fileIndex, filepath.Base(file.Name()))
			fileIndex++
		}
	}

	newFile, err := os.Create(newFileName)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer func(newFile *os.File) {
		err := newFile.Close()
		if err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}(newFile)

	log.Println("Receiving image...")
	client.receiveImage(conn, newFile)
}

func main() {
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

	log.SetOutput(logFile)

	args := os.Args

	if len(args) > 3 || len(args) < 2 {
		fmt.Println("Usage: ./client <image_file_path> <server_address>")
		log.Fatal("Invalid number of arguments")
	}

	imageFilePath := args[1]
	log.Printf("Image file path: %s", imageFilePath)

	host := defaultHost
	port := defaultPort
	if len(args) == 3 {
		tmpHost, tmpPort, err := net.SplitHostPort(args[2])
		if err != nil {
			log.Fatalf("Invalid server address format: %v", err)
		}
		host = tmpHost
		port = tmpPort
	}
	log.Printf("Server address: %s:%s", host, port)

	client := newClient(host, port)
	client.run(imageFilePath)
}
