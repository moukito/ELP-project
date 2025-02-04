package utils

/*
Package utils provides functionality for image analysis and geometry processing, including a breadth-first search (BFS) implementation for detecting contours in binary images.

---

### FindContoursBFSWithDefault(img *image.Gray) []geometry.Contour
Finds contours within the full bounds of a binary grayscale image using a breadth-first search.

- **Parameters**:
  - img: A binary grayscale image (`*image.Gray`). Non-zero pixels are treated as "white" (foreground), and zero pixels are treated as "black" (background).
- **Returns**:
  - contours: A slice of `geometry.Contour`, where each contour represents a connected component of white pixels in the image.
- **Behavior**:
  - Utilizes `FindContoursBFS` with the full bounds of the provided image as its region of interest.

---

### FindContoursBFS(img *image.Gray, bounds image.Rectangle) []geometry.Contour
Finds contours within a specific region of a binary grayscale image using BFS.

- **Parameters**:
  - img: A binary grayscale image (`*image.Gray`).
  - bounds: An `image.Rectangle` defining the region of interest in the image to process.
- **Returns**:
  - contours: A slice of `geometry.Contour`, each representing a connected component of white pixels in the specified region.
- **Behavior**:
  - Iterates over each pixel in the region defined by `bounds`.
  - For every unvisited white pixel (foreground), initiates a BFS to explore all connected white pixels, marking each as visited.
  - Explores in 8 possible directions (up, down, left, right, and diagonals) defined by the `directions` variable.
  - Connected components with fewer than 50 pixels are ignored to reduce noise.
  - Returns all identified contours with more than 50 pixels.

---

### Key Features
- **Contour Detection**:
  - Implements a BFS-based approach to find connected components with white pixels in binary images.
- **Customizable Bounds**:
  - Allows the user to limit processing to a specific rectangular region of the input image.
- **Noise Reduction**:
  - Filters out small contours (less than 50 pixels) to focus on significant components.

---

### Dependencies
- **geometry**:
  - The `geometry.Point` type is used to represent 2D points (x, y).
  - The `geometry.Contour` type represents a slice of `geometry.Point`, signifying all points in a contour.
- **imageUtils**:
  - The `imageUtils.IsWhite` function is used to determine if a given pixel in the grayscale image is part of the foreground.

---

### Example Usage
```go
package main

import (
	"ELP-project/internal/imageUtils"
	"ELP-project/internal/utils"
	"image"
	"image/png"
	"os"
)

func main() {
	// Load a binary grayscale image
	file, _ := os.Open("binary_image.png")
	defer file.Close()
	img, _, _ := image.Decode(file)
	grayImg := img.(*image.Gray)

	// Find contours in the image
	contours := utils.FindContoursBFSWithDefault(grayImg)

	// Process contours (example: print the number of detected contours)
	println("Number of contours:", len(contours))
}
```

---

### Contour Filtering
By default, only contours with more than 50 pixels are returned. Adjusting the threshold for contour size can be achieved by modifying the relevant `if` condition within the `FindContoursBFS` function.

### Key Behavior
- **8-Directional Search**:
  - Ensures all neighbors (vertical, horizontal, and diagonal) are considered during BFS traversal.
- **Memory Efficiency**:
  - Uses a `visited` map to avoid revisiting already-processed pixels and reduce redundant computation.
- **Dynamic Adaptation**:
  - Can be applied to entire images or specific subregions, enabling flexibility in use cases like ROI-specific contour detection.

*/

import (
	"ELP-project/internal/geometry"
	"ELP-project/internal/imageUtils"
	"image"
)

var directions = []geometry.Point{
	{0, 1}, {1, 0}, {0, -1}, {-1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1},
}

func FindContoursBFSWithDefault(img *image.Gray) []geometry.Contour {
	return FindContoursBFS(img, img.Bounds())
}

func FindContoursBFS(img *image.Gray, bounds image.Rectangle) []geometry.Contour {
	visited := make(map[geometry.Point]bool)
	var contours []geometry.Contour

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			p := geometry.Point{X: x, Y: y}

			if imageUtils.IsWhite(img, x, y) && !visited[p] {
				var contour geometry.Contour
				queue := []geometry.Point{p}

				for len(queue) > 0 {
					curr := queue[0]
					queue = queue[1:]

					if visited[curr] {
						continue
					}
					visited[curr] = true
					contour = append(contour, curr)

					for _, d := range directions {
						neighbor := geometry.Point{X: curr.X + d.X, Y: curr.Y + d.Y}
						if imageUtils.IsWhite(img, neighbor.X, neighbor.Y) && !visited[neighbor] {
							queue = append(queue, neighbor)
						}
					}
				}
				if len(contour) > 50 {
					contours = append(contours, contour)
				}
			}
		}
	}

	return contours
}
