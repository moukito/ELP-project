package utils

import (
	"ELP-project/internal/geometry"
	"image"
	"image/color"
	"image/draw"
)

func ExtractRegion(img image.Image, quad geometry.Contour) *image.RGBA {
	bounds := img.Bounds()
	mask := image.NewRGBA(bounds)

	draw.Draw(mask, bounds, &image.Uniform{C: color.Black}, image.Point{}, draw.Src)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if isInsideQuad(x, y, quad) {
				mask.Set(x, y, img.At(x, y))
			}
		}
	}

	return mask
}

func isInsideQuad(x, y int, quad geometry.Contour) bool {
	count := 0
	n := len(quad)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		if (quad[i].Y > y) != (quad[j].Y > y) &&
			x < (quad[j].X-quad[i].X)*(y-quad[i].Y)/(quad[j].Y-quad[i].Y)+quad[i].X {
			count++
		}
	}
	return count%2 == 1
}
