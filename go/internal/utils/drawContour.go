package utils

/*
Package utils provides tools for image manipulation and processing, including methods for
drawing geometry and highlighting contours on images.

---

### DrawContour(img image.Image, contour geometry.Contour) *image.RGBA
Draws a contour on the given image and returns a new RGBA image with the highlighted contour.

- **Parameters**:
  - img: The input image (`image.Image`) on which the contour will be drawn.
  - contour: A `geometry.Contour` object representing the list of points that form the contour.

- **Returns**:
  - A new image (`*image.RGBA`) with the contour overlaid on the input image.

- **Behavior**:
  - Creates a new RGBA image with the same bounds as the input image.
  - Copies the content of the input image into the new image.
  - Draws each point of the contour in red (`color.RGBA{R: 255, A: 255}`) onto the new image.

- **Applications**:
  - Highlighting detected shapes or contours in an image.
  - Visualization of geometric data superimposed on images.

---

### Example Usage:
```go
package main

import (
	"ELP-project/internal/geometry"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"utils"
)

func main() {
	// Open and decode the input image
	file, _ := os.Open("input.jpg")
	defer file.Close()
	img, _, _ := image.Decode(file)

	// Define a contour to overlay on the image
	contour := geometry.Contour{
		{X: 50, Y: 50},
		{X: 100, Y: 50},
		{X: 100, Y: 100},
		{X: 50, Y: 100},
	}

	// Draw the contour on the image
	output := utils.DrawContour(img, contour)

	// Save the result
	outputFile, _ := os.Create("output.jpg")
	defer outputFile.Close()
	jpeg.Encode(outputFile, output, nil)
}
```

---

### Key Features:
- **Contour Drawing**:
  - Easily overlays contours or geometric shapes on images for visualization and debugging.
- **Red Color Highlight**:
  - Marked points in the contour are drawn in red for clear visibility.
*/

import (
	"ELP-project/internal/geometry"
	"image"
	"image/color"
	"image/draw"
)

func DrawContour(img image.Image, contour geometry.Contour) *image.RGBA {
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)

	draw.Draw(output, bounds, img, bounds.Min, draw.Src)

	red := color.RGBA{R: 255, A: 255}
	for _, p := range contour {
		output.Set(p.X, p.Y, red)
	}

	return output
}
