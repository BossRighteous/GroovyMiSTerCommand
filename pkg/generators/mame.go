package generators

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BossRighteous/GroovyMiSTerCommand/pkg/command"
)

type MameListGame struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	CloneOf      string `json:"cloneof"`
	Manufacturer string `json:"manufacturer"`
	Year         string `json:"year"`
	Genre        string `json:"genre"`
}

type MameList []MameListGame

type MameCommandVars struct {
	MachineName string `json:"MACHINE_NAME"`
}

type MameCmd struct {
	Cmd  string          `json:"cmd"`
	Vars MameCommandVars `json:"vars"`
}

func loadMamelist(path string) (*MameList, error) {
	var mamelist MameList
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dat, &mamelist)
	if err != nil {
		return nil, err
	}
	return &mamelist, nil
}

func verifyRomExists(romsDir string, name string) bool {
	path := filepath.Join(romsDir, name+".zip")
	if _, err := os.Stat(path); err == nil {
		return true
	}
	return false
}

func writeGMCtoDir(dir string, filename string, content []byte) {
	os.MkdirAll(dir, os.ModePerm)
	gmcPath := filepath.Join(dir, filename+".gmc")
	fo, err := os.Create(gmcPath)
	if err != nil {
		fmt.Printf("Unable to create file, may exist %s\n", gmcPath)
		return
	}
	length, err := fo.Write(content)
	if err != nil || length != len(content) {
		log.Fatal(err)
	}
	if err := fo.Close(); err != nil {
		log.Fatal(err)
	}
}

func GenerateMameGMCs(config command.GMCConfigMameGenerator) {
	mamelist, err := loadMamelist(config.MamelistPath)
	if err != nil {
		fmt.Println("Mamelist could not be loaded from JSON path")
		log.Fatal(err)
	}

	// Allocate folder
	gmcPath := filepath.Join("Groovy", "MAME")
	os.MkdirAll(gmcPath, os.ModePerm)

	allPath := filepath.Join(gmcPath, "All")
	os.MkdirAll(allPath, os.ModePerm)

	genresPath := filepath.Join(gmcPath, "Genres")
	os.MkdirAll(genresPath, os.ModePerm)

	manusPath := filepath.Join(gmcPath, "Manufacturers")
	os.MkdirAll(manusPath, os.ModePerm)

	yearsPath := filepath.Join(gmcPath, "Years")
	os.MkdirAll(yearsPath, os.ModePerm)

	for _, item := range *mamelist {
		//fmt.Println(item.Name)
		exists := verifyRomExists(config.RomsDir, item.Name)
		if !exists {
			//fmt.Println(item.Name, "not found")
			continue
		}

		cmd := MameCmd{
			Cmd: "mame",
			Vars: MameCommandVars{
				MachineName: item.Name,
			},
		}
		jsonBytes, err := json.Marshal(cmd)
		if err != nil {
			fmt.Println("error:", err)
		}

		writeGMCtoDir(allPath, item.Description, jsonBytes)
		writeGMCtoDir(filepath.Join(genresPath, item.Genre), item.Description, jsonBytes)
		writeGMCtoDir(filepath.Join(manusPath, item.Manufacturer), item.Description, jsonBytes)
		writeGMCtoDir(filepath.Join(yearsPath, item.Year), item.Description, jsonBytes)
	}
}
