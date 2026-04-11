package main

import (
	"errors"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	outCmd := exec.Command(cmd[0], cmd[1:]...)
	outCmd.Stdout = os.Stdout
	outCmd.Stderr = os.Stderr
	outCmd.Stdin = os.Stdin
	outCmd.Env = buildEnv(os.Environ(), env)

	err := outCmd.Run()

	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return 1
	}

	return 0
}

func buildEnv(currentEnv []string, env Environment) []string {
	resMap := make(map[string]string)
	for _, e := range currentEnv {
		if k, v, ok := strings.Cut(e, "="); ok {
			resMap[k] = v
		}
	}

	for k, v := range env {
		if v.NeedRemove {
			delete(resMap, k)
		} else {
			resMap[k] = v.Value
		}
	}

	res := make([]string, 0, len(resMap))
	for k, v := range resMap {
		res = append(res, k+"="+v)
	}
	return res
}
