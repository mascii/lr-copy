package extractor

import (
	"bytes"
	"os"
	"time"

	heicexif "github.com/dsoprea/go-heic-exif-extractor"
	"github.com/rwcarlsen/goexif/exif"
)

func LoadShootingDateFromHeic(filePath string) (*time.Time, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	mc, err := heicexif.NewHeicExifMediaParser().Parse(f, 0)
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
