package cpplan

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

// DirEntrySubset is a subset of fs.DirEntry.
type DirEntrySubset interface {
	Name() string
	IsDir() bool
}

type FilePathMapping struct {
	SrcFilePath string
	DstFilePath string
}

type generateCopyPlanConfig struct {
	srcDirPath               string
	dstBaseDirPath           string
	separate                 bool
	loadShootingDateFromExif func(filePath string) (*time.Time, error)
}

func NewGenerateCopyPlanConfig(srcDirPath, dstBaseDirPath string, separate bool) generateCopyPlanConfig {
	return generateCopyPlanConfig{
		srcDirPath:     srcDirPath,
		dstBaseDirPath: dstBaseDirPath,
		separate:       separate,
		loadShootingDateFromExif: func(filePath string) (*time.Time, error) {
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
		},
	}
}

func GenerateCopyPlan[T DirEntrySubset](files []T, cfg generateCopyPlanConfig) []*FilePathMapping {
	mapping := make(map[string]*time.Time, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !isJpegFile(file.Name()) {
			continue
		}

		srcFullPath := path.Join(cfg.srcDirPath, file.Name())

		date, err := cfg.loadShootingDateFromExif(srcFullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load %s (%v)\n", srcFullPath, err)
			continue
		}

		fileNameWithoutExt := getFileNameWithoutExt(file.Name())
		mapping[fileNameWithoutExt] = date
	}

	plan := make([]*FilePathMapping, 0, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		date, ok := mapping[getFileNameWithoutExt(file.Name())]
		if !ok {
			continue
		}

		category := ""
		if cfg.separate && !isJpegFile(file.Name()) {
			category = strings.ToUpper(path.Ext(file.Name())[1:]) // "ORF", "ARW", etc.
		}

		plan = append(plan, &FilePathMapping{
			SrcFilePath: path.Join(cfg.srcDirPath, file.Name()),
			DstFilePath: path.Join(cfg.dstBaseDirPath, category, dateToLightroomFormat(date), file.Name()),
		})
	}

	return plan
}

func isJpegFile(fileName string) bool {
	ext := path.Ext(fileName)
	switch strings.ToLower(ext) {
	case ".jpg", ".jpeg":
		return true
	default:
		return false
	}
}

func getFileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, path.Ext(fileName))
}

func dateToLightroomFormat(d *time.Time) string {
	return d.Format("2006/2006-01-02") // Lightroom のフォルダ名の形式に合わせる
}
