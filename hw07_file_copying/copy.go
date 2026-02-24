package main

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFileNotFound          = errors.New("file not found")
)

const bufferDefaultSize = 1024 * 32

var ShowProgressBar = false

func Copy(fromPath, toPath string, offset, limit int64) error {
	if err := validate(fromPath, offset); err != nil {
		return err
	}

	if fromPath == toPath {
		return nil
	}

	var bar *pb.ProgressBar
	if ShowProgressBar {
		var err error
		bar, err = newProgressBar(fromPath, offset, limit)
		if err != nil {
			return err
		}
		bar.Start()
		defer bar.Finish()
	}

	buffer := make(chan []byte)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go func() {
		defer wg.Done()
		readFile(ctx, cancel, fromPath, buffer, offset, limit)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		writeFile(ctx, cancel, toPath, buffer, bar)
	}()

	wg.Wait()
	return nil
}

func readFile(ctx context.Context, cancel context.CancelFunc, from string, buf chan<- []byte, offset, limit int64) {
	defer close(buf)

	f, err := os.Open(from)
	if err != nil {
		cancel()
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			cancel()
			panic(err)
		}
	}(f)

	fi, err := f.Stat()
	if err != nil {
		cancel()
		return
	}

	remaining := fi.Size() - offset
	if remaining <= 0 {
		return
	}

	if limit <= 0 || limit > remaining {
		limit = remaining
	}

	left := limit
	bufferSize := int(min(int64(bufferDefaultSize), left))
	if bufferSize <= 0 {
		return
	}

	buffer := make([]byte, bufferSize)

	_, err = f.Seek(offset, 0)
	if err != nil {
		cancel()
		return
	}

	for left > 0 {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if int64(len(buffer)) > left {
			buffer = buffer[:left]
		}

		n, err := f.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			cancel()
			return
		}

		if n > 0 {
			left -= int64(n)
			buf <- buffer[:n]
		}
	}
}

func writeFile(ctx context.Context, cancel context.CancelFunc, toPath string, bufChan <-chan []byte, bar *pb.ProgressBar) {
	out, err := os.Create(toPath)
	if err != nil {
		cancel()
		return
	}
	defer out.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case buf, ok := <-bufChan:
			if !ok {
				return
			}
			n, err := out.Write(buf)
			if bar != nil && n > 0 {
				bar.Add(n)
			}
			if err != nil {
				cancel()
				return
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

func newProgressBar(from string, offset, limit int64) (*pb.ProgressBar, error) {
	fs, err := os.Stat(from)
	if err != nil {
		return nil, err
	}

	remaining := fs.Size() - offset
	if remaining < 0 {
		remaining = 0
	}

	if limit <= 0 || limit > remaining {
		limit = remaining
	}

	bar := pb.New64(limit)
	bar.SetTemplateString(`{{percent .}} {{bar . }} {{counters .}}`)
	return bar, nil
}
