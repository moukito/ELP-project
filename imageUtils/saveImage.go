package imageUtils

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

// SaveImage saves an image to a file in the specified format.
func SaveImage(img image.Image, filePath string, format string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	switch strings.ToLower(format) {
	case "jpg", "jpeg":
		return jpeg.Encode(file, img, nil)
	case "png":
		return png.Encode(file, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
