package extractor

import (
	"os"
	"time"

	"github.com/mascii/lr-copy/movmeta"
)

func LoadShootingDateFromMov(filePath string) (*time.Time, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	t, err := movmeta.GetVideoCreationTimeMetadata(f)

	return &t, err
}
