package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var source = "./testdata/input.txt"

func TestCopy(t *testing.T) {

	t.Run("copying not existing file", func(t *testing.T) {
		err := Copy("not_existing_file.txt", "copied_file.txt", 0, 0)
		if err == nil || !errors.Is(err, ErrFileNotFound) {
			t.Errorf("invalid error for not existing file: %v", err)
		}
	})

	t.Run("copying unsupported file", func(t *testing.T) {
		err := Copy("/dev/urandom", "copied_file.txt", 0, 0)
		if err == nil || !errors.Is(err, ErrUnsupportedFile) {
			t.Errorf("invalid error for unsupported file: %v", err)
		}
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "test_file.*.txt")
		if err != nil {
			t.Fatal(err)
		}
		defer func(name string) {
			err := os.Remove(name)
			if err != nil {
				t.Fatal(err)
			}
		}(tmpFile.Name())

		content := []byte("test content")
		if _, err := tmpFile.Write(content); err != nil {
			t.Fatal(err)
		}
		if err := tmpFile.Close(); err != nil {
			t.Fatal(err)
		}

		err = Copy(tmpFile.Name(), "copied_file.txt", int64(len(content)+1), 1000)
		if err == nil || !errors.Is(err, ErrOffsetExceedsFileSize) {
			t.Errorf("invalid error for offset exceeding file size: %v", err)
		}
	})

	t.Run("copying itself", func(t *testing.T) {
		desc := source
		err := Copy(source, desc, 0, 0)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("copying entire file", func(t *testing.T) {
		testCase(t, 0, 0, "out_offset0_limit0.txt")
	})

	t.Run("copy file: offset 0, limit 10", func(t *testing.T) {
		testCase(t, 0, 10, "out_offset0_limit10.txt")
	})

	t.Run("copy file: offset 0, limit 10000 exceeds file size", func(t *testing.T) {
		testCase(t, 0, 10000, "out_offset0_limit10000.txt")
	})

	t.Run("copy file: offset 100, limit 1000", func(t *testing.T) {
		testCase(t, 100, 1000, "out_offset100_limit1000.txt")
	})

	t.Run("copy file: offset 6000, limit 1000", func(t *testing.T) {
		testCase(t, 6000, 1000, "out_offset6000_limit1000.txt")
	})
}

func addCopyPostfix(path string) string {
	ext := filepath.Ext(path)
	name := strings.TrimSuffix(path, ext)
	return name + "_copy" + ext
}

func testCase(t *testing.T, offset, limit int, testFile string) {
	desc := addCopyPostfix("./testdata/" + testFile)
	//defer func(name string) {
	//	err := os.Remove(name)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//}(desc)

	err := Copy(source, desc, int64(offset), int64(limit))
	if err != nil {
		t.Error(err)
	}
	mustFilesEqual(t, desc, "./testdata/"+testFile)
}

func mustFilesEqual(t *testing.T, f1, f2 string) {
	expectedContent, err := os.ReadFile(f1)
	assert.NoError(t, err)

	actualContent, err := os.ReadFile(f2)
	assert.NoError(t, err)

	assert.Equal(t, expectedContent, actualContent, "file contents should match")
}
