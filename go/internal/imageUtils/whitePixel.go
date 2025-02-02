package imageUtils

import "image"

func IsWhite(img *image.Gray, x, y int) bool {
	return img.GrayAt(x, y).Y > 128
}
