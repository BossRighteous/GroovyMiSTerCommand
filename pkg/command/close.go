//go:build !windows

package command

import (
	"os"
	"strings"
	"syscall"
)

func CloseProcess(proc *os.Process, exitStrat string) error {
	if strings.ToLower(exitStrat) == "kill" {
		return proc.Kill()
	}
	return proc.Signal(syscall.SIGINT)
}
