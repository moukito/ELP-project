package imageUtils

/*
Package imageUtils provides utilities for image processing, with a focus on operations like
grayscale conversion for improved analysis in image manipulation tasks.

---

### Grayscale(img image.Image) *image.Gray
Converts an image to a grayscale representation based on perceptual luminance.

- **Parameters**:
  - `img`: The input image (`image.Image`) to be converted to grayscale.

- **Returns**:
  - A new grayscale image (`*image.Gray`) with the same dimensions as the input image.

- **Behavior**:
  - Iterates through every pixel of the input image.
  - Computes the grayscale value based on the formula:
    ```
    grayValue = 0.299 * Red + 0.587 * Green + 0.114 * Blue
    ```
  - Sets the calculated grayscale value (`grayValue`) to the corresponding pixel in the new grayscale image, preserving the spatial dimensions.

- **Use Case**:
  - Useful as a preprocessing step for tasks like edge detection, filtering, or any image analysis where color information is not needed.

---

### Key Features:
- **Perceptual Luminance**:
  - Grayscale conversion uses weighted contributions from each color channel (`R`, `G`, `B`) to match human visual system sensitivity.
- **Compatibility**:
  - Supports any `image.Image` compatible input, converting it to a `*image.Gray` format.

---

### Example Usage:

```go
package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"imageUtils"
)

func main() {
	// Load an input image file
	file, _ := os.Open("input.jpg")
	defer file.Close()
	img, _, _ := image.Decode(file)

	// Convert the image to grayscale
	grayImg := imageUtils.Grayscale(img)

	// Save the converted grayscale image
	outputFile, _ := os.Create("output.png")
	defer outputFile.Close()
	png.Encode(outputFile, grayImg)
}
```
*/

import (
	"image"
	"image/color"
)

func Grayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImage := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			r8, g8, b8 := r>>8, g>>8, b>>8
			grayValue := uint8(0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8))

			grayImage.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}

	return grayImage
}
