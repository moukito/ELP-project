package utils

import (
	"math"
)

// generateGaussianKernel crée un noyau gaussien 2D.
func generateGaussianKernel(size int, sigma float64) [][]float64 {
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

	// Normaliser le noyau pour que la somme = 1
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			kernel[i][j] /= sum
		}
	}

	return kernel
}
