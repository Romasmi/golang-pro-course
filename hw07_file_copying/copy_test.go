package main

import (
	"errors"
	"os"
	"testing"
)

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
		source := "./testdata/input.txt"
		desc := source
		err := Copy(source, desc, 0, 0)
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("copying entire file", func(t *testing.T) {
		source := "./testdata/input.txt"
		desc := "./testdata/input_copy.txt"

		err := Copy(source, desc, 0, 0)
		if err != nil {
			t.Error(err)
		}
		MustSizeEqual(t, source, desc)
	})
}

func MustSizeEqual(t *testing.T, f1, f2 string) {
	sourceInfo, err := os.Stat(f1)
	if err != nil {
		t.Fatal(err)
	}
	descInfo, err := os.Stat(f2)
	if err != nil {
		t.Fatal(err)
	}
	if sourceInfo.Size() != descInfo.Size() {
		t.Errorf("invalid file sizes: %d != %d", sourceInfo.Size(), descInfo.Size())
	}
}
