package utils

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// MatrixToImage convertit une matrice de pixels RGBA en une image de type image.RGBA.
func MatrixToImage(matrix [][][4]uint8) *image.RGBA {
	height := len(matrix)
	width := len(matrix[0])

	// Créer une nouvelle image RGBA
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Parcourir chaque pixel et définir sa couleur dans l'image
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := matrix[y][x]
			col := color.RGBA{
				R: pixel[0],
				G: pixel[1],
				B: pixel[2],
				A: pixel[3],
			}
			img.Set(x, y, col)
		}
	}

	return img
}

// SaveImage enregistre une image au format PNG.
func SaveImage(img *image.RGBA, outputPath string) error {
	// Créer un fichier pour écrire l'image
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encoder l'image en PNG
	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
