package generators

import (
	"encoding/json"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/BossRighteous/GroovyMiSTerCommand/pkg/command"
)

type GroovyMiSTerCommandNoRaw struct {
	Cmd  string            `json:"cmd"`
	Vars map[string]string `json:"vars"`
}

func scanGlobFiles(dir string, extensions []string) ([]string, error) {
	if len(extensions) == 0 {
		extensions = append(extensions, "*")
	}
	levels := 4
	files := make([]string, 0)
	globPattern := dir + "/*"
	for levels > 0 {
		levels--
		for _, ext := range extensions {
			extPattern := globPattern + "." + ext
			fmt.Println("scanning glob", extPattern)
			additional, err := filepath.Glob(extPattern)
			if err != nil {
				fmt.Println(err)
				continue
			}
			files = append(files, additional...)
		}
		globPattern += "/*"
	}

	return files, nil
}

func GenerateDirectoryGMCs(config command.GMCConfigDirectoryGenerator) {
	dirPath := filepath.Clean(config.Dir)
	fmt.Println(dirPath)
	files, err := scanGlobFiles(dirPath, config.Extensions)
	if err != nil {
		log.Fatal("Dir scan failed")
	}

	gmcPath := filepath.Join("Groovy", config.Name)

	for _, file := range files {
		fmt.Println(file)
		fileDir := filepath.Dir(file)
		relPath := file[len(dirPath):]
		relDir := filepath.Dir(relPath)

		gmc := GroovyMiSTerCommandNoRaw{
			Cmd:  config.Template.Cmd,
			Vars: config.Template.Vars,
		}

		varsRep := map[string]string{
			"ROM_FULL_PATH":     file,
			"ROM_FULL_DIR":      fileDir,
			"ROM_RELATIVE_PATH": relPath,
			"ROM_RELATIVE_DIR":  relDir,
		}

		for k, v := range gmc.Vars {
			value := v
			for kRep, vRep := range varsRep {
				pattern := "${" + kRep + "}"
				value = strings.Replace(value, pattern, vRep, -1)
			}
			gmc.Vars[k] = value
		}

		jsonBytes, err := json.Marshal(gmc)
		if err != nil {
			fmt.Println("error:", err)
		}
		WriteGMCtoDir(filepath.Join(gmcPath, relDir), filepath.Base(file), jsonBytes)
	}
}
