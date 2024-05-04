package extractor

import (
	"bytes"
	"time"

	heicexif "github.com/dsoprea/go-heic-exif-extractor"
	"github.com/rwcarlsen/goexif/exif"
)

func LoadShootingDateFromHeic(filePath string) (*time.Time, error) {
	mc, err := heicexif.NewHeicExifMediaParser().ParseFile(filePath)
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
