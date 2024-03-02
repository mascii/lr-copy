package cpplan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DirEntryMock struct {
	name  string
	isDir bool
}

func (d *DirEntryMock) Name() string { return d.name }
func (d *DirEntryMock) IsDir() bool  { return d.isDir }

func TestHasNoFilesToCopy(t *testing.T) {
	d := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	p := Plan{
		mapping: map[string]*time.Time{
			"example001": &d,
			"example002": &d,
		},
		srcDirPath:     "/path/to/photos",
		dstBaseDirPath: "/home/user/photos",
		separate:       true,
	}

	assert.False(t, p.HasNoFilesToCopy())

	p = Plan{
		mapping:        map[string]*time.Time{},
		srcDirPath:     "/path/to/photos",
		dstBaseDirPath: "/home/user/photos",
		separate:       true,
	}

	assert.True(t, p.HasNoFilesToCopy())
}

func Test_FindFilePathMapping(t *testing.T) {
	d := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		separate bool
		file     *DirEntryMock
		ok       bool
		expected *FilePathMapping
	}{
		{
			separate: true,
			file: &DirEntryMock{
				name:  "directory_name",
				isDir: true,
			},
			expected: nil,
			ok:       false,
		},
		{
			separate: true,
			file: &DirEntryMock{
				name:  "example001.jpg",
				isDir: false,
			},
			expected: &FilePathMapping{
				SrcFilePath: "/path/to/photos/example001.jpg",
				DstFilePath: "/home/user/photos/2022/2022-01-01/example001.jpg",
			},
			ok: true,
		},
		{
			separate: true,
			file: &DirEntryMock{
				name:  "example001.raw",
				isDir: false,
			},
			expected: &FilePathMapping{
				SrcFilePath: "/path/to/photos/example001.raw",
				DstFilePath: "/home/user/photos/RAW/2022/2022-01-01/example001.raw",
			},
			ok: true,
		},
		{
			separate: false,
			file: &DirEntryMock{
				name:  "example001.jpg",
				isDir: false,
			},
			expected: &FilePathMapping{
				SrcFilePath: "/path/to/photos/example001.jpg",
				DstFilePath: "/home/user/photos/2022/2022-01-01/example001.jpg",
			},
			ok: true,
		},
		{
			separate: false,
			file: &DirEntryMock{
				name:  "example001.raw",
				isDir: false,
			},
			expected: &FilePathMapping{
				SrcFilePath: "/path/to/photos/example001.raw",
				DstFilePath: "/home/user/photos/2022/2022-01-01/example001.raw",
			},
			ok: true,
		},
		{
			separate: true,
			file: &DirEntryMock{
				name:  "example003.jpg",
				isDir: false,
			},
			expected: nil,
			ok:       false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.file.name, func(t *testing.T) {
			p := Plan{
				mapping: map[string]*time.Time{
					"example001": &d,
					"example002": &d,
				},
				srcDirPath:     "/path/to/photos",
				dstBaseDirPath: "/home/user/photos",
				separate:       tc.separate,
			}

			m, ok := p.FindFilePathMapping(tc.file)
			assert.Equal(t, tc.ok, ok)
			assert.Equal(t, tc.expected, m)
		})
	}
}

func Test_GenerateCopyPlan(t *testing.T) {
	date1 := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC)

	files := []*DirEntryMock{
		{
			name:  "example001.jpg",
			isDir: false,
		},
		{
			name:  "example001.raw",
			isDir: false,
		},
		{
			name:  "example002.jpg",
			isDir: false,
		},
		{
			name:  "example002.raw",
			isDir: false,
		},
		{
			name:  "directory_name",
			isDir: true,
		},
	}
	cfg := GenerateCopyPlanConfig{
		SrcDirPath:     "/path/to/photos",
		DstBaseDirPath: "/home/user/photos",
		Separate:       true,
		LoadShootingDateFromExif: func(filePath string) (*time.Time, error) {
			switch filePath {
			case "/path/to/photos/example001.jpg":
				return &date1, nil
			case "/path/to/photos/example002.jpg":
				return &date2, nil
			default:
				assert.FailNow(t, "unexpected file path: %s", filePath)
				panic(filePath)
			}
		},
	}

	plan := GenerateCopyPlan(files, cfg)
	assert.Equal(t, Plan{
		mapping: map[string]*time.Time{
			"example001": &date1,
			"example002": &date2,
		},
		srcDirPath:     "/path/to/photos",
		dstBaseDirPath: "/home/user/photos",
		separate:       true,
	}, plan)
}

func Test_isJpegFile(t *testing.T) {
	testCases := []struct {
		fileName string
		expected bool
	}{
		{"example.jpg", true},
		{"example.jpeg", true},
		{"example.JPG", true},
		{"example.JPEG", true},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			result := isJpegFile(tc.fileName)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func Test_getFileNameWithoutExt(t *testing.T) {
	testCases := []struct {
		fileName string
		expected string
	}{
		{"example.raw", "example"},
		{"example.jpg", "example"},
		{"example.JPEG", "example"},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			result := getFileNameWithoutExt(tc.fileName)
			assert.Equal(t, tc.expected, result)
		})
	}
}
