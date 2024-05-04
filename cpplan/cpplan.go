package cpplan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mascii/lr-copy/extractor"
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
	srcDirPath             string
	dstBaseDirPath         string
	separate               bool
	shootingDateExtractors map[string]func(filePath string) (*time.Time, error)
}

func NewGenerateCopyPlanConfig(srcDirPath, dstBaseDirPath string, separate bool) generateCopyPlanConfig {
	return generateCopyPlanConfig{
		srcDirPath:     srcDirPath,
		dstBaseDirPath: dstBaseDirPath,
		separate:       separate,
		shootingDateExtractors: map[string]func(filePath string) (*time.Time, error){
			"JPG":  extractor.LoadShootingDateFromJpeg,
			"JPEG": extractor.LoadShootingDateFromJpeg,
		},
	}
}

func GenerateCopyPlan[T DirEntrySubset](files []T, cfg generateCopyPlanConfig) []*FilePathMapping {
	mapping := make(map[string]*time.Time, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		extractor, ok := cfg.shootingDateExtractors[getExtByFileName(file.Name())]
		if !ok {
			continue
		}

		srcFullPath := filepath.Join(cfg.srcDirPath, file.Name())

		date, err := extractor(srcFullPath)
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
		if cfg.separate {
			if _, ok = cfg.shootingDateExtractors[getExtByFileName(file.Name())]; !ok {
				category = getExtByFileName(file.Name()) // "ORF", "ARW", etc.
			}
		}

		plan = append(plan, &FilePathMapping{
			SrcFilePath: filepath.Join(cfg.srcDirPath, file.Name()),
			DstFilePath: filepath.Join(cfg.dstBaseDirPath, category, dateToLightroomFormat(date), file.Name()),
		})
	}

	return plan
}

func getExtByFileName(fileName string) string {
	return strings.ToUpper(filepath.Ext(fileName)[1:]) // "JPG", "JPEG", etc.
}

func getFileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func dateToLightroomFormat(d *time.Time) string {
	return d.Format("2006/2006-01-02") // Lightroom のフォルダ名の形式に合わせる
}
