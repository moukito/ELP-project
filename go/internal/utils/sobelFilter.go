package utils

import (
	"image"
	"image/color"
	"math"
)

// GenerateSobelKernel génère dynamiquement un noyau de Sobel de taille variable.
func GenerateSobelKernel(size int) ([][]float64, [][]float64) {
	if size%2 == 0 {
		panic("Sobel kernel size must be odd")
	}

	// Si la taille est 3x3, utiliser le noyau standard
	if size == 3 {
		return [][]float64{
				{-1, 0, 1},
				{-2, 0, 2},
				{-1, 0, 1},
			}, [][]float64{
				{-1, -2, -1},
				{0, 0, 0},
				{1, 2, 1},
			}
	}

	// Si la taille est 5x5, utiliser la version fixe du noyau
	if size == 5 {
		return [][]float64{
				{-2, -1, 0, 1, 2},
				{-3, -2, 0, 2, 3},
				{-4, -3, 0, 3, 4},
				{-3, -2, 0, 2, 3},
				{-2, -1, 0, 1, 2},
			}, [][]float64{
				{-2, -2, -4, -2, -2},
				{-1, -1, -2, -1, -1},
				{0, 0, 0, 0, 0},
				{1, 1, 2, 1, 1},
				{2, 2, 4, 2, 2},
			}
	}

	kernelX := make([][]float64, size)
	kernelY := make([][]float64, size)
	radius := size / 2
	sigma := float64(size) / 3 // Ajustement empirique du sigma

	// Création d'une approximation du gradient de Sobel
	for i := 0; i < size; i++ {
		kernelX[i] = make([]float64, size)
		kernelY[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			x, y := float64(j-radius), float64(i-radius)

			// Approximation des dérivées de Sobel en X et Y
			kernelX[i][j] = -x * math.Exp(-(x*x+y*y)/(2*sigma*sigma))
			kernelY[i][j] = -y * math.Exp(-(x*x+y*y)/(2*sigma*sigma))
		}
	}

	// Normalisation des noyaux
	normalizeKernel(kernelX)
	normalizeKernel(kernelY)

	return kernelX, kernelY
}

// Normalise un noyau pour que la somme absolue des valeurs soit 1
func normalizeKernel(kernel [][]float64) {
	sum := 0.0
	for _, row := range kernel {
		for _, val := range row {
			sum += math.Abs(val)
		}
	}

	if sum != 0 {
		for i := range kernel {
			for j := range kernel[i] {
				kernel[i][j] /= sum
			}
		}
	}
}

func ComputeDynamicThresholds(img *image.Gray, alpha float64) (float64, float64) {
	bounds := img.Bounds()
	totalGradient := 0.0
	count := 0

	// Appliquer Sobel pour obtenir la magnitude des gradients
	sobelX, sobelY := GenerateSobelKernel(5)
	gradient, _ := ApplySobelEdgeDetection(img, sobelX, sobelY)

	// Calcul de la moyenne des gradients
	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			totalGradient += float64(gradient.GrayAt(x, y).Y)
			count++
		}
	}

	meanGradient := totalGradient / float64(count)

	// Définition des seuils bas et haut
	highThreshold := alpha * meanGradient
	lowThreshold := 0.4 * highThreshold

	return lowThreshold, highThreshold
}

// ApplySobelEdgeDetection applique un noyau de Sobel dynamique
func ApplySobelEdgeDetection(img *image.Gray, kernelX, kernelY [][]float64) (*image.Gray, [][]float64) {
	bounds := img.Bounds()
	output := image.NewGray(bounds)
	gradientAngles := make([][]float64, bounds.Max.Y)
	radius := len(kernelX) / 2

	// Allocation des tableaux
	for i := range gradientAngles {
		gradientAngles[i] = make([]float64, bounds.Max.X)
	}

	// Convolution du noyau Sobel
	for y := bounds.Min.Y + radius; y < bounds.Max.Y-radius; y++ {
		for x := bounds.Min.X + radius; x < bounds.Max.X-radius; x++ {
			var gx, gy float64

			for ky := -radius; ky <= radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					px := x + kx
					py := y + ky

					gray := float64(img.GrayAt(px, py).Y)
					gx += gray * kernelX[ky+radius][kx+radius]
					gy += gray * kernelY[ky+radius][kx+radius]
				}
			}

			// Magnitude et angle du gradient
			magnitude := math.Sqrt(gx*gx + gy*gy)
			angle := math.Atan2(gy, gx) * (180 / math.Pi)

			output.SetGray(x, y, color.Gray{Y: uint8(math.Min(magnitude, 255))})
			gradientAngles[y][x] = angle
		}
	}

	return output, gradientAngles
}
