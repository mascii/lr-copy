//go:build windows
// +build windows

package cpplan

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type dirEntryMock struct {
	name string
}

func (d *dirEntryMock) Name() string { return d.name }

func Test_GenerateCopyPlan(t *testing.T) {
	date1 := time.Date(2024, 2, 12, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2024, 3, 2, 0, 0, 0, 0, time.UTC)
	files := []*dirEntryMock{
		{
			name: "example001.jpg",
		},
		{
			name: "example001.raw",
		},
		{
			name: "example002.jpg",
		},
		{
			name: "example002.raw",
		},
		{
			name: "error.jpg",
		},
	}
	loadShootingDateFromJpeg := func(filePath string) (*time.Time, error) {
		switch filePath {
		case "C:\\path\\to\\photos\\example001.jpg":
			return &date1, nil
		case "C:\\path\\to\\photos\\example002.jpg":
			return &date2, nil
		case "C:\\path\\to\\photos\\error.jpg":
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
					SrcFilePath: "C:\\path\\to\\photos\\example001.jpg",
					DstFilePath: "D:\\photos\\2024\\2024-02-12\\example001.jpg",
				},
				{
					SrcFilePath: "C:\\path\\to\\photos\\example001.raw",
					DstFilePath: "D:\\photos\\RAW\\2024\\2024-02-12\\example001.raw",
				},
				{
					SrcFilePath: "C:\\path\\to\\photos\\example002.jpg",
					DstFilePath: "D:\\photos\\2024\\2024-03-02\\example002.jpg",
				},
				{
					SrcFilePath: "C:\\path\\to\\photos\\example002.raw",
					DstFilePath: "D:\\photos\\RAW\\2024\\2024-03-02\\example002.raw",
				},
			},
		},
		{
			separate: false,
			expected: []*FilePathMapping{
				{
					SrcFilePath: "C:\\path\\to\\photos\\example001.jpg",
					DstFilePath: "D:\\photos\\2024\\2024-02-12\\example001.jpg",
				},
				{
					SrcFilePath: "C:\\path\\to\\photos\\example001.raw",
					DstFilePath: "D:\\photos\\2024\\2024-02-12\\example001.raw",
				},
				{
					SrcFilePath: "C:\\path\\to\\photos\\example002.jpg",
					DstFilePath: "D:\\photos\\2024\\2024-03-02\\example002.jpg",
				},
				{
					SrcFilePath: "C:\\path\\to\\photos\\example002.raw",
					DstFilePath: "D:\\photos\\2024\\2024-03-02\\example002.raw",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("separate=%v", tc.separate), func(t *testing.T) {
			cfg := generateCopyPlanConfig{
				srcDirPath:     "C:\\path\\to\\photos",
				dstBaseDirPath: "D:\\photos",
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
