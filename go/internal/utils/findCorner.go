package utils

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
