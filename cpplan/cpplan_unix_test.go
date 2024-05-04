//go:build unix
// +build unix

package cpplan

import (
	"errors"
	"fmt"
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
	loadShootingDateFromJpeg := func(filePath string) (*time.Time, error) {
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
	}

	testCases := []struct {
		separate bool
		expected []*FilePathMapping
	}{
		{
			separate: true,
			expected: []*FilePathMapping{
				{
					SrcFilePath: "/path/to/photos/example001.jpg",
					DstFilePath: "/home/user/photos/2024/2024-02-12/example001.jpg",
				},
				{
					SrcFilePath: "/path/to/photos/example001.raw",
					DstFilePath: "/home/user/photos/RAW/2024/2024-02-12/example001.raw",
				},
				{
					SrcFilePath: "/path/to/photos/example002.jpg",
					DstFilePath: "/home/user/photos/2024/2024-03-02/example002.jpg",
				},
				{
					SrcFilePath: "/path/to/photos/example002.raw",
					DstFilePath: "/home/user/photos/RAW/2024/2024-03-02/example002.raw",
				},
			},
		},
		{
			separate: false,
			expected: []*FilePathMapping{
				{
					SrcFilePath: "/path/to/photos/example001.jpg",
					DstFilePath: "/home/user/photos/2024/2024-02-12/example001.jpg",
				},
				{
					SrcFilePath: "/path/to/photos/example001.raw",
					DstFilePath: "/home/user/photos/2024/2024-02-12/example001.raw",
				},
				{
					SrcFilePath: "/path/to/photos/example002.jpg",
					DstFilePath: "/home/user/photos/2024/2024-03-02/example002.jpg",
				},
				{
					SrcFilePath: "/path/to/photos/example002.raw",
					DstFilePath: "/home/user/photos/2024/2024-03-02/example002.raw",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("separate=%v", tc.separate), func(t *testing.T) {
			cfg := generateCopyPlanConfig{
				srcDirPath:     "/path/to/photos",
				dstBaseDirPath: "/home/user/photos",
				separate:       tc.separate,
				shootingDateExtractors: map[string]func(filePath string) (*time.Time, error){
					"JPG": loadShootingDateFromJpeg,
				},
			}
			plan := GenerateCopyPlan(files, cfg)
			assert.Equal(t, tc.expected, plan)
		})
	}
}
