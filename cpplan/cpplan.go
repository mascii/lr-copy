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

// DirEntrySubset is a subset of fs.DirEntry.
type DirEntrySubset interface {
	Name() string
	IsDir() bool
}

type FilePathMapping struct {
	SrcFilePath string
	DstFilePath string
}

type Plan struct {
	mapping        map[string]*time.Time // mapping[ファイル名(拡張子なし)] = 撮影日時
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

func (p Plan) HasNoFilesToCopy() bool {
	return len(p.mapping) == 0
}

func (p Plan) FindFilePathMapping(file DirEntrySubset) (_ *FilePathMapping, ok bool) {
	if file.IsDir() {
		return nil, false
	}
	date, ok := p.mapping[getFileNameWithoutExt(file.Name())]
	if !ok {
		return nil, false
	}

	category := ""
	if p.separate && !isJpegFile(file.Name()) {
		category = strings.ToUpper(path.Ext(file.Name())[1:]) // "ORF", "ARW", etc.
	}

	return &FilePathMapping{
		SrcFilePath: path.Join(p.srcDirPath, file.Name()),
		DstFilePath: path.Join(p.getDstDirPath(category, date), file.Name()),
	}, true
}

func GenerateCopyPlan(files []fs.DirEntry, srcDirPath, dstBaseDirPath string, separate bool) Plan {
	mapping := make(map[string]*time.Time, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !isJpegFile(file.Name()) {
			continue
		}

		srcFullPath := path.Join(srcDirPath, file.Name())

		date, err := loadShootingDateFromExif(srcFullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load %s (%v)\n", srcFullPath, err)
			continue
		}

		fileNameWithoutExt := getFileNameWithoutExt(file.Name())
		mapping[fileNameWithoutExt] = date
	}

	return Plan{
		mapping:        mapping,
		srcDirPath:     srcDirPath,
		dstBaseDirPath: dstBaseDirPath,
		separate:       separate,
	}
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
