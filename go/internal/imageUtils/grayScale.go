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
			r8, g8, b8 := r>>8, g>>8, b>>8 // Convert to 8-bit
			grayValue := uint8(0.299*float64(r8) + 0.587*float64(g8) + 0.114*float64(b8))

			// Updating the pixel in the grayscale image.
			grayImage.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}

	return grayImage
}
