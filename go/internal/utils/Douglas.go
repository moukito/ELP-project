package utils

import (
	"ELP-project/internal/geometry"
	"math"
)

func DouglasPucker(points geometry.Contour, epsilon float64) geometry.Contour {
	if len(points) < 3 {
		return points
	}

	var maxDist float64
	index := 0
	first, last := points[0], points[len(points)-1]

	for i := 1; i < len(points)-1; i++ {
		dist := perpendicularDistance(points[i], first, last)
		if dist > maxDist {
			maxDist = dist
			index = i
		}
	}

	if maxDist > epsilon {
		firstPart := DouglasPucker(points[:index+1], epsilon)
		secondPart := DouglasPucker(points[index:], epsilon)

		return append(firstPart[:len(firstPart)-1], secondPart...)
	} else {
		return []geometry.Point{first, last}
	}
}

func perpendicularDistance(point geometry.Point, lineStart geometry.Point, lineEnd geometry.Point) float64 {
	numerator := math.Abs(float64((lineEnd.Y-lineStart.Y)*point.X - (lineEnd.X-lineStart.X)*point.Y + lineEnd.X*lineStart.Y - lineEnd.Y*lineStart.X))
	denominator := math.Sqrt(math.Pow(float64(lineEnd.Y-lineStart.Y), 2) + math.Pow(float64(lineEnd.X-lineStart.X), 2))
	return numerator / denominator
}
