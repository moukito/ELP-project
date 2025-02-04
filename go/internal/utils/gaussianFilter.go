package utils

/*
Package utils provides tools for image processing, including Gaussian kernel generation and its application for image filtering.

---

### GenerateGaussianKernel(size int, sigma float64) [][]float64
Generates a 2D Gaussian kernel of the specified size. Gaussian kernels are commonly used for operations such as blurring in image processing.

- **Parameters**:
  - `size` (int): The size of the kernel. Must be odd; otherwise, the function will panic.
  - `sigma` (float64): The standard deviation of the Gaussian distribution, controlling the spread of the kernel.
- **Returns**:
  - `[][]float64`: A 2D slice representing the Gaussian kernel, normalized so that the sum of its elements equals 1.
- **Behavior**:
  - Creates a `size x size` matrix based on the Gaussian formula. Each element represents the corresponding weight.
  - Normalizes the kernel so that the total sum equals 1, ensuring it's suitable for use in convolution.

#### Panics:
- If `size` is even, as Gaussian kernels require odd dimensions.

#### Example Usage:
```go
kernel := GenerateGaussianKernel(5, 1.0) // Generates a 5x5 Gaussian kernel with a sigma of 1.0
```

---

### ApplyKernel(img *image.Gray, kernel [][]float64) *image.Gray
Applies a 2D convolution using a specified kernel (e.g., a Gaussian kernel) to a grayscale image.

- **Parameters**:
  - `img` (*image.Gray): The grayscale image onto which the kernel is applied.
  - `kernel` ([][]float64): The kernel to use for convolution (like a Gaussian kernel).
- **Returns**:
  - `*image.Gray`: A new grayscale image where the kernel has been applied.

#### Behavior:
- Iterates over the image pixels and calculates a weighted sum for each pixel based on the kernel.
- Accounts for image boundaries by excluding out-of-bounds pixels during convolution.
- Creates and returns a new grayscale image resulting from the convolution.

#### Example Usage:
```go
img := ... // Load or create a grayscale image
kernel := GenerateGaussianKernel(5, 1.0) // Create a 5x5 Gaussian kernel
blurredImg := ApplyKernel(img, kernel)   // Apply the kernel to blur the image
```

---

### Key Features:
- **Dynamic Gaussian Kernel Generation**:
  - Easily create Gaussian kernels of various sizes to match specific filter requirements.
- **Image Filtering**:
  - Apply Gaussian filters for blurring, noise reduction, or pre-processing steps before edge detection.
- **Efficient Convolution**:
  - Handles convolution with arbitrary kernels, making this tool flexible for tasks beyond Gaussian smoothing.

---

### Example Workflow:
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

	// Generate a Gaussian kernel and apply it to the image
	kernel := utils.GenerateGaussianKernel(5, 1.0) // 5x5 kernel with sigma 1.0
	blurredImg := utils.ApplyKernel(grayImg, kernel)

	// Save the blurred image
	outputFile, _ := os.Create("blurred.png")
	defer outputFile.Close()
	png.Encode(outputFile, blurredImg)
}
```
*/

import (
	"image"
	"image/color"
	"math"
)

func GenerateGaussianKernel(size int, sigma float64) [][]float64 {
	if size%2 == 0 {
		panic("Gaussian kernel size must be odd")
	}

	kernel := make([][]float64, size)
	sum := 0.0
	radius := size / 2

	for i := 0; i < size; i++ {
		kernel[i] = make([]float64, size)
		for j := 0; j < size; j++ {
			x, y := float64(i-radius), float64(j-radius)
			kernel[i][j] = (1 / (2 * math.Pi * sigma * sigma)) * math.Exp(-(x*x+y*y)/(2*sigma*sigma))
			sum += kernel[i][j]
		}
	}

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			kernel[i][j] /= sum
		}
	}

	return kernel
}

func ApplyKernel(img *image.Gray, kernel [][]float64) *image.Gray {
	bounds := img.Bounds()
	output := image.NewGray(bounds)
	radius := len(kernel) / 2

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var sum float64
			var weightSum float64

			for ky := -radius; ky <= radius; ky++ {
				for kx := -radius; kx <= radius; kx++ {
					pixelX := x + kx
					pixelY := y + ky
					if pixelX >= bounds.Min.X && pixelX < bounds.Max.X && pixelY >= bounds.Min.Y && pixelY < bounds.Max.Y {
						gray := float64(img.GrayAt(pixelX, pixelY).Y)
						sum += gray * kernel[ky+radius][kx+radius]
						weightSum += kernel[ky+radius][kx+radius]
					}
				}
			}

			output.SetGray(x, y, color.Gray{Y: uint8(sum / weightSum)})
		}
	}

	return output
}
