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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}
