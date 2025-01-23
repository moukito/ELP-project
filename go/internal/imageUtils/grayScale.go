package imageUtils

import (
	"image"
	"image/color"
)

// Grayscale converts an image to grayscale.
func Grayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImage := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			// Standard luminance formula: 0.299*R + 0.587*G + 0.114*B
			grayValue := uint8((299*r + 587*g + 114*b) / 1000 >> 8)
			grayImage.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}

	return grayImage
}

func GrayscaleWrapper(img image.Image) (image.Image, error) {
	return Grayscale(img), nil
}
