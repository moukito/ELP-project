package geometry

type Point struct {
	X, Y int
}

type Contour []Point

type ContourWithArea struct {
	Contour Contour
	Area    float64
}
