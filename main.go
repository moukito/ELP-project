package main

import (
	"ELP-project/utils"
	"fmt"
)

func main() {
	imagePath := "image.png"
	matrix, err := utils.ImageToMatrix(imagePath)
	if err != nil {
		fmt.Println("Erreur :", err)
		return
	}
	for x := 0; x < 5; x++ {
		fmt.Printf("Pixel (%d, %d): %v\n", x, 0, matrix[0][x])
	}
}
