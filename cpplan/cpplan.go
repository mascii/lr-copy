package cpplan

import (
	"io/fs"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/mascii/lr-copy/extractor"
)

// DirEntrySubset is a subset of fs.DirEntry.
type DirEntrySubset interface {
	Name() string
}

type FilePathMapping struct {
	SrcFilePath string
	DstFilePath string
}

type generateCopyPlanConfig struct {
	srcDirPath             string
	dstBaseDirPath         string
	separate               bool
	fallback               bool
	shootingDateExtractors map[string]func(filePath string) (*time.Time, error)
}

func NewGenerateCopyPlanConfig(srcDirPath, dstBaseDirPath string, separate, fallback bool) generateCopyPlanConfig {
	return generateCopyPlanConfig{
		srcDirPath:     srcDirPath,
		dstBaseDirPath: dstBaseDirPath,
		separate:       separate,
		fallback:       fallback,
		shootingDateExtractors: map[string]func(filePath string) (*time.Time, error){
			"JPG":  extractor.LoadShootingDateFromJpeg,
			"JPEG": extractor.LoadShootingDateFromJpeg,
			"HEIC": extractor.LoadShootingDateFromHeic,
			"MOV":  extractor.LoadShootingDateFromMov,
			"PNG":  extractor.LoadShootingDateFromPng,
		},
	}
}

func GenerateCopyPlan[T DirEntrySubset](files []T, cfg generateCopyPlanConfig) []*FilePathMapping {
	mapping := make(map[string]*time.Time, len(files))
	for _, file := range files {
		extractor, ok := cfg.shootingDateExtractors[getExtByFileName(file.Name())]
		if !ok {
			continue
		}

		srcFullPath := filepath.Join(cfg.srcDirPath, file.Name())
		fileNameWithoutExt := getFileNameWithoutExt(file.Name())

		if t, err := extractor(srcFullPath); err == nil {
			mapping[fileNameWithoutExt] = t
		} else if file, ok := any(file).(fs.DirEntry); cfg.fallback && ok {
			info, err2 := file.Info()
			if err2 != nil {
				log.Fatal(err2)
			}
			mt := info.ModTime()
			mapping[fileNameWithoutExt] = &mt
			log.Printf("Failed to load the shooting date in %s (%v), mod time will be used as a fallback (%v)\n", srcFullPath, err, mt)
		} else {
			log.Printf("Failed to load the shooting date in %s (%v)\n", srcFullPath, err)
		}
	}

	plan := make([]*FilePathMapping, 0, len(files))
	for _, file := range files {
		date, ok := mapping[getFileNameWithoutExt(file.Name())]
		if !ok {
			continue
		}

		category := ""
		if cfg.separate {
			ext := getExtByFileName(file.Name())
			if _, ok = cfg.shootingDateExtractors[ext]; !ok {
				category = ext // "ORF", "ARW", etc.
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
	ext := filepath.Ext(fileName)
	if len(ext) == 0 {
		return ""
	}
	return strings.ToUpper(ext[1:]) // "JPG", "JPEG", etc.
}

func getFileNameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func dateToLightroomFormat(d *time.Time) string {
	return d.Format("2006/2006-01-02") // Lightroom のフォルダ名の形式に合わせる
}
