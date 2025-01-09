package utils

import (
	"math"
)

// GenerateGaussianKernel crée un noyau gaussien 2D.
func GenerateGaussianKernel(size int, sigma float64) [][]float64 {
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

// ApplyGaussianFilter applique un filtre gaussien à une matrice RGBA.
func ApplyGaussianFilter(matrix [][][4]uint8, kernel [][]float64) [][][4]uint8 {
	height := len(matrix)
	width := len(matrix[0])
	radius := len(kernel) / 2

	// Créer une matrice pour le résultat
	result := make([][][4]uint8, height)
	for i := range result {
		result[i] = make([][4]uint8, width)
	}

	// Parcourir chaque pixel
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Variables pour accumuler les valeurs RGB
			var r, g, b, a float64

			// Appliquer le noyau gaussien
			for ky := -radius; ky < radius; ky++ {
				for kx := -radius; kx < radius; kx++ {
					ny, nx := y+ky, x+kx
					if ny >= 0 && ny < height && nx >= 0 && nx < width {
						weight := kernel[ky+radius][kx+radius]
						pixel := matrix[ny][nx]
						r += weight * float64(pixel[0])
						g += weight * float64(pixel[1])
						b += weight * float64(pixel[2])
						a += weight * float64(pixel[3])
					}
				}
			}

			// Stocker le résultat
			result[y][x] = [4]uint8{uint8(r), uint8(g), uint8(b), uint8(a)}
		}
	}

	return result
}
