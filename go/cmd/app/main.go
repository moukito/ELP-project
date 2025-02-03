package main

import (
	"ELP-project/internal/geometry"
	"ELP-project/internal/imageUtils"
	"ELP-project/internal/utils"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
)

// Main Canny filter pipeline.
func main() {
	// Input/output paths
	inputPath := "./go/image2.jpg"
	outputPath := "output.jpg"

	// Load image
	img, format, err := imageUtils.LoadImage(inputPath)
	if err != nil {
		log.Fatalf("Failed to load input image: %v", err)
	}

	// Convert to grayscale
	grayImg := imageUtils.Grayscale(img)

	edges := utils.ApplyCannyEdgeDetection(grayImg)

	// Save edges to a file for visualization
	edgesFile, err := os.Create("edges.jpg")
	if err != nil {
		log.Fatalf("Failed to create edges file: %v", err)
	}
	defer edgesFile.Close()
	jpeg.Encode(edgesFile, edges, nil)
	fmt.Println("Edges saved to edges.jpg")

	contours := utils.FindContoursBFS(edges)
	contourComplet := utils.FindQuadrilateral(contours)
	fmt.Println(len(contourComplet.Contour))

	/*if len(contourA4) != 4 {
		fmt.Println("No contour found.")
		return
	}*/

	// Sauvegarder l'image avec le contour dessiné
	contourImg := utils.DrawContour(img, contourComplet.Contour)
	outFile, err := os.Create("contour_detected.jpg")
	if err != nil {
		fmt.Println("Erreur de création :", err)
		return
	}
	defer outFile.Close()
	jpeg.Encode(outFile, contourImg, nil)
	fmt.Println("Image avec contour sauvegardée dans contour_detected.jpg")

	corner1 := geometry.Point{
		X: 1500,
		Y: 2000,
	}
	corner2 := geometry.Point{
		X: 1500,
		Y: 2000,
	}
	corner3 := geometry.Point{
		X: 1500,
		Y: 2000,
	}
	corner4 := geometry.Point{
		X: 1500,
		Y: 2000,
	}
	for contour := range contourComplet.Contour {
		if contourComplet.Contour[contour].X < corner1.X {
			corner1.X = contourComplet.Contour[contour].X
		}
		if contourComplet.Contour[contour].Y < corner1.Y {
			corner1.Y = contourComplet.Contour[contour].Y
		}
		if contourComplet.Contour[contour].X > corner2.X {
			corner2.X = contourComplet.Contour[contour].X
		}
		if contourComplet.Contour[contour].Y < corner2.Y {
			corner2.Y = contourComplet.Contour[contour].Y
		}
		if contourComplet.Contour[contour].X < corner3.X {
			corner3.X = contourComplet.Contour[contour].X
		}
		if contourComplet.Contour[contour].Y > corner3.Y {
			corner3.Y = contourComplet.Contour[contour].Y
		}
		if contourComplet.Contour[contour].X > corner4.X {
			corner4.X = contourComplet.Contour[contour].X
		}
		if contourComplet.Contour[contour].Y > corner4.Y {
			corner4.Y = contourComplet.Contour[contour].Y
		}
	}
	contourA4 := geometry.Contour{
		corner1,
		corner2,
		corner4,
		corner3,
	}
	fmt.Printf("Contour A4 points: %+v\n", contourA4)

	// Extraire uniquement la région intérieure du contour
	extractedRegion := utils.ExtractRegion(img, contourA4)
	outFile, err = os.Create("extracted_region.jpg")
	if err != nil {
		fmt.Println("Erreur de création :", err)
		return
	}
	defer outFile.Close()
	jpeg.Encode(outFile, extractedRegion, nil)
	fmt.Println("Image extraite sauvegardée dans extracted_region.jpg")
	/*
		// Définir un rectangle A4 cible
		targetSize := [4]utils.Point2f{
			{0, 0},
			{3072, 0},
			{3072, 4096},
			{0, 4096},
		}

		// Calcul de l'homographie
		homography := utils.ComputeHomographyMatrix([4]utils.Point2f{
			{float64(contourA4[0].X), float64(contourA4[0].Y)},
			{float64(contourA4[1].X), float64(contourA4[1].Y)},
			{float64(contourA4[2].X), float64(contourA4[2].Y)},
			{float64(contourA4[3].X), float64(contourA4[3].Y)},
		}, targetSize)

		// Appliquer la transformation
		warped := utils.ApplyPerspectiveTransform(img, homography, 3072, 4096)*/

	rect := image.Rect(contourA4[0].X, contourA4[0].Y, contourA4[2].X, contourA4[2].Y)
	warped := image.NewRGBA(rect)
	draw.Draw(warped, rect, img, image.Pt(contourA4[0].X, contourA4[0].Y), draw.Src)

	// Save the result
	err = imageUtils.SaveImage(warped, outputPath, format)
	if err != nil {
		log.Fatalf("Failed to save output image: %v", err)
	}

	fmt.Println("Canny filter applied and output saved to", outputPath)

	//img2 := main2.MaskOutsideCorners(edgesBackup, 128, 0.5)

	//imageUtils.SaveImage(img2, "image_with_corner.jpg", format)

}
