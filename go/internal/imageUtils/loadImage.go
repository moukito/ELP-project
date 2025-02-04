package imageUtils

/*
Package imageUtils provides utilities for loading and processing images, including convenient methods for file handling and image decoding.

---

### LoadImage(filePath string) (image.Image, string, error)
Loads an image from the specified file path and decodes it into an `image.Image` object.

- **Parameters**:
  - filePath: The path to the image file on the disk.

- **Returns**:
  - `image.Image`: The decoded image object.
  - `string`: The format of the image (e.g., "jpeg", "png").
  - `error`: An error if the image could not be loaded or decoded.

- **Behavior**:
  - Opens the file located at `filePath`.
  - Attempts to decode the image using standard Go image decoders.
  - Closes the file after decoding.
  - If an error occurs during file opening or decoding, the error is returned.

- **Panics**:
  - If the file cannot be closed after processing, a panic will be triggered.

---

### Key Features:
- **File Decoding**:
  - Supports multiple image formats (JPEG, PNG, etc.) thanks to the standard `image` package.
- **Safe File Handling**:
  - Ensures files are closed properly after processing, even in the event of errors.

---

### Example Usage:
```go
package main

import (
	"image/jpeg"
	"imageUtils"
	"os"
)

func main() {
	// Load an image from file
	img, format, err := imageUtils.LoadImage("example.jpg")
	if err != nil {
		panic(err)
	}

	// Save the loaded image to a new file
	outFile, err := os.Create("output.jpg")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Re-encode the image as JPEG
	err = jpeg.Encode(outFile, img, nil)
	if err != nil {
		panic(err)
	}

	println("Image loaded and re-saved in format:", format)
}
```
*/

import (
	"image"
	"os"
)

func LoadImage(filePath string) (image.Image, string, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return nil, "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}
