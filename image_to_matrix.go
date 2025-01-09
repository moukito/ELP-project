package main

import (
	"fmt"
	"image"
	_ "image/jpeg" // Pour décoder les images JPEG
	_ "image/png"  // Pour décoder les images PNG
	"os"
)

func image_to_matrix(imagePath string) ([][][4]uint8, error) {
	// Ouvrir le fichier image
	file, err := os.Open(imagePath) // Remplacez "image.jpg" par le chemin de votre image
	if err != nil {
		return nil, fmt.Errorf("Erreur lors de l'ouverture de l'image : %w", err)
	}
	defer file.Close()

	// Décoder l'image
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("erreur lors du décodage de l'image : %w", err)
	}

	// Obtenir la taille de l'image
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// Créer une matrice 2D pour les pixels
	matrix := make([][][4]uint8, height)
	for y := 0; y < height; y++ {
		row := make([][4]uint8, width)
		for x := 0; x < width; x++ {
			// Extraire la couleur du pixel
			r, g, b, a := img.At(x, y).RGBA()
			row[x] = [4]uint8{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
		}
		matrix[y] = row
	}

	return matrix, nil
}

func main() {
	imagePath := "image.png"
	matrix, err := image_to_matrix(imagePath)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	for x := 0; x < 5; x++ {
		fmt.Printf("Pixel (%d, %d): %v\n", x, 0, matrix[0][x])
	}
}
