package utils

/*
Package utils provides functionality for processing contours and identifying geometrical shapes, such as the largest quadrilateral in a given set of contours.

---

### FindQuadrilateral(contours []geometry.Contour) geometry.ContourWithArea
Identifies the largest quadrilateral in a provided set of contours by calculating the area of each.

- **Parameters**:
  - `contours`: A slice of contours (`[]geometry.Contour`), each represented as a closed sequence of 2D points.
- **Returns**:
  - `geometry.ContourWithArea`: A structure containing the largest quadrilateral (`Contour`) and its associated area (`Area`).

#### Behavior:
- Iterates through the list of contours.
- Calculates the area of each contour using the `polygonArea` function.
- Identifies the contour with the maximum area as the best quadrilateral.
- Returns the largest quadrilateral along with its area.

#### Example Usage:
```go
var contours []geometry.Contour = ... // Load or generate contours
largestQuad := utils.FindQuadrilateral(contours)
fmt.Printf("Area of largest quadrilateral: %f\n", largestQuad.Area)
```

---

### polygonArea(points geometry.Contour) float64
Calculates the area of a given polygon represented by a contour.

- **Parameters**:
  - `points`: A sequence of 2D points forming a closed polygon (`geometry.Contour`).
- **Returns**:
  - `float64`: The computed area of the polygon.

#### Behavior:
- Implements the shoelace formula to efficiently compute the area of a polygon.
- Ensures that the returned value is always positive by taking the absolute value.

#### Example:
```go
area := utils.polygonArea(contour)
fmt.Printf("Area: %f\n", area)
```

---

### Key Features:
- **Contour Processing**:
  - Simplifies contours by ensuring a minimum distance between points.
- **Polygon Area Computation**:
  - Efficiently calculates the area of 2D polygons using a straightforward algorithm.
- **Shape Identification**:
  - Identifies the largest quadrilateral in a set of contours by comparing their areas.

---

### Example Workflow:
```go
package main

import (
	"ELP-project/internal/geometry"
	"utils"
	"fmt"
)

func main() {
	// Load or generate contours (example data)
	contours := []geometry.Contour{ ... }

// Find the largest quadrilateral
bestQuadrilateral := utils.FindQuadrilateral(contours)

// Output the result
fmt.Printf("Largest Quadrilateral Area: %f\n", bestQuadrilateral.Area)
}
```
*/

import (
	"ELP-project/internal/geometry"
	"math"
)

func FindQuadrilateral(contours []geometry.Contour) geometry.ContourWithArea {
	var bestQuad geometry.Contour
	maxArea := 0.0

	for _, contour := range contours {
		area := polygonArea(contour)
		if area > maxArea {
			maxArea = area
			bestQuad = contour
		}
	}
	return geometry.ContourWithArea{Contour: bestQuad, Area: maxArea}
}

func polygonArea(points geometry.Contour) float64 {
	n := len(points)
	area := 0.0

	for i := 0; i < n; i++ {
		j := (i + 1) % n
		area += float64(points[i].X*points[j].Y - points[j].X*points[i].Y)
	}

	return math.Abs(area) / 2.0
}
