package utils

import (
	"ELP-project/internal/geometry"
	"image"
	"image/color"
	"image/draw"
)

func DrawContour(img image.Image, contour geometry.Contour) *image.RGBA {
	bounds := img.Bounds()
	output := image.NewRGBA(bounds)

	draw.Draw(output, bounds, img, bounds.Min, draw.Src)

	red := color.RGBA{R: 255, A: 255}
	for _, p := range contour {
		output.Set(p.X, p.Y, red)
	}

	return output
}
