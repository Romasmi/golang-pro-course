package main

import (
	"context"
	"errors"
	"os"
	"sync"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotFound          = errors.New("file not found")
)

var bufferSize = 1024 * 32

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validate(fromPath, offset); err != nil {
		return err
	}
	if fromPath == toPath {
		return nil
	}

	buffer := make(chan []byte)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	wg.Go(func() {
		readFile(ctx, cancel, fromPath, buffer, offset, limit)
	})

	wg.Go(func() {
		writeFile(ctx, cancel, toPath, buffer)
	})

	wg.Wait()
	return nil
}

func readFile(ctx context.Context, cancel context.CancelFunc, fromPath string, bufChan chan<- []byte, offset, limit int64) {
	defer close(bufChan)
	f, err := os.Open(fromPath)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)
	if err != nil {
		cancel()
		return
	}

	fi, err := f.Stat()
	if err != nil {
		cancel()
		return
	}
	if limit == 0 {
		limit = fi.Size()
	}
	left := limit
	bufferSize = int(min(int64(bufferSize), limit))
	buffer := make([]byte, bufferSize)
	_, err = f.Seek(offset, 0)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			b, err := f.Read(buffer)
			if err != nil {
				cancel()
				return
			}
			left -= int64(b)
			bufChan <- buffer[:b]
			if left <= 0 {
				return
			}
			if left < int64(bufferSize) {
				buffer = buffer[:left]
			}
		}
	}
}

func writeFile(ctx context.Context, cancel context.CancelFunc, toPath string, bufChan <-chan []byte) {
	for {
		select {
		case <-ctx.Done():
			return
		case buf, ok := <-bufChan:
			if !ok {
				return
			}
			err := os.WriteFile(toPath, buf, 0644)
			if err != nil {
				cancel()
			}
		}
	}
}

func validate(fromPath string, offset int64) error {
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

	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
