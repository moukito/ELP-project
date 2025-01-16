package utils

import (
	"image"
	"image/color"
	"math"
)

// ApplySobelEdgeDetection detects edges using the Sobel operator.
func ApplySobelEdgeDetection(img *image.Gray) (*image.Gray, [][]float64) {
	bounds := img.Bounds()
	edgeImage := image.NewGray(bounds)

	// Sobel kernels for calculating Gx and Gy
	sobelX := [3][3]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	sobelY := [3][3]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	// Create a 2D slice to store the gradient magnitudes
	gradientMagnitudes := make([][]float64, bounds.Max.Y)
	for i := range gradientMagnitudes {
		gradientMagnitudes[i] = make([]float64, bounds.Max.X)
	}

	// Iterate over each pixel in the image
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var gx, gy float64 // Gradients in x and y directions

			// Apply the Sobel kernels
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					px := x + kx
					py := y + ky
					// Check if the kernel's position is within image bounds
					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						gray := float64(img.GrayAt(px, py).Y)
						gx += gray * sobelX[ky+1][kx+1]
						gy += gray * sobelY[ky+1][kx+1]
					}
				}
			}

			// Calculate the gradient magnitude (edge strength)
			magnitude := math.Sqrt(gx*gx + gy*gy)
			gradientMagnitudes[y][x] = magnitude

			// Normalize and set the pixel value in the edge image
			edgeImage.SetGray(x, y, color.Gray{Y: uint8(math.Min(magnitude, 255))})
		}
	}

	return edgeImage, gradientMagnitudes
}
