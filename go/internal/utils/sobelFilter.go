package utils

/*
Package utils provides tools for image processing, including dynamic Sobel kernel generation,
edge detection, and threshold calculations.

---

### GenerateSobelKernel(size int) ([][]float64, [][]float64)
Dynamically generates a Sobel kernel of variable size for edge detection.

- **Parameters**:
  - size: The size of the kernel (must be odd). Supports standard 3x3 and 5x5 kernels with predefined values or creates larger kernels dynamically.
- **Returns**:
  - Two Sobel kernels: one for the X-gradient (`[][]float64`) and one for the Y-gradient (`[][]float64`).
- **Behavior**:
  - For size 3 or 5, predefined kernels are returned.
  - For larger sizes, Gaussian-like approximation is applied to generate Sobel derivatives.
  - Kernels are normalized to ensure their absolute sum equals 1.
- **Panics**:
  - If `size` is even, since Sobel kernels require odd dimensions.

---

### normalizeKernel(kernel [][]float64)
Normalizes a Sobel kernel so that the sum of its absolute values equals 1.

- **Parameters**:
  - kernel: A 2D slice representing the kernel.
- **Behavior**:
  - Computes the sum of the absolute values of all elements.
  - Scales each element by dividing by the total sum.

---

### ComputeDynamicThresholds(img *image.Gray, alpha float64) (float64, float64)
Calculates dynamic thresholds for edge detection based on image gradients.

- **Parameters**:
  - img: A grayscale image (`*image.Gray`).
  - alpha: A multiplier for the high threshold.
- **Returns**:
  - lowThreshold: The lower bound for edge detection.
  - highThreshold: The upper bound for edge detection.
- **Behavior**:
  - Applies a 5x5 Sobel filter to compute the gradient magnitude of the image.
  - Calculates the average gradient magnitude and sets `highThreshold` as `alpha * meanGradient`.
  - `lowThreshold` is set to 40% of `highThreshold`.

---

### ApplySobelEdgeDetection(img *image.Gray, kernelX, kernelY [][]float64) (*image.Gray, [][]float64)
Applies a Sobel edge detection filter to a grayscale image.

- **Parameters**:
  - img: A grayscale image (`*image.Gray`).
  - kernelX: A Sobel kernel for detecting X-gradients (`[][]float64`).
  - kernelY: A Sobel kernel for detecting Y-gradients (`[][]float64`).
- **Returns**:
  - output: A new grayscale image (`*image.Gray`) representing the magnitude of the gradient.
  - gradientAngles: A 2D slice of gradient angles (`[][]float64`), where each value corresponds to the angle of the gradient at a pixel.
- **Behavior**:
  - Convolves the input image with the provided Sobel kernels in both X and Y directions.
  - Computes the gradient magnitude (`sqrt(gx^2 + gy^2)`) and angle (`atan2(gy, gx)`) for each pixel.
  - Clamps the gradient magnitude to a maximum value of 255 for 8-bit images.
  - Returns the filtered image and gradient orientations.

---

### Key Features:
- **Dynamic Kernel Generation**:
  - Easily customize Sobel kernels to adapt to specific image resolutions and requirements.
- **Edge Detection**:
  - Apply Sobel filters to highlight edges of various orientations.
- **Thresholding**:
  - Dynamically computed thresholds allow robust edge detection independent of image intensity.

---

### Example Usage:
```go
package main

import (
	"image"
	"image/color"
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

	// Generate Sobel kernels and apply edge detection
	sobelX, sobelY := utils.GenerateSobelKernel(5)
	edges, _ := utils.ApplySobelEdgeDetection(grayImg, sobelX, sobelY)

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
	"math"
)

func GenerateSobelKernel(size int) ([][]float64, [][]float64) {
	if size%2 == 0 {
		panic("Sobel kernel size must be odd")
	}

	if size == 3 {
		return [][]float64{
				{-1, 0, 1},
				{-2, 0, 2},
				{-1, 0, 1},
			}, [][]float64{
				{-1, -2, -1},
				{0, 0, 0},
				{1, 2, 1},
			}
	}

	if size == 5 {
		return [][]float64{
				{-2, -1, 0, 1, 2},
				{-3, -2, 0, 2, 3},
				{-4, -3, 0, 3, 4},
				{-3, -2, 0, 2, 3},
				{-2, -1, 0, 1, 2},
			}, [][]float64{
				{-2, -2, -4, -2, -2},
				{-1, -1, -2, -1, -1},
				{0, 0, 0, 0, 0},
				{1, 1, 2, 1, 1},
				{2, 2, 4, 2, 2},
			}
	}

	kernelX := make([][]float64, size)
	kernelY := make([][]float64, size)
	radius := size / 2
	sigma := float64(size) / 3

	for i := 0; i < size; i++ {
		kernelX[i] = make([]float64, size)
		kernelY[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			x, y := float64(j-radius), float64(i-radius)

			kernelX[i][j] = -x * math.Exp(-(x*x+y*y)/(2*sigma*sigma))
			kernelY[i][j] = -y * math.Exp(-(x*x+y*y)/(2*sigma*sigma))
		}
	}

	normalizeKernel(kernelX)
	normalizeKernel(kernelY)

	return kernelX, kernelY
}

func normalizeKernel(kernel [][]float64) {
	sum := 0.0
	for _, row := range kernel {
		for _, val := range row {
			sum += math.Abs(val)
		}
	}

	if sum != 0 {
		for i := range kernel {
			for j := range kernel[i] {
				kernel[i][j] /= sum
			}
		}
	}
}

func ComputeDynamicThresholds(img *image.Gray, alpha float64) (float64, float64) {
	bounds := img.Bounds()
	totalGradient := 0.0
	count := 0

	sobelX, sobelY := GenerateSobelKernel(5)
	gradient, _ := ApplySobelEdgeDetection(img, sobelX, sobelY)

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			totalGradient += float64(gradient.GrayAt(x, y).Y)
			count++
		}
	}

	meanGradient := totalGradient / float64(count)

	highThreshold := alpha * meanGradient
	lowThreshold := 0.4 * highThreshold

	return lowThreshold, highThreshold
}

func ApplySobelEdgeDetection(img *image.Gray, kernelX, kernelY [][]float64) (*image.Gray, [][]float64) {
	bounds := img.Bounds()
	output := image.NewGray(bounds)
	gradientAngles := make([][]float64, bounds.Max.Y)
	radius := len(kernelX) / 2

	for i := range gradientAngles {
		gradientAngles[i] = make([]float64, bounds.Max.X)
	}

	for y := bounds.Min.Y + radius; y < bounds.Max.Y-radius; y++ {
		for x := bounds.Min.X + radius; x < bounds.Max.X-radius; x++ {
			var gx, gy float64

			for ky := -radius; ky <= radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					px := x + kx
					py := y + ky

					gray := float64(img.GrayAt(px, py).Y)
					gx += gray * kernelX[ky+radius][kx+radius]
					gy += gray * kernelY[ky+radius][kx+radius]
				}
			}

			magnitude := math.Sqrt(gx*gx + gy*gy)
			angle := math.Atan2(gy, gx) * (180 / math.Pi)

			output.SetGray(x, y, color.Gray{Y: uint8(math.Min(magnitude, 255))})
			gradientAngles[y][x] = angle
		}
	}

	return output, gradientAngles
}
