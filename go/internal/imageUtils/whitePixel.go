package imageUtils

/*
Package imageUtils provides utility functions for image processing, including pixel value checks.

---

### IsWhite(img *image.Gray, x, y int) bool
Checks if the pixel at (x, y) in a grayscale image is considered "white".

- **Parameters**:
  - img: A grayscale image (`*image.Gray`) where pixels are evaluated.
  - x: The x-coordinate of the pixel.
  - y: The y-coordinate of the pixel.
- **Returns**:
  - A boolean value (`true` if the pixel is "white", otherwise `false`).

- **Behavior**:
  - Accesses the pixel value at the specified coordinates.
  - Compares the grayscale value (`Y`) to 128 (on a scale of 0 to 255).
  - If `Y > 128`, the pixel is considered "white" and the function returns `true`.
  - Otherwise, the function returns `false`.

---

### Key Features:
- Provides a simple way to evaluate pixel brightness in grayscale images.
- Threshold-based approach ensures consistency across different images.

---

### Example Usage:
```go
package main

import (
	"fmt"
	"image"
	"image/color"
	"imageUtils"
)

func main() {
	// Create a simple grayscale image
	img := image.NewGray(image.Rect(0, 0, 10, 10))

	// Set a pixel to a specific gray value
	img.SetGray(5, 5, color.Gray{Y: 200})

	// Check if the pixel at (5, 5) is "white"
	isWhite := imageUtils.IsWhite(img, 5, 5)
	fmt.Printf("Is the pixel at (5, 5) white? %v\n", isWhite) // Output: true
}
```
*/

import "image"

func IsWhite(img *image.Gray, x, y int) bool {
	return img.GrayAt(x, y).Y > 128
}
