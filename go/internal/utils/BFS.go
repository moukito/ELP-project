package utils

import (
	"ELP-project/internal/geometry"
	"ELP-project/internal/imageUtils"
	"image"
)

// Directions possibles (haut, bas, gauche, droite, diagonales)
var directions = []geometry.Point{
	{0, 1}, {1, 0}, {0, -1}, {-1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1},
}

func FindContoursBFS(img *image.Gray) []geometry.Contour {
	bounds := img.Bounds()
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
						if neighbor.X >= bounds.Min.X && neighbor.X < bounds.Max.X &&
							neighbor.Y >= bounds.Min.Y && neighbor.Y < bounds.Max.Y &&
							imageUtils.IsWhite(img, neighbor.X, neighbor.Y) && !visited[neighbor] {
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
