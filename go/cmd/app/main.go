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
	inputPath := "./go/input.jpg"
	outputPath := "output.jpg"

	// Load image
	img, format, err := imageUtils.LoadImage(inputPath)
	if err != nil {
		log.Fatalf("Failed to load input image: %v", err)
	}

	// Convert to grayscale
	grayImg := imageUtils.Grayscale(img)

	edges := utils.ApplyCannyEdgeDetection(grayImg, 50, 150)

	// Save the result
	err = imageUtils.SaveImage(edges, outputPath, format)
	if err != nil {
		log.Fatalf("Failed to save output image: %v", err)
	}

	fmt.Println("Canny filter applied and output saved to", outputPath)
}
