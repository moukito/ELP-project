package geometry

/*
Package geometry provides data structures to represent points, contours, and contours with associated geometric properties.

---

### Point
Represents a coordinate in 2D space.

- **Fields**:
  - `X`: The X-coordinate of the point (integer).
  - `Y`: The Y-coordinate of the point (integer).

---

### Contour
A type alias representing a collection of `Point`, intended to define a closed shape or outline.

- **Behavior**:
  - A `Contour` represents the sequence of points making up a polygon or curve.
  - Can be used to define boundaries or regions in geometrical computations.

---

### ContourWithArea
Extends the concept of a `Contour` by associating it with an area measurement.

- **Fields**:
  - `Contour`: The sequence of `Point` instances representing the shape.
  - `Area`: The area enclosed by the contour (float64).

- **Usage**:
  - Useful in applications requiring both the shape and the measure of the enclosed surface, such as in computational geometry or image analysis.

---

### Example Usage:
```go
package main

import (
	"fmt"
	"geometry"
)

func main() {
	// Define points
	p1 := geometry.Point{X: 0, Y: 0}
	p2 := geometry.Point{X: 5, Y: 0}
	p3 := geometry.Point{X: 5, Y: 5}
	p4 := geometry.Point{X: 0, Y: 5}

	// Create a contour
	contour := geometry.Contour{p1, p2, p3, p4}

	// Create a ContourWithArea
	contourWithArea := geometry.ContourWithArea{
		Contour: contour,
		Area:    25.0, // Computed as an example
	}

	fmt.Printf("Contour: %v\n", contourWithArea.Contour)
	fmt.Printf("Area: %f\n", contourWithArea.Area)
}
```
*/

type Point struct {
	X, Y int
}

type Contour []Point

type ContourWithArea struct {
	Contour Contour
	Area    float64
}
