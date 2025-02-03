package utils

import "ELP-project/internal/geometry"

func FindCorner(contour geometry.Contour, center geometry.Point) geometry.Contour {
	corner1 := center
	corner2 := center

	for _, point := range contour {
		if point.X < corner1.X {
			corner1 = point
		}
		if point.Y < corner1.Y {
			corner1 = point
		}
		if point.X > corner2.X {
			corner2 = point
		}
		if point.Y > corner2.Y {
			corner2 = point
		}
	}

	return geometry.Contour{corner1, corner2}
}
