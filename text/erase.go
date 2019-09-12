package text

import (
	"image"
	"image/draw"
)

func Erase(src image.Image, areas []image.Rectangle) image.Image {
	b := src.Bounds()
	newImg := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newImg, b, src, b.Min, draw.Src)

	for _, area := range areas {
		for x := area.Min.X; x < area.Max.X; x++ {
			for y := area.Min.Y; y < area.Max.Y; y++ {
				c := newImg.At(x-1, y-1)
				newImg.Set(x, y, c)
			}
		}
	}

	return newImg
}
