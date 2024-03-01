package cpplan

import (
	"io/fs"
	"os"
	"path"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type Plan map[string]DirectoryMapping

func (p Plan) Lookup(fileName string) (DirectoryMapping, bool) {
	dm, ok := p[getFileNameWithoutExt(fileName)]
	return dm, ok

}

type DirectoryMapping struct {
	SrcDir string
	DstDir string
}

func (p *DirectoryMapping) SrcFullPath(fileName string) string {
	return path.Join(p.SrcDir, fileName)
}
func (p *DirectoryMapping) DstFullPath(fileName string) string {
	return path.Join(p.DstDir, fileName)
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
			return nil, err
		}

		fileNameWithoutExt := getFileNameWithoutExt(file.Name())
		plan[fileNameWithoutExt] = DirectoryMapping{
			SrcDir: srcDirPath,
			DstDir: path.Join(dstDirPath, t.Format("2006/2006-01-02")),
		}
	}

	return plan, nil
}

func isJpegFile(name string) bool {
	ext := path.Ext(name)
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
