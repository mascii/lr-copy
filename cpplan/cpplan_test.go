package cpplan

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_getExtByFileName(t *testing.T) {
	testCases := []struct {
		fileName string
		expected string
	}{
		{"example.jpg", "JPG"},
		{"example.jpeg", "JPEG"},
		{"example.JPG", "JPG"},
		{"example.JPEG", "JPEG"},
		{"example.orf", "ORF"},
		{"example.ORF", "ORF"},
		{"example", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.fileName, func(t *testing.T) {
			result := getExtByFileName(tc.fileName)
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
