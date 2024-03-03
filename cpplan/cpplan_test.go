package cpplan

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dirEntryMock struct {
	name  string
	isDir bool
}

func (d *dirEntryMock) Name() string { return d.name }
func (d *dirEntryMock) IsDir() bool  { return d.isDir }

func Test_HasNoFilesToCopy(t *testing.T) {
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
		file     *dirEntryMock
		ok       bool
		expected *FilePathMapping
	}{
		{
			separate: true,
			file: &dirEntryMock{
				name:  "directory_name",
				isDir: true,
			},
			expected: nil,
			ok:       false,
		},
		{
			separate: true,
			file: &dirEntryMock{
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
			file: &dirEntryMock{
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
			file: &dirEntryMock{
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
			file: &dirEntryMock{
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
			file: &dirEntryMock{
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

	files := []*dirEntryMock{
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
			name:  "error.jpg",
			isDir: false,
		},
		{
			name:  "directory_name",
			isDir: true,
		},
	}
	cfg := generateCopyPlanConfig{
		srcDirPath:     "/path/to/photos",
		dstBaseDirPath: "/home/user/photos",
		separate:       true,
		loadShootingDateFromExif: func(filePath string) (*time.Time, error) {
			switch filePath {
			case "/path/to/photos/example001.jpg":
				return &date1, nil
			case "/path/to/photos/example002.jpg":
				return &date2, nil
			case "/path/to/photos/error.jpg":
				return nil, errors.New("error.jpg")
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

func Test_dateToLightroomFormat(t *testing.T) {
	d := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, "2024/2024-02-12", dateToLightroomFormat(&d))
}
