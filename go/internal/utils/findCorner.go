package utils

/*
Package utils provides tools for geometric operations including corner detection
from contours based on a reference center point.

---

### FindCorner(contour geometry.Contour, center geometry.Point) geometry.Contour
Finds the bounding corners of a contour.

- **Parameters**:
  - `contour geometry.Contour`: A collection of points representing the contour.
  - `center geometry.Point`: A reference point used as the initial bounding box coordinates.

- **Returns**:
  - `geometry.Contour`: A two-point contour containing the top-left (`corner1`) and
    bottom-right (`corner2`) corners of the bounding rectangle for the input contour.

- **Behavior**:
  - Iterates through all points in the input `contour` and updates the bounding points:
    - `corner1` is updated to ensure it holds the minimum `X` and `Y` values.
    - `corner2` is updated to ensure it holds the maximum `X` and `Y` values.
  - Effectively computes a bounding box for the entire contour.

---

### Key Features:
- **Bounding Box Calculation**:
  - Identifies the smallest rectangle that completely contains the input contour.

---

### Example Usage:
```go
package main

import (
	"ELP-project/internal/geometry"
	"utils"
	"fmt"
)

func main() {
	// Example contour
	contour := geometry.Contour{
		{X: 1, Y: 2},
		{X: 3, Y: 4},
		{X: 0, Y: 1},
		{X: 2, Y: 5},
	}

	// Reference center point
	center := geometry.Point{X: 0, Y: 0}

	// Find bounding corners
	corners := utils.FindCorner(contour, center)

	// Print corners
	fmt.Printf("Top-Left Corner: %+v\n", corners[0])
	fmt.Printf("Bottom-Right Corner: %+v\n", corners[1])
}
```
*/

import "ELP-project/internal/geometry"

func FindCorner(contour geometry.Contour, center geometry.Point) geometry.Contour {
	corner1 := center
	corner2 := center

	for _, point := range contour {
		if point.X < corner1.X {
			corner1.X = point.X
		}
		if point.Y < corner1.Y {
			corner1.Y = point.Y
		}
		if point.X > corner2.X {
			corner2.X = point.X
		}
		if point.Y > corner2.Y {
			corner2.Y = point.Y
		}
	}

	return geometry.Contour{corner1, corner2}
}
