package main

import (
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gomager/border"
	"gomager/crop"
	"gomager/text"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
)

const (
	formatPNG  = "png"
	formatJPEG = "jpeg"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("file name not provided")
	}

	fileName := os.Args[0]

	areas, err := text.Find(fileName)
	if err != nil {
		log.Fatalf("failed to find text on image: %q", err.Error())
	}

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("failed to open Open %q, error %q", fileName, err.Error())
	}

	defer closeWrapper(file)

	img, format, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to Decode file %q, error %q", fileName, err.Error())
	}

	if len(areas) > 0 {
		img = text.Erase(img, areas)
	}

	area, found := border.Detect(img)
	if found {
		img, err = crop.Crop(img, area)
		if err != nil {
			log.Fatalf("image Crop failed: %q", err.Error())
		}
	}

	resultFileName := fmt.Sprintf("%s\result-%s", path.Dir(fileName), fileName)
	if err := saveImage(format, resultFileName, img); err != nil {
		log.Fatalf("failed to save image: %q", err.Error())
	}

}

func saveImage(format, path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to Create file %q", path)
	}

	defer closeWrapper(file)

	switch format {
	case formatPNG:
		err = png.Encode(file, img)
	case formatJPEG:
		err = jpeg.Encode(file, img, nil)
	default:
		err = errors.Errorf("unsupported image format %q", format)
	}

	return err
}

func closeWrapper(f io.Closer) {
	if err := f.Close(); err != nil {
		log.Errorf("error while file Close: %q", err.Error())
	}
}
