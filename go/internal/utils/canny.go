package utils

/*
Package utils provides tools for advanced image processing, including the implementation of the Canny edge detection algorithm.

---

### nonMaxSuppression(gradient image.Gray, angles [][]float64) *image.Gray
Performs Non-Maximum Suppression (NMS) to thin edges by suppressing non-edge gradients.

- **Parameters**:
  - gradient: A grayscale image (`image.Gray`) representing the gradient magnitudes.
  - angles: A 2D slice of angles (`[][]float64`) representing the gradient directions.

- **Returns**:
  - A new grayscale image with thinned edges (`*image.Gray`).

- **Behavior**:
  - Based on gradient angles, compares the current pixel's magnitude with neighboring pixels along the gradient direction.
  - Keeps the pixel if it is the local maximum; otherwise, suppresses it (sets it to 0).
  - Handles different gradient directions (horizontal, vertical, and diagonals) accordingly.

---

### hysteresisThresholding(img *image.Gray, lowThreshold, highThreshold float64) *image.Gray
Applies hysteresis thresholding to classify edges as strong, weak, or non-edges.

- **Parameters**:
  - img: A grayscale image (`*image.Gray`) containing edge gradients.
  - lowThreshold: The lower threshold for edge detection.
  - highThreshold: The upper threshold for edge detection.

- **Returns**:
  - A grayscale image (`*image.Gray`) with edges classified as strong or non-edges.

- **Behavior**:
  - Pixels with magnitude above `highThreshold` are classified as strong edges.
  - Pixels with magnitude between `lowThreshold` and `highThreshold` are weak edges.
  - Weak edges are only preserved if they are connected to strong edges; otherwise, they are discarded.

---

### isConnectedToStrong(img *image.Gray, x, y int, strong uint8) bool
A helper function to check if a weak edge is connected to any strong edge.

- **Parameters**:
  - img: A grayscale image (`*image.Gray`) containing edges after initial classification.
  - x, y: Coordinates of the weak edge.
  - strong: The intensity value identifying strong edges.

- **Returns**:
  - A boolean value (`true` if connected to a strong edge, `false` otherwise).

- **Behavior**:
  - Checks neighboring pixels (in an 8-connected neighborhood) to determine if any pixel is classified as a strong edge.

---

### ApplyCannyEdgeDetection(img *image.Gray) *image.Gray
The main function to apply the complete Canny edge detection pipeline to a grayscale image.

- **Parameters**:
  - img: A grayscale image (`*image.Gray`) to process.

- **Returns**:
  - A grayscale image (`*image.Gray`) with detected edges.

- **Behavior**:
  1. Applies Gaussian blurring to reduce noise using `GenerateGaussianKernel` and `ApplyKernel`.
  2. Computes gradient magnitudes and directions using Sobel filters by calling `GenerateSobelKernel` and `ApplySobelEdgeDetection`.
  3. Applies Non-Maximum Suppression (`nonMaxSuppression`) to thin the edges.
  4. Calculates dynamic thresholds using `ComputeDynamicThresholds`.
  5. Applies hysteresis thresholding (`hysteresisThresholding`) to finalize edge classification.
  6. Returns the final edge-detected image.

---

### Key Features:
- **Edge Preservation**:
  - By applying Non-Maximum Suppression, only the most prominent edges are preserved.
- **Dynamic Thresholding**:
  - Automatically computes thresholds for hysteresis, adapting to the input image's intensity distribution.
- **Noise Reduction**:
  - Employs Gaussian blurring to minimize the effect of noise on edge detection.

---

### Example Usage:
```go
package main

import (
	"image"
	"image/png"
	"os"
	"utils"
)

func main() {
	// Load a grayscale image
	file, _ := os.Open("input.png")
	defer file.Close()
	img, _, _ := image.Decode(file)
	grayImg := img.(*image.Gray)

	// Apply the Canny edge detection algorithm
	edges := utils.ApplyCannyEdgeDetection(grayImg)

	// Save the result
	outputFile, _ := os.Create("edges.png")
	defer outputFile.Close()
	png.Encode(outputFile, edges)
}
```
*/

import (
	"image"
	"image/color"
)

func nonMaxSuppression(gradient image.Gray, angles [][]float64) *image.Gray {
	bounds := gradient.Bounds()
	suppressed := image.NewGray(bounds)

	for y := 1; y < bounds.Max.Y-1; y++ {
		for x := 1; x < bounds.Max.X-1; x++ {
			angle := angles[y][x]
			mag := gradient.GrayAt(x, y).Y
			n1, n2 := uint8(0), uint8(0)

			if (angle >= -22.5 && angle <= 22.5) || (angle >= 157.5 || angle <= -157.5) {
				n1, n2 = gradient.GrayAt(x-1, y).Y, gradient.GrayAt(x+1, y).Y
			} else if (angle > 22.5 && angle <= 67.5) || (angle < -112.5 && angle >= -157.5) {
				n1, n2 = gradient.GrayAt(x-1, y-1).Y, gradient.GrayAt(x+1, y+1).Y
			} else if (angle > 67.5 && angle <= 112.5) || (angle < -67.5 && angle >= -112.5) {
				n1, n2 = gradient.GrayAt(x, y-1).Y, gradient.GrayAt(x, y+1).Y
			} else {
				n1, n2 = gradient.GrayAt(x-1, y+1).Y, gradient.GrayAt(x+1, y-1).Y
			}

			if mag >= n1 && mag >= n2 {
				suppressed.SetGray(x, y, color.Gray{Y: mag})
			} else {
				suppressed.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}
	return suppressed
}

func hysteresisThresholding(img *image.Gray, lowThreshold, highThreshold float64) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)

	strong := uint8(255)
	weak := uint8(75)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.GrayAt(x, y).Y
			if float64(pixel) >= highThreshold {
				output.SetGray(x, y, color.Gray{Y: strong})
			} else if float64(pixel) >= lowThreshold {
				output.SetGray(x, y, color.Gray{Y: weak})
			} else {
				output.SetGray(x, y, color.Gray{Y: 0})
			}
		}
	}

	for y := 1; y < bounds.Max.Y-1; y++ {
		for x := 1; x < bounds.Max.X-1; x++ {
			if output.GrayAt(x, y).Y == weak {
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
	kernel := GenerateGaussianKernel(5, 1.4)
	blurred := ApplyKernel(img, kernel)

	lowThreshold, highThreshold := ComputeDynamicThresholds(blurred, 1.5)

	sobelX, sobelY := GenerateSobelKernel(3)
	edges, gradientAngles := ApplySobelEdgeDetection(blurred, sobelX, sobelY)

	nms := nonMaxSuppression(*edges, gradientAngles)

	finalEdges := hysteresisThresholding(nms, lowThreshold, highThreshold)

	return finalEdges
}
