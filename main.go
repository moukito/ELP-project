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

	kernel := utils.GenerateGaussianKernel(2, 1.0)
	fmt.Printf("Pixel (%d, %d): %v\n", 0, 0, kernel[0][0])
	fmt.Printf("Pixel (%d, %d): %v\n", 0, 1, kernel[0][1])
	fmt.Printf("Pixel (%d, %d): %v\n", 1, 0, kernel[1][0])
	fmt.Printf("Pixel (%d, %d): %v\n", 1, 1, kernel[1][1])
	result := utils.ApplyGaussianFilter(matrix, kernel)
	fmt.Printf("Pixel (%d, %d): %v\n", 5, 0, result[5][10])

	utils.SaveImage(utils.MatrixToImage(result), "output.png")
}
