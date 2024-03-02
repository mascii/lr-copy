package cpplan

import (
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type DirectoryMapping struct {
	shootingDate *time.Time
}

type FilePathMapping struct {
	SrcFilePath string
	DstFilePath string
}

type Plan struct {
	mapping        map[string]*DirectoryMapping
	srcDirPath     string
	dstBaseDirPath string
	separate       bool
}

func (p Plan) getDstDirPath(category string, date *time.Time) string {
	return path.Join(
		p.dstBaseDirPath,
		category,
		date.Format("2006/2006-01-02"), // Lightroom のフォルダ名の形式に合わせる
	)
}

func (p Plan) FindFilePathMapping(file fs.DirEntry) (_ *FilePathMapping, ok bool) {
	if file.IsDir() {
		return nil, false
	}
	dm, ok := p.mapping[getFileNameWithoutExt(file.Name())]
	if !ok {
		return nil, false
	}

	category := ""
	if p.separate && !isJpegFile(file.Name()) {
		category = strings.ToUpper(path.Ext(file.Name())[1:]) // "ORF", "ARW", etc.
	}

	return &FilePathMapping{
		SrcFilePath: path.Join(p.srcDirPath, file.Name()),
		DstFilePath: path.Join(p.getDstDirPath(category, dm.shootingDate), file.Name()),
	}, ok
}

func GenerateCopyPlan(files []fs.DirEntry, srcDirPath, dstBaseDirPath string, separate bool) (Plan, error) {
	mapping := make(map[string]*DirectoryMapping, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !isJpegFile(file.Name()) {
			continue
		}

		srcFullPath := path.Join(srcDirPath, file.Name())

		shootingDate, err := loadShootingDateFromExif(srcFullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load %s (%v)\n", srcFullPath, err)
			continue
		}

		fileNameWithoutExt := getFileNameWithoutExt(file.Name())
		mapping[fileNameWithoutExt] = &DirectoryMapping{
			shootingDate: shootingDate,
		}
	}

	return Plan{
		mapping:        mapping,
		srcDirPath:     srcDirPath,
		dstBaseDirPath: dstBaseDirPath,
		separate:       separate,
	}, nil
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

func loadShootingDateFromExif(filePath string) (*time.Time, error) {
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
