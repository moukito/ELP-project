package imageUtils

/*
Package imageUtils provides functionality for saving images in various formats, such as JPEG and PNG.

---

### SaveImage(img image.Image, filePath string, format string) error
Saves an image to the specified file in the given format.

- **Parameters**:
  - `img`: The image to save. Must implement the `image.Image` interface.
  - `filePath`: The path where the image file will be saved.
  - `format`: The format of the image to save. Supported formats are "jpg", "jpeg", and "png".

- **Returns**:
  - An error (`error`) if any issue occurs during the saving process, or `nil` if the operation succeeds.

- **Behavior**:
  - The function creates a file at the specified `filePath` and saves the provided image in the specified format.
  - If the format is "jpg" or "jpeg", the image is saved in JPEG format using the `image/jpeg` package.
  - If the format is "png", the image is saved in PNG format using the `image/png` package.
  - If the format is unsupported, the function returns an error indicating the unsupported format.
  - Closes the created file after saving the image.

- **Panics**:
  - If the file fails to close after being written.

---

### Supported Formats:
- **JPEG**:
  - Extensions: `jpg`, `jpeg`
  - Saves the image using the `jpeg.Encode` function with default encoding options.
- **PNG**:
  - Extension: `png`
  - Saves the image using the `png.Encode` function.
- **Unsupported Formats**:
  - Returns an error with a message indicating the format is not supported.

---

### Example Usage:
```go
package main

import (
	"image"
	"image/color"
	"image/draw"
	"imageUtils"
)

func main() {
	// Create a simple image with a colored rectangle
	rect := image.Rect(0, 0, 200, 100)
	img := image.NewRGBA(rect)
	draw.Draw(img, rect, &image.Uniform{C: color.RGBA{255, 0, 0, 255}}, image.Point{}, draw.Src)

	// Save the image as PNG
	err := imageUtils.SaveImage(img, "output.png", "png")
	if err != nil {
		panic(err)
	}

	// Save the image as JPEG
	err = imageUtils.SaveImage(img, "output.jpg", "jpg")
	if err != nil {
		panic(err)
	}
}
```

---

### Key Features:
- **Multi-format Support**:
  - Support for saving images as PNG and JPEG.
- **Simple Interface**:
  - Standardized function for saving images in different formats.
- **Error Handling**:
  - Returns descriptive errors if the file cannot be created, a format is unsupported, or an encoding operation fails.
*/

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func SaveImage(img image.Image, filePath string, format string) error {
	file, err := os.Create(filePath)

	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	switch strings.ToLower(format) {
	case "jpg", "jpeg":
		return jpeg.Encode(file, img, nil)
	case "png":
		return png.Encode(file, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
