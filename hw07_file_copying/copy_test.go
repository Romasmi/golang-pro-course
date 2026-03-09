package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var source = "./testdata/input.txt"

func TestCopy(t *testing.T) {
	unreadableFilename := "./testdata/unreadable.txt"
	existingFilename := "./testdata/existing_file.txt"

	// prepare an unreadable file
	err := os.WriteFile(unreadableFilename, []byte("unreadable"), 0644)
	assert.NoError(t, err)
	defer os.Remove(unreadableFilename)
	err = os.Chmod(unreadableFilename, 0000)
	assert.NoError(t, err)

	// prepare an existing file
	err = os.WriteFile(existingFilename, []byte("some content"), 0644)
	assert.NoError(t, err)

	// prepare a large file
	largeFileSrc, largeFileDst := generateLargeFile(t)
	defer os.Remove(largeFileSrc)
	defer os.Remove(largeFileDst)

	tests := []struct {
		name         string
		from         string
		to           string
		offset       int64
		limit        int64
		wantErr      error
		expectedFile string
	}{
		{
			name:    "file not found",
			from:    "not_existing_file.txt",
			to:      "copied_file.txt",
			wantErr: ErrFileNotFound,
		},
		{
			name:    "unsupported file",
			from:    "/dev/urandom",
			to:      "copied_file.txt",
			wantErr: ErrUnsupportedFile,
		},
		{
			name:    "offset exceeds file size",
			from:    source,
			to:      "copied_file.txt",
			offset:  10000,
			limit:   1000,
			wantErr: ErrOffsetExceedsFileSize,
		},
		{
			name:    "copying itself",
			from:    source,
			to:      source,
			wantErr: nil,
		},
		{
			name:         "valid full copy",
			from:         source,
			to:           addCopyPostfix("./testdata/out_offset0_limit0.txt"),
			expectedFile: "./testdata/out_offset0_limit0.txt",
			wantErr:      nil,
		},
		{
			name:         "valid partial copy: offset 0, limit 10",
			from:         source,
			to:           addCopyPostfix("./testdata/out_offset0_limit10.txt"),
			offset:       0,
			limit:        10,
			expectedFile: "./testdata/out_offset0_limit10.txt",
			wantErr:      nil,
		},
		{
			name:         "valid partial copy: offset 0, limit 10000",
			from:         source,
			to:           addCopyPostfix("./testdata/out_offset0_limit10000.txt"),
			offset:       0,
			limit:        10000,
			expectedFile: "./testdata/out_offset0_limit10000.txt",
			wantErr:      nil,
		},
		{
			name:         "valid partial copy: offset 100, limit 1000",
			from:         source,
			to:           addCopyPostfix("./testdata/out_offset100_limit1000.txt"),
			offset:       100,
			limit:        1000,
			expectedFile: "./testdata/out_offset100_limit1000.txt",
			wantErr:      nil,
		},
		{
			name:         "valid partial copy: offset 6000, limit 1000",
			from:         source,
			to:           addCopyPostfix("./testdata/out_offset6000_limit1000.txt"),
			offset:       6000,
			limit:        1000,
			expectedFile: "./testdata/out_offset6000_limit1000.txt",
			wantErr:      nil,
		},
		{
			name:    "empty file with offset 0",
			from:    "./testdata/empty.txt",
			to:      "empty_copy.txt",
			offset:  0,
			limit:   0,
			wantErr: nil,
		},
		{
			name:    "overwrite existing destination file",
			from:    "./testdata/empty.txt",
			to:      existingFilename,
			wantErr: nil,
		},
		{
			name:         "limit exactly equal to file size",
			from:         source,
			to:           addCopyPostfix("./testdata/out_offset0_limit_full.txt"),
			offset:       0,
			limit:        10000,
			expectedFile: "./testdata/out_offset0_limit0.txt",
			wantErr:      nil,
		},
		{
			name:    "unreadable file",
			from:    unreadableFilename,
			to:      "/tmp/unreadable_copy.txt",
			wantErr: os.ErrPermission,
		},
		{
			name:    "destination directory doesn't exist",
			from:    source,
			to:      "./non_existent_dir/copy.txt",
			wantErr: os.ErrNotExist,
		},
		{
			name:         "full multi-chunk copy",
			from:         largeFileDst,
			to:           addCopyPostfix(largeFileDst),
			expectedFile: largeFileDst,
			wantErr:      nil,
		},
		{
			name:         "full multi-chunk copy",
			from:         largeFileDst,
			to:           addCopyPostfix(largeFileDst),
			expectedFile: largeFileDst,
			wantErr:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantErr == nil && tt.from != tt.to {
				defer os.Remove(tt.to)
			}

			err := Copy(tt.from, tt.to, tt.offset, tt.limit)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			if tt.expectedFile != "" {
				mustFilesEqual(t, tt.to, tt.expectedFile)
			}
		})
	}
}

func addCopyPostfix(path string) string {
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(path, ext)
	return name + "_copy" + ext
}

func mustFilesEqual(t *testing.T, f1, f2 string) {
	t.Helper()

	expectedContent, err := os.ReadFile(f1)
	assert.NoError(t, err)

	actualContent, err := os.ReadFile(f2)
	assert.NoError(t, err)

	assert.Equal(t, expectedContent, actualContent, "file contents should match")
}

func generateLargeFile(t *testing.T) (string, string) {
	// Generate a file significantly larger than bufferDefaultSize (32KB).
	const fileSize = 512 * 1024

	content := make([]byte, fileSize)
	for i := range content {
		content[i] = byte(i % 251)
	}

	src, err := os.CreateTemp("", "copy_src_*.bin")
	assert.NoError(t, err)

	_, err = src.Write(content)
	assert.NoError(t, err)
	assert.NoError(t, src.Close())

	dst, err := os.CreateTemp("", "copy_dst_*.bin")
	assert.NoError(t, err)

	return src.Name(), dst.Name()
}
