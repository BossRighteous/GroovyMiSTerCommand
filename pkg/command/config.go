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

type GMCConfigMameGenerator struct {
	RomsDir      string `json:"roms_dir"`
	MamelistPath string `json:"mamelist_path"`
}

type GMCConfigRetroarchGenerator struct {
	PlaylistsDir string `json:"playlists_dir"`
}

type GMCConfigDirectoryGenerator struct {
	Name       string              `json:"name"`
	Dir        string              `json:"dir"`
	Extensions []string            `json:"extensions"`
	Template   GroovyMiSTerCommand `json:"template"`
}

type GMCConfigGenerators struct {
	Mame        GMCConfigMameGenerator        `json:"mame"`
	Retroarch   GMCConfigRetroarchGenerator   `json:"retroarch"`
	Directories []GMCConfigDirectoryGenerator `json:"directories"`
}

type GMCConfig struct {
	MisterHost      string              `json:"mister_host"`
	DisplayMessages bool                `json:"display_messages"`
	Commands        []GMCConfigCommand  `json:"commands"`
	Generators      GMCConfigGenerators `json:"generators"`
	CmdMap          map[string]GMCConfigCommand
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
	config.CmdMap = make(map[string]GMCConfigCommand)
	for i := range config.Commands {
		cmd := config.Commands[i]
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
