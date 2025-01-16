package main

import (
	imageUtils2 "ELP-project/internal/imageUtils"
	utils2 "ELP-project/internal/utils"
	"ELP-project/server"
	"fmt"
	"log"
)

// Main Canny filter pipeline.
func main() {
	// Input/output paths
	inputPath := "input.png"
	outputPath := "output.png"

	// Load image
	img, format, err := imageUtils2.LoadImage(inputPath)
	if err != nil {
		log.Fatalf("Failed to load input image: %v", err)
	}

	// Convert to grayscale
	grayImg := imageUtils2.Grayscale(img)

	// Customize Gaussian kernel
	kernelSize := 5
	kernelSigma := 1.4
	gaussianKernel := utils2.GenerateGaussianKernel(kernelSize, kernelSigma)

	// Apply Gaussian blur
	blurredImg := utils2.ApplyKernel(grayImg, gaussianKernel)

	// Apply Sobel edge detection
	edges, _ := utils2.ApplySobelEdgeDetection(blurredImg)

	// Save the result
	err = imageUtils2.SaveImage(edges, outputPath, format)
	if err != nil {
		log.Fatalf("Failed to save output image: %v", err)
	}

	fmt.Println("Canny filter applied and output saved to", outputPath)

	server := server.New(&server.Config{
		Host: "localhost",
		Port: "3333",
	})
	server.Run()
}
