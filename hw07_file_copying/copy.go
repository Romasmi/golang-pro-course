package main

import (
	"errors"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotFound          = errors.New("file not found")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if !fileExists(fromPath) {
		return ErrFileNotFound
	}

	fileInfo, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if !fileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	if offset >= fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// Place your code here.
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
