package main

import (
	"ELP-project/internal/imageUtils"
	"ELP-project/internal/utils"
	"fmt"
	"log"
)

// Main Canny filter pipeline.
func main() {
	// Input/output paths
	inputPath := "input.png"
	outputPath := "output.png"

	// Load image
	img, format, err := imageUtils.LoadImage(inputPath)
	if err != nil {
		log.Fatalf("Failed to load input image: %v", err)
	}

	// Convert to grayscale
	grayImg := imageUtils.Grayscale(img)

	// Customize Gaussian kernel
	kernelSize := 15
	kernelSigma := float64(kernelSize) / 6
	gaussianKernel := utils.GenerateGaussianKernel(kernelSize, kernelSigma)

	// Apply Gaussian blur
	blurredImg := utils.ApplyKernel(grayImg, gaussianKernel)

	// Apply Sobel edge detection
	edges, _ := utils.ApplySobelEdgeDetection(blurredImg)

	// Save the result
	err = imageUtils.SaveImage(edges, outputPath, format)
	if err != nil {
		log.Fatalf("Failed to save output image: %v", err)
	}

	fmt.Println("Canny filter applied and output saved to", outputPath)
}
