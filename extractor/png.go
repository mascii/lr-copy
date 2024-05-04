package extractor

import (
	"bytes"
	"time"

	pngstructure "github.com/dsoprea/go-png-image-structure"
	"github.com/rwcarlsen/goexif/exif"
)

func LoadShootingDateFromPng(filePath string) (*time.Time, error) {
	mc, err := pngstructure.NewPngMediaParser().ParseFile(filePath)
	if err != nil {
		return nil, err
	}

	_, rawExif, err := mc.Exif()
	if err != nil {
		return nil, err
	}

	x, err := exif.Decode(bytes.NewReader(rawExif))
	if err != nil {
		return nil, err
	}

	t, err := x.DateTime()
	if err != nil {
		return nil, err
	}

	return &t, nil
}
