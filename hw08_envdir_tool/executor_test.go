package main

import (
	"os"
	"testing"
)

func TestRunCmd(t *testing.T) {
	t.Run("basic success", func(t *testing.T) {
		env := Environment{
			"FOO": {Value: "BAR", NeedRemove: false},
		}

		cmd := []string{"sh", "-c", "exit 0"}
		code := RunCmd(cmd, env)
		if code != 0 {
			t.Errorf("RunCmd() = %d, want 0", code)
		}
	})

	t.Run("command not found", func(t *testing.T) {
		code := RunCmd([]string{"non-existent-command-12345"}, nil)
		if code == 0 {
			t.Errorf("RunCmd() with non-existent command returned 0, want non-zero")
		}
	})

	t.Run("command returns error code", func(t *testing.T) {
		cmd := []string{"sh", "-c", "exit 42"}
		code := RunCmd(cmd, nil)
		if code != 42 {
			t.Errorf("RunCmd() = %d, want 42", code)
		}
	})

	t.Run("environment variables set", func(t *testing.T) {
		env := Environment{
			"TEST_VAR": {Value: "123", NeedRemove: false},
		}
		// Using sh -c '[ "$TEST_VAR" = "123" ]' to check the value.
		cmd := []string{"sh", "-c", "[ \"$TEST_VAR\" = \"123\" ]"}
		code := RunCmd(cmd, env)
		if code != 0 {
			t.Errorf("RunCmd() = %d, want 0 when environment variable is set", code)
		}
	})

	t.Run("environment variables removed", func(t *testing.T) {
		os.Setenv("REMOVE_ME", "I am here")
		defer os.Unsetenv("REMOVE_ME")

		env := Environment{
			"REMOVE_ME": {NeedRemove: true},
		}
		// Using sh -c '[ -z "$REMOVE_ME" ]' to check that the variable is empty/unset.
		cmd := []string{"sh", "-c", "[ -z \"$REMOVE_ME\" ]"}
		code := RunCmd(cmd, env)
		if code != 0 {
			t.Errorf("RunCmd() = %d, want 0 when environment variable is removed", code)
		}
	})
}
