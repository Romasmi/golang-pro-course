package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("valid files", func(t *testing.T) {
		dir := t.TempDir()

		createTestFile(t, dir, "FOO", "123")
		createTestFile(t, dir, "BAR", "value")

		got, err := ReadDir(dir)
		if err != nil {
			t.Errorf("ReadDir() error = %v", err)
			return
		}

		want := Environment{
			"FOO": {Value: "123", NeedRemove: false},
			"BAR": {Value: "value", NeedRemove: false},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("ReadDir() got = %v, want %v", got, want)
		}
	})

	t.Run("empty file - remove env var", func(t *testing.T) {
		dir := t.TempDir()

		createTestFile(t, dir, "UNSET", "")

		got, err := ReadDir(dir)
		if err != nil {
			t.Errorf("ReadDir() error = %v", err)
			return
		}

		want := Environment{
			"UNSET": {NeedRemove: true},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("ReadDir() got = %v, want %v", got, want)
		}
	})

	t.Run("trim trailing spaces and tabs", func(t *testing.T) {
		dir := t.TempDir()

		// "value  \t " should be trimmed to "value"
		createTestFile(t, dir, "VAR", "value  \t \nother line")

		got, err := ReadDir(dir)
		if err != nil {
			t.Errorf("ReadDir() error = %v", err)
			return
		}

		want := Environment{
			"VAR": {Value: "value", NeedRemove: false},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("ReadDir() got = %v, want %v", got, want)
		}
	})

	t.Run("replace null bytes with newline", func(t *testing.T) {
		dir := t.TempDir()

		// "hello\x00world" should be "hello\nworld"
		createTestFile(t, dir, "NULL", "hello\x00world")

		got, err := ReadDir(dir)
		if err != nil {
			t.Errorf("ReadDir() error = %v", err)
			return
		}

		want := Environment{
			"NULL": {Value: "hello\nworld", NeedRemove: false},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("ReadDir() got = %v, want %v", got, want)
		}
	})

	t.Run("invalid name with equals", func(t *testing.T) {
		dir := t.TempDir()

		createTestFile(t, dir, "INVALID=NAME", "value")

		_, err := ReadDir(dir)
		if err == nil {
			t.Error("ReadDir() expected error for invalid filename, got nil")
		}
	})

	t.Run("ignore subdirectories", func(t *testing.T) {
		dir := t.TempDir()
		err := os.Mkdir(filepath.Join(dir, "SUBDIR"), 0o755)
		if err != nil {
			t.Fatal(err)
		}

		createTestFile(t, dir, "FOO", "123")

		got, err := ReadDir(dir)
		if err != nil {
			t.Errorf("ReadDir() error = %v", err)
			return
		}

		want := Environment{
			"FOO": {Value: "123", NeedRemove: false},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("ReadDir() got = %v, want %v", got, want)
		}
	})

	t.Run("non-existent directory", func(t *testing.T) {
		_, err := ReadDir("/non/existent/path/for/envdir/test")
		if err == nil {
			t.Error("ReadDir() expected error for non-existent directory, got nil")
		}
	})
}

func createTestFile(t *testing.T, dir, name, content string) {
	t.Helper()
	err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644)
	if err != nil {
		t.Fatal(err)
	}
}
