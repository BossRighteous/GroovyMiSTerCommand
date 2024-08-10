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

type MameGMCVars struct {
	MachineName string `json:"MACHINE_NAME"`
}

type MameGMC struct {
	Cmd  string      `json:"cmd"`
	Vars MameGMCVars `json:"vars"`
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

func GenerateMameGMCs(config command.GMCConfigMameGenerator) {
	mamelistPath := filepath.Clean(config.MamelistPath)
	mamelist, err := loadMamelist(mamelistPath)
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

	romsDir := filepath.Clean(config.RomsDir)
	for _, item := range *mamelist {
		//fmt.Println(item.Name)
		exists := verifyRomExists(romsDir, item.Name)
		if !exists {
			//fmt.Println(item.Name, "not found")
			continue
		}

		gmc := MameGMC{
			Cmd: "mame",
			Vars: MameGMCVars{
				MachineName: item.Name,
			},
		}
		jsonBytes, err := json.Marshal(gmc)
		if err != nil {
			fmt.Println("error:", err)
		}

		WriteGMCtoDir(allPath, item.Description, jsonBytes)
		WriteGMCtoDir(filepath.Join(genresPath, item.Genre), item.Description, jsonBytes)
		WriteGMCtoDir(filepath.Join(manusPath, item.Manufacturer), item.Description, jsonBytes)
		WriteGMCtoDir(filepath.Join(yearsPath, item.Year), item.Description, jsonBytes)
	}
}
