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
		err := os.WriteFile(filepath.Join(dir, "FOO"), []byte("123"), 0644)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(filepath.Join(dir, "BAR"), []byte("value"), 0644)
		if err != nil {
			t.Fatal(err)
		}

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
		err := os.WriteFile(filepath.Join(dir, "UNSET"), []byte{}, 0644)
		if err != nil {
			t.Fatal(err)
		}

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
		err := os.WriteFile(filepath.Join(dir, "VAR"), []byte("value  \t \nother line"), 0644)
		if err != nil {
			t.Fatal(err)
		}

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
		err := os.WriteFile(filepath.Join(dir, "NULL"), []byte("hello\x00world"), 0644)
		if err != nil {
			t.Fatal(err)
		}

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
		err := os.WriteFile(filepath.Join(dir, "INVALID=NAME"), []byte("value"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		_, err = ReadDir(dir)
		if err == nil {
			t.Error("ReadDir() expected error for invalid filename, got nil")
		}
	})

	t.Run("ignore subdirectories", func(t *testing.T) {
		dir := t.TempDir()
		err := os.Mkdir(filepath.Join(dir, "SUBDIR"), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(filepath.Join(dir, "FOO"), []byte("123"), 0644)
		if err != nil {
			t.Fatal(err)
		}

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
