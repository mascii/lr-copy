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
	srcDir string
	dstDir string
}

type FilePathMapping struct {
	SrcFilePath string
	DstFilePath string
}

type Plan map[string]*DirectoryMapping

func (p Plan) FindFilePathMapping(file fs.DirEntry) (_ *FilePathMapping, ok bool) {
	if file.IsDir() {
		return nil, false
	}
	dm, ok := p[getFileNameWithoutExt(file.Name())]
	if !ok {
		return nil, false
	}

	return &FilePathMapping{
		SrcFilePath: path.Join(dm.srcDir, file.Name()),
		DstFilePath: path.Join(dm.dstDir, file.Name()),
	}, ok
}

func (p Plan) GetDstDirs() map[string]struct{} {
	paths := make(map[string]struct{}, len(p))
	for _, dm := range p {
		paths[dm.dstDir] = struct{}{}
	}
	return paths
}

func GenerateCopyPlan(files []fs.DirEntry, srcDirPath, dstDirPath string) (Plan, error) {
	plan := make(Plan, len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !isJpegFile(file.Name()) {
			continue
		}

		srcFullPath := path.Join(srcDirPath, file.Name())

		t, err := loadShootingDateFromExif(srcFullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load %s (%v)\n", srcFullPath, err)
			continue
		}

		fileNameWithoutExt := getFileNameWithoutExt(file.Name())
		plan[fileNameWithoutExt] = &DirectoryMapping{
			srcDir: srcDirPath,
			dstDir: path.Join(dstDirPath, t.Format("2006/2006-01-02")), // Lightroom のフォルダ名の形式に合わせる
		}
	}

	return plan, nil
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
