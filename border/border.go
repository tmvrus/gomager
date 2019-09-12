package border

import (
	"image"
	"image/color"

	"github.com/anthonynsimon/bild/effect"
)

const (
	almostBlackThreshold = 20000 //pretty subjective but better than nothing
)

func Detect(src image.Image) (area image.Rectangle, detected bool) {

	src = effect.Grayscale(src)
	bounds := src.Bounds()
	area.Min.X = bounds.Min.X
	area.Max.X = bounds.Max.X
	area.Min.Y = bounds.Min.Y
	area.Max.Y = bounds.Max.Y

	// up to down
	for y := bounds.Min.Y; y <= bounds.Max.Y-1; y++ {
		if !IsBorderRow(bounds.Min.X, bounds.Max.X-1, y, src) {
			break
		}

		area.Min.Y = y
		detected = true
	}

	// down to up
	for y := bounds.Max.Y - 1; y >= bounds.Min.Y; y-- {
		if !IsBorderRow(bounds.Min.X, bounds.Max.X-1, y, src) {
			break
		}

		area.Max.Y = y
		detected = true
	}

	// looks like image is whole black
	if area.Min.Y == area.Max.Y {
		detected = false
	}

	return
}

func IsBorderRow(from, to, row int, src image.Image) bool {
	for x := from; x <= to; x++ {
		if !isBorderPixel(almostBlackThreshold, src.At(x, row)) {
			return false
		}
	}

	return true
}

func isBorderPixel(threshold uint32, c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r <= threshold && g <= threshold && b <= threshold
}
