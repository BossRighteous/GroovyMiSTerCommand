package command

import (
	"encoding/json"
	"strings"
)

type GroovyMiSTerCommand struct {
	Cmd  string            `json:"cmd"`
	Vars map[string]string `json:"vars"`
	Raw  []byte
}

func ParseGMC(cmdBytes []byte) (GroovyMiSTerCommand, error) {
	cmd := GroovyMiSTerCommand{
		Raw: cmdBytes,
	}
	err := json.Unmarshal(cmdBytes, &cmd)
	return cmd, err
}

func ReplaceArgVars(args []string, vars map[string]string) []string {
	nArgs := make([]string, len(args))
	for i := range args {
		arg := args[i]
		for k, v := range vars {
			pattern := "${" + k + "}"
			arg = strings.Replace(arg, pattern, v, -1)
		}
		nArgs[i] = arg
	}
	return nArgs
}
