package imageUtils

import (
	"image"
	"os"
)

// LoadImage loads an image from a given file path.
func LoadImage(filePath string) (image.Image, string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}
