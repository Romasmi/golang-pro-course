package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/sync/errgroup"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

type EnvVar struct {
	Name  string
	Value EnvValue
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := getFolderFiles(dir)
	if err != nil {
		return nil, err
	}
	envs := make(Environment, len(files))

	var (
		g  errgroup.Group
		mu sync.Mutex
	)

	for _, file := range files {
		g.Go(func() error {
			env, err := getEnv(filepath.Join(dir, file))
			if err != nil {
				return err
			}

			mu.Lock()
			envs[env.Name] = env.Value
			mu.Unlock()
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return envs, nil
}

func getFolderFiles(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	out := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			out = append(out, entry.Name())
		}
	}
	return out, nil
}

func getEnv(filepath string) (*EnvVar, error) {
	info, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}

	if !isValidEnvName(info.Name()) {
		return nil, fmt.Errorf("invalid env name: %s for file %s", info.Name(), filepath)
	}

	if info.Size() == 0 {
		return &EnvVar{info.Name(), EnvValue{NeedRemove: true}}, nil
	}

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var firstLine string
	if scanner.Scan() {
		firstLine = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &EnvVar{info.Name(), EnvValue{Value: processEnvValue(firstLine)}}, nil
}

func processEnvValue(value string) string {
	out := strings.TrimRightFunc(value, unicode.IsSpace)
	return string(bytes.ReplaceAll([]byte(out), []byte{0x00}, []byte("\n")))
}

func isValidEnvName(name string) bool {
	return !strings.Contains(name, "=")
}
