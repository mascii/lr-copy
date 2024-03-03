package cpplan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
