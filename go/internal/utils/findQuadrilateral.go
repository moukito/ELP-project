package utils

import (
	"ELP-project/internal/geometry"
	"math"
)

func FindQuadrilateral(contours []geometry.Contour) geometry.Contour {
	var bestQuad geometry.Contour
	maxArea := 0.0

	for _, contour := range contours {
		//preprocessed := reduceByDistance(contour, 10.0) // Minimum 10 pixels de distance entre points
		//approx := DouglasPucker(contour /*preprocessed*/, 50.0)
		approx := contour
		if true /*len(approx) == 4*/ {
			area := polygonArea(approx)
			if area > maxArea {
				maxArea = area
				bestQuad = approx
			}
		}
	}
	return bestQuad
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

func reduceByDistance(points geometry.Contour, minDist float64) geometry.Contour {
	var reduced geometry.Contour
	reduced = append(reduced, points[0]) // Conservez toujours le premier point

	for i := 1; i < len(points)-1; i++ {
		prevPoint := reduced[len(reduced)-1]
		dist := math.Hypot(float64(points[i].X-prevPoint.X), float64(points[i].Y-prevPoint.Y))

		if dist > minDist {
			reduced = append(reduced, points[i])
		}
	}

	reduced = append(reduced, points[len(points)-1]) // Conservez toujours le dernier point
	return reduced
}
