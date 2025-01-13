package utils

import (
	"image"
	"image/color"
	"math"
)

// GenerateGaussianKernel dynamically creates a Gaussian kernel of any size and sigma.
func GenerateGaussianKernel(size int, sigma float64) [][]float64 {
	if size%2 == 0 {
		panic("Gaussian kernel size must be odd")
	}

	kernel := make([][]float64, size)
	sum := 0.0
	radius := size / 2

	for i := 0; i < size; i++ {
		kernel[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			x, y := float64(i-radius), float64(j-radius)
			kernel[i][j] = (1 / (2 * math.Pi * sigma * sigma)) * math.Exp(-(x*x+y*y)/(2*sigma*sigma))
			sum += kernel[i][j]
		}
	}

	// Normalize the kernel
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			kernel[i][j] /= sum
		}
	}

	return kernel
}

// ApplyKernel applies a convolution kernel to a grayscale image (e.g., Gaussian blur or Sobel).
func ApplyKernel(img *image.Gray, kernel [][]float64) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)
	radius := len(kernel) / 2

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sum float64
			var weightSum float64

			// Apply kernel around the pixel
			for ky := -radius; ky <= radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					pixelX := x + kx
					pixelY := y + ky
					if pixelX >= bounds.Min.X && pixelX < bounds.Max.X && pixelY >= bounds.Min.Y && pixelY < bounds.Max.Y {
						gray := float64(img.GrayAt(pixelX, pixelY).Y)
						sum += gray * kernel[ky+radius][kx+radius]
						weightSum += kernel[ky+radius][kx+radius]
					}
				}
			}

			// Store the result
			output.SetGray(x, y, color.Gray{Y: uint8(sum / weightSum)})
		}
	}

	return output
}
