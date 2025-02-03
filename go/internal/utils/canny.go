package utils

import (
	"image"
	"image/color"
)

// Non-maximum suppression
func nonMaxSuppression(gradient image.Gray, angles [][]float64) *image.Gray {
	bounds := gradient.Bounds()
	suppressed := image.NewGray(bounds)

	for y := 1; y < bounds.Max.Y-1; y++ {
		for x := 1; x < bounds.Max.X-1; x++ {
			angle := angles[y][x]
			mag := gradient.GrayAt(x, y).Y
			n1, n2 := uint8(0), uint8(0)

			// Quantize the angle to 4 directions (0, 45, 90, 135 degrees)
			if (angle >= -22.5 && angle <= 22.5) || (angle >= 157.5 || angle <= -157.5) {
				n1, n2 = gradient.GrayAt(x-1, y).Y, gradient.GrayAt(x+1, y).Y
			} else if (angle > 22.5 && angle <= 67.5) || (angle < -112.5 && angle >= -157.5) {
				n1, n2 = gradient.GrayAt(x-1, y-1).Y, gradient.GrayAt(x+1, y+1).Y
			} else if (angle > 67.5 && angle <= 112.5) || (angle < -67.5 && angle >= -112.5) {
				n1, n2 = gradient.GrayAt(x, y-1).Y, gradient.GrayAt(x, y+1).Y
			} else {
				n1, n2 = gradient.GrayAt(x-1, y+1).Y, gradient.GrayAt(x+1, y-1).Y
			}

			// Suppress non-maximum values
			if mag >= n1 && mag >= n2 {
				suppressed.SetGray(x, y, color.Gray{Y: mag})
			} else {
				suppressed.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
	return suppressed
}

// Seuil d'hystérésis
func hysteresisThresholding(img *image.Gray, lowThreshold, highThreshold float64) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)

	strong := uint8(255)
	weak := uint8(75)

	// Étape 1 : Classification des pixels
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.GrayAt(x, y).Y
			if float64(pixel) >= highThreshold {
				output.SetGray(x, y, color.Gray{Y: strong}) // Contour fort
			} else if float64(pixel) >= lowThreshold {
				output.SetGray(x, y, color.Gray{Y: weak}) // Contour faible
			} else {
				output.SetGray(x, y, color.Gray{Y: 0}) // Supprimé
			}
		}
	}

	// Étape 2 : Connexion des contours faibles aux contours forts
	for y := 1; y < bounds.Max.Y-1; y++ {
		for x := 1; x < bounds.Max.X-1; x++ {
			if output.GrayAt(x, y).Y == weak {
				// Vérifie si un pixel fort est voisin
				if isConnectedToStrong(output, x, y, strong) {
					output.SetGray(x, y, color.Gray{Y: strong})
				} else {
					output.SetGray(x, y, color.Gray{Y: 0})
				}
			}
		}
	}

	return output
}

// Vérifie si un pixel faible est connecté à un pixel fort
func isConnectedToStrong(img *image.Gray, x, y int, strong uint8) bool {
	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}
	for _, d := range directions {
		if img.GrayAt(x+d.dx, y+d.dy).Y == strong {
			return true
		}
	}
	return false
}

func ApplyCannyEdgeDetection(img *image.Gray) *image.Gray {
	// 1. Appliquer un flou Gaussien pour réduire le bruit
	kernel := GenerateGaussianKernel(5 /*5*/, 1.4 /*1.4*/)
	blurred := ApplyKernel(img, kernel)

	lowThreshold, highThreshold := ComputeDynamicThresholds(blurred, 1.5)

	// 2. Appliquer le filtre de Sobel pour obtenir les gradients
	sobelX, sobelY := GenerateSobelKernel(3)
	edges, gradientAngles := ApplySobelEdgeDetection(blurred, sobelX, sobelY)

	// 3. Suppression des non-maxima
	nms := nonMaxSuppression(*edges, gradientAngles)

	// 4. Seuil d'hystérésis
	finalEdges := hysteresisThresholding(nms, lowThreshold, highThreshold)

	return finalEdges
}
