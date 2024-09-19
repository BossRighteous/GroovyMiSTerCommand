package command

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"time"
)

type ProcessStateError struct {
	e string
}

func (err *ProcessStateError) Error() string {
	return err.e
}

type RunResult struct {
	Code         int
	Message      string
	MessageLines []string
	BlitMessage  bool
}

type CommandRunner struct {
	Cmd        *exec.Cmd
	Config     *GMCConfig
	ResultChan chan RunResult
}

func (cmdr *CommandRunner) Cancel() error {
	if cmdr.Cmd != nil {
		if cmdr.Cmd.Process != nil {
			fmt.Println("Killing running process")
			err := cmdr.Cmd.Process.Kill()
			//return cmdr.Cmd.Process.Signal(syscall.SIGTERM)
			time.Sleep(time.Millisecond * 250)
			return err
		} else if cmdr.Cmd.ProcessState != nil {
			fmt.Println("Process already completed, cleaning up")
			fmt.Println(cmdr.Cmd.ProcessState)
			cmdr.Cmd = nil
			return &ProcessStateError{e: cmdr.Cmd.ProcessState.String()}
		}
		cmdr.Cmd = nil
	}
	return nil
}

func (cmdr *CommandRunner) IsRunning() bool {
	return cmdr.Cmd != nil
}

func (cmdr *CommandRunner) ReplaceArgVars(args []string, vars map[string]string) []string {
	return ReplaceArgVars(args, vars)
}

func (cmdr *CommandRunner) Run(gmc GroovyMiSTerCommand) RunResult {
	cErr := cmdr.Cancel()
	if cErr != nil {
		return RunResult{
			Code:        2,
			Message:     cErr.Error(),
			BlitMessage: true,
		}
	}
	fmt.Println("Cancel Error: ", cErr)

	cfgCmd, ok := cmdr.Config.CmdMap[gmc.Cmd]
	if !ok {
		return RunResult{
			Code:        3,
			Message:     "cmd key not found in server config whitelist, aborting",
			BlitMessage: true,
		}
	}

	args := cmdr.ReplaceArgVars(cfgCmd.ExecArgs, gmc.Vars)

	fmt.Println(cfgCmd.WorkDir, cfgCmd.ExecBin, args)
	bin := filepath.Join(cfgCmd.WorkDir, cfgCmd.ExecBin)
	cmd := exec.Command(bin, args...)
	cmdr.Cmd = cmd
	if cfgCmd.WorkDir != "" {
		cmd.Dir = cfgCmd.WorkDir
	}

	stdout, _ := cmd.StderrPipe()
	scanner := bufio.NewScanner(stdout)
	outLines := make([]string, 10)
	err := cmd.Start()
	if err != nil {
		return RunResult{
			Code:        4,
			Message:     err.Error(),
			BlitMessage: true,
		}
	}
	go func() {
		// wait async
		result := RunResult{
			Code:        0,
			BlitMessage: true,
		}

		for scanner.Scan() {
			// Do something with the line here.
			fmt.Println(scanner.Text())
			outLines = append(outLines[1:10], scanner.Text())
		}
		if err := cmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {

				fmt.Println("Exit Status", exiterr.ExitCode())
				result.Code = exiterr.ExitCode()
				result.MessageLines = outLines[:]
			} else {
				result.Code = 1
				result.Message = err.Error()
			}
		}
		cmdr.Cmd = nil
		cmdr.ResultChan <- result
	}()
	return RunResult{
		Code:    -1,
		Message: "Pending Command Resolution",
	}
}
