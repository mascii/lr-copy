package extractor

import (
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func LoadShootingDateFromJpeg(filePath string) (*time.Time, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return nil, err
	}

	t, err := x.DateTime()
	if err != nil {
		return nil, err
	}

	return &t, nil
}
