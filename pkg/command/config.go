package command

import (
	"encoding/json"
	"os"
)

type GMCConfigCommand struct {
	Cmd      string   `json:"cmd"`
	WorkDir  string   `json:"work_dir"`
	ExecBin  string   `json:"exec_bin"`
	ExecArgs []string `json:"exec_args"`
}

type GMCConfig struct {
	MisterHost string             `json:"mister_host"`
	Commands   []GMCConfigCommand `json:"commands"`
	CmdMap     map[string]GMCConfigCommand
}

func LoadConfigFromPath(path string) (*GMCConfig, error) {
	var config GMCConfig
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dat, &config)
	if err != nil {
		return nil, err
	}

	for i := range config.Commands {
		cmd := config.Commands[i]
		config.CmdMap = make(map[string]GMCConfigCommand)
		config.CmdMap[cmd.Cmd] = cmd
	}

	// Add built in commands
	config.CmdMap["unload"] = GMCConfigCommand{
		Cmd:      "unload",
		ExecBin:  "echo",
		ExecArgs: []string{""},
	}

	return &config, nil
}
