package command

import "encoding/json"

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
