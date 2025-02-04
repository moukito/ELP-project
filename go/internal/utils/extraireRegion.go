package utils

/*
Package utils provides tools for image manipulation and processing, including functionality to extract a region of interest (ROI) from an image based on a specified contour.

---

### ExtractRegion(img image.Image, quad geometry.Contour) *image.RGBA
Extracts a specific region of an image defined by a quadrilateral contour.

- **Parameters**:
  - `img`: The input image (`image.Image`) from which a region will be extracted.
  - `quad`: A quadrilateral contour (`geometry.Contour`) defining the region of interest. The contour is treated as a polygon with a set of vertices.

- **Returns**:
  - A new RGBA image (`*image.RGBA`) where pixels inside the defined region retain their original values, and pixels outside the region are black.

- **Behavior**:
  - Creates a black "mask" image with the same bounds as the input image.
  - Iterates through each pixel in the image to determine if it lies within the specified quadrilateral using the helper function `isInsideQuad`.
  - Pixels inside the region are copied from the input image into the output image.
  - Pixels outside the region are set to black.

---

### isInsideQuad(x, y int, quad geometry.Contour) bool
Determines if a given point lies inside a quadrilateral contour, implementing a point-in-polygon algorithm.

- **Parameters**:
  - `x`, `y`: The coordinates of the point to test.
  - `quad`: A quadrilateral contour (`geometry.Contour`) representing the polygon.

- **Returns**:
  - `true` if the point lies inside the quadrilateral; otherwise, `false`.

- **Behavior**:
  - Uses an edge-crossing algorithm to determine the number of times a horizontal ray from the test point intersects the edges of the polygon.
  - A point is considered inside if the number of intersections is odd.

---

### Key Features:
- **Region of Interest (ROI)**:
  - Provides functionality to extract specific areas of interest from an image, preserving only the relevant content while masking the rest.
- **Geometric Contour Support**:
  - The quadrilateral contour used for ROI extraction is flexible and can define any convex or concave polygons.

---

### Example Usage:
```go
package main

import (
	"ELP-project/internal/geometry"
	"image"
	"image/color"
	"image/draw"
	"os"
	"utils"
)

func main() {
	// Open an input image
	file, _ := os.Open("input.jpg")
	defer file.Close()
	img, _, _ := image.Decode(file)

	// Define a quadrilateral contour
	quad := geometry.Contour{
		{X: 50, Y: 50},
		{X: 200, Y: 50},
		{X: 200, Y: 200},
		{X: 50, Y: 200},
	}

	// Extract the region of interest
	extractedImage := utils.ExtractRegion(img, quad)

	// Save the extracted region
	outputFile, _ := os.Create("output.png")
	defer outputFile.Close()
	png.Encode(outputFile, extractedImage)
}
```
---

### Notes:
- The function assumes that the input quad is a polygon where no two edges intersect except at the vertices.
- The output image will have the same dimensions as the input image, with irrelevant areas masked in black.
*/

import (
	"ELP-project/internal/geometry"
	"image"
	"image/color"
	"image/draw"
)

func ExtractRegion(img image.Image, quad geometry.Contour) *image.RGBA {
	bounds := img.Bounds()
	mask := image.NewRGBA(bounds)

	draw.Draw(mask, bounds, &image.Uniform{C: color.Black}, image.Point{}, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if isInsideQuad(x, y, quad) {
				mask.Set(x, y, img.At(x, y))
			}
		}
	}

	return mask
}

func isInsideQuad(x, y int, quad geometry.Contour) bool {
	count := 0
	n := len(quad)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		if (quad[i].Y > y) != (quad[j].Y > y) &&
			x < (quad[j].X-quad[i].X)*(y-quad[i].Y)/(quad[j].Y-quad[i].Y)+quad[i].X {
			count++
		}
	}
	return count%2 == 1
}
