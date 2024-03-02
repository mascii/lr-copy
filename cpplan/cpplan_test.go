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

func TestFindFilePathMapping_Separate_Is_True(t *testing.T) {
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

	m, ok := p.FindFilePathMapping(&DirEntryMock{
		name:  "directory_name",
		isDir: true,
	})
	assert.False(t, ok)
	assert.Nil(t, m)

	m, ok = p.FindFilePathMapping(&DirEntryMock{
		name:  "example001.jpg",
		isDir: false,
	})
	assert.True(t, ok)
	assert.Equal(t, "/path/to/photos/example001.jpg", m.SrcFilePath)
	assert.Equal(t, "/home/user/photos/2022/2022-01-01/example001.jpg", m.DstFilePath)

	m, ok = p.FindFilePathMapping(&DirEntryMock{
		name:  "example001.raw",
		isDir: false,
	})
	assert.True(t, ok)
	assert.Equal(t, "/path/to/photos/example001.raw", m.SrcFilePath)
	assert.Equal(t, "/home/user/photos/RAW/2022/2022-01-01/example001.raw", m.DstFilePath)

	m, ok = p.FindFilePathMapping(&DirEntryMock{
		name:  "example003.jpg",
		isDir: false,
	})
	assert.False(t, ok)
	assert.Nil(t, m)

}

func Test_FindFilePathMapping_Separate_Is_False(t *testing.T) {
	d := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)

	p := Plan{
		mapping: map[string]*time.Time{
			"example001": &d,
			"example002": &d,
		},
		srcDirPath:     "/path/to/photos",
		dstBaseDirPath: "/home/user/photos",
		separate:       false,
	}

	m, ok := p.FindFilePathMapping(&DirEntryMock{
		name:  "directory_name",
		isDir: true,
	})
	assert.False(t, ok)
	assert.Nil(t, m)

	m, ok = p.FindFilePathMapping(&DirEntryMock{
		name:  "example001.jpg",
		isDir: false,
	})
	assert.True(t, ok)
	assert.Equal(t, "/path/to/photos/example001.jpg", m.SrcFilePath)
	assert.Equal(t, "/home/user/photos/2022/2022-01-01/example001.jpg", m.DstFilePath)

	m, ok = p.FindFilePathMapping(&DirEntryMock{
		name:  "example001.raw",
		isDir: false,
	})
	assert.True(t, ok)
	assert.Equal(t, "/path/to/photos/example001.raw", m.SrcFilePath)
	assert.Equal(t, "/home/user/photos/2022/2022-01-01/example001.raw", m.DstFilePath)

	m, ok = p.FindFilePathMapping(&DirEntryMock{
		name:  "example003.jpg",
		isDir: false,
	})
	assert.False(t, ok)
	assert.Nil(t, m)
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
