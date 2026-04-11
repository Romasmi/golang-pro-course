package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/samber/lo"
)

var (
	ErrFolderNotExist = errors.New("folder does not exist")
	ErrEmptyFile      = errors.New("file is empty")
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

	for _, file := range files {
		env, err := getEnv(dir + "/" + file)
		if err != nil {
			return nil, err
		}
		envs[env.Name] = env.Value
	}
	return envs, nil
}

func getFolderFiles(path string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	return lo.Map(entries, func(entry os.DirEntry, _ int) string {
		return entry.Name()
	}), nil
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
		return &EnvVar{processEnvValue(info.Name()), EnvValue{NeedRemove: true}}, nil
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
	return &EnvVar{processEnvValue(info.Name()), EnvValue{Value: firstLine}}, nil
}

func processEnvValue(value string) string {
	out := strings.TrimRightFunc(value, unicode.IsSpace)
	return string(bytes.ReplaceAll([]byte(out), []byte{0x00}, []byte("\t")))
}

func isValidEnvName(name string) bool {
	return !strings.Contains(name, "=")
}
