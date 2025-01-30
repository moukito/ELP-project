package utils

import (
	"image"
	"image/color"
	"math"
)

// convertToBinaryMatrix convert an gray image to an binary image represented by a bidimensional matrix
func convertToBinaryMatrix(img *image.Gray, threshold uint8) [][]int {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Créez une matrice pour stocker les valeurs binaires
	binaryMatrix := make([][]int, height)
	for i := range binaryMatrix {
		binaryMatrix[i] = make([]int, width)
	}

	// Parcourez chaque pixel de l'image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Obtenez la valeur de gris du pixel
			grayValue := img.GrayAt(x, y).Y
			// Appliquez le seuil
			if grayValue > threshold {
				binaryMatrix[y][x] = 1 // Blanc
			} else {
				binaryMatrix[y][x] = 0 // Noir
			}
		}
	}

	return binaryMatrix
}

// findContours apply DFS to find contours of A4
func findContours(image [][]int) [][][]int {
	height, width := len(image), len(image[0])
	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}

	var contours [][][]int
	directions := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if image[y][x] == 1 && !visited[y][x] {
				contour := [][]int{}
				stack := [][]int{{y, x}}

				for len(stack) > 0 {
					point := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					py, px := point[0], point[1]

					if visited[py][px] {
						continue
					}
					visited[py][px] = true
					contour = append(contour, []int{px, py})

					for _, d := range directions {
						ny, nx := py+d[0], px+d[1]
						if ny >= 0 && ny < height && nx >= 0 && nx < width && image[ny][nx] == 1 && !visited[ny][nx] {
							stack = append(stack, []int{ny, nx})
						}
					}
				}
				contours = append(contours, contour)
			}
		}
	}
	return contours
}

// perpendicularDistance calculate the distance between one point and a line defined by two other points
func perpendicularDistance(point, lineStart, lineEnd []int) float64 {
	x0, y0 := float64(point[0]), float64(point[1])
	x1, y1 := float64(lineStart[0]), float64(lineStart[1])
	x2, y2 := float64(lineEnd[0]), float64(lineEnd[1])

	num := math.Abs((y2-y1)*x0 - (x2-x1)*y0 + x2*y1 - y2*x1)
	denom := math.Sqrt(math.Pow(y2-y1, 2) + math.Pow(x2-x1, 2))
	return num / denom
}

// douglasPeucker apply Douglas-Peucker algorithm to determine possible outline
func douglasPeucker(points [][]int, epsilon float64) [][]int {
	if len(points) < 3 {
		return points
	}

	var maxDist float64
	index := 0

	for i := 1; i < len(points)-1; i++ {
		dist := perpendicularDistance(points[i], points[0], points[len(points)-1])
		if dist > maxDist {
			maxDist = dist
			index = i
		}
	}

	if maxDist > epsilon {
		left := douglasPeucker(points[:index+1], epsilon)
		right := douglasPeucker(points[index:], epsilon)
		return append(left[:len(left)-1], right...)
	}

	return [][]int{points[0], points[len(points)-1]}
}

// pointSetAfterDouglasPeucker stock all outline from douglasPeucker
func pointSetAfterDouglasPeucker(contours [][][]int, epsilon float64) [][][]int {
	var pointSet [][][]int

	for _, contour := range contours {
		simplified := douglasPeucker(contour, epsilon)
		pointSet = append(pointSet, simplified)
	}

	return pointSet
}

// diagonal calculate the diagonal with the biggest outline
func diagonal(pointSet [][][]int) ([]int, []int) {
	norm := 0.0
	var cornerA, cornerC []int
	for _, points := range pointSet {
		for _, point1 := range points {
			for _, point2 := range points {
				vx := float64(point2[0] - point1[0])
				vy := float64(point2[1] - point1[1])
				dist := math.Sqrt(vx*vx + vy*vy)
				if dist > norm {
					norm = dist
					cornerA = point1
					cornerC = point2
				}
			}
		}
	}
	return cornerA, cornerC
}

/* func cornerConstructor(cornerA, cornerC []int) ([]int, []int, []int, []int) {
	x1, y1 := float64(cornerA[0]), float64(cornerA[1])
	x3, y3 := float64(cornerC[0]), float64(cornerC[1])

	cx := (x1 + x3) / 2
	cy := (y1 + y3) / 2

	vdx := x3 - x1
	vdy := y3 - y1

	d := math.Sqrt(vdx*vdx + vdy*vdy)
	h := d / math.Sqrt(3)
	w := h * math.Sqrt(2)

	angle := math.Atan2(vdy, vdx) // Angle en radians

	x2 := cx + (w/2)*math.Cos(angle-math.Pi/2)
	y2 := cy + (w/2)*math.Sin(angle-math.Pi/2)
	x4 := cx - (w/2)*math.Cos(angle-math.Pi/2)
	y4 := cy - (w/2)*math.Sin(angle-math.Pi/2)

	cornerB := []int{int(math.Round(x2)), int(math.Round(y2))}
	cornerD := []int{int(math.Round(x4)), int(math.Round(y4))}

	return cornerA, cornerB, cornerC, cornerD
} */

// cornerConstructor determine 4 corners of an A4 paper
func cornerConstructor(cornerA, cornerC []int) ([]int, []int, []int, []int) {
	cornerB := []int{cornerA[0], cornerC[1]}
	cornerD := []int{cornerC[0], cornerA[1]}
	return cornerA, cornerB, cornerC, cornerD
}

// isInside checks whether a point is inside the rectangle defined by the corners
func isInside(p, cornerA, cornerB, cornerC, cornerD []int) bool {
	totalArea := triangleArea(cornerA, cornerB, cornerC) + triangleArea(cornerA, cornerC, cornerD)
	areaSum := triangleArea(p, cornerA, cornerB) + triangleArea(p, cornerB, cornerC) + triangleArea(p, cornerC, cornerD) + triangleArea(p, cornerD, cornerA)
	return totalArea == areaSum
}

// triangleArea calculate the area of a triangle using the determinant formula
func triangleArea(p1, p2, p3 []int) int {
	return abs((p1[0]*(p2[1]-p3[1]) + p2[0]*(p3[1]-p1[1]) + p3[0]*(p1[1]-p2[1])) / 2)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// MaskOutsideCorners mask all pixels which are outside rectangle
func MaskOutsideCorners(img *image.Gray, threshold uint8, epsilon float64) *image.Gray {
	binaryMatrix := convertToBinaryMatrix(img, threshold)
	contours := findContours(binaryMatrix)
	pointSet := pointSetAfterDouglasPeucker(contours, epsilon)
	cornerA, cornerB, cornerC, cornerD := cornerConstructor(diagonal(pointSet))
	bounds := img.Bounds()
	newImg := image.NewGray(bounds)

	// Parcours tous les pixels de l'image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			p := []int{x, y}

			// Vérifie si le point est à l'intérieur du quadrilatère
			if isInside(p, cornerA, cornerB, cornerC, cornerD) {
				newImg.SetGray(x, y, img.GrayAt(x, y)) // Conserve la couleur originale
			} else {
				newImg.SetGray(x, y, color.Gray{0}) // Met en noir
			}
		}
	}

	return newImg
}
