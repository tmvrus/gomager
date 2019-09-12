package crop

import (
	"image"

	"github.com/oliamb/cutter"
	"github.com/pkg/errors"
)

func Crop(src image.Image, area image.Rectangle) (image.Image, error) {
	cropped, err := cutter.Crop(src, cutter.Config{
		Width:  area.Max.X - area.Min.X,
		Height: area.Max.Y - area.Min.Y,
		Anchor: image.Point{
			X: area.Min.X,
			Y: area.Min.Y,
		},
		Options: cutter.Copy,
	})

	if err != nil {
		return src, errors.Wrap(err, "failed to Crop")
	}

	return cropped, nil
}
