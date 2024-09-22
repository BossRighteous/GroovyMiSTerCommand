//go:build windows

package command

import (
	"os"
	"strings"

	"github.com/BossRighteous/GroovyMiSTerCommand/pkg/winutils"
)

func CloseProcess(proc *os.Process, exitStrat string) error {
	if strings.ToLower(exitStrat) == "kill" {
		return proc.Kill()
	}
	// "quit"
	return winutils.CloseWindow(proc.Pid, true)

}
