package main

import (
	"os"
	"os/exec"
)

const (
	ReturnCodeOK  = 0
	ReturnCodeErr = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) int {
	if err := setEnvs(env); err != nil {
		return ReturnCodeErr
	}

	if len(cmd) < 1 {
		return ReturnCodeErr
	}

	cmdName := cmd[0]
	args := cmd[1:]
	if len(cmdName) < 1 {
		return ReturnCodeErr
	}

	wcExel := exec.Command(cmdName, args...)
	wcExel.Stdout = os.Stdout
	wcExel.Stderr = os.Stderr
	wcExel.Stdin = os.Stdin
	if err := wcExel.Run(); err != nil {
		return ReturnCodeErr
	}

	return ReturnCodeOK
}

func setEnvs(env Environment) error {
	for name, envVal := range env {
		if err := os.Unsetenv(name); err != nil {
			return err
		}

		if !envVal.NeedRemove {
			if err := os.Setenv(name, envVal.Value); err != nil {
				return err
			}
		}
	}

	return nil
}
