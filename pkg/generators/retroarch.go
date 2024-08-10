package generators

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/BossRighteous/GroovyMiSTerCommand/pkg/command"
)

type RetroarchPlaylistItem struct {
	Path     string `json:"path"`
	Label    string `json:"label"`
	CorePath string `json:"core_path"`
	CoreName string `json:"core_name"`
	Crc32    string `json:"crc32"`
	DbName   string `json:"db_name"`
}

type RetroarchPlaylist struct {
	Version            string                  `json:"version"`
	DefaultCorePath    string                  `json:"default_core_path"`
	DefaultCoreName    string                  `json:"default_core_name"`
	LabelDisplayMode   int                     `json:"label_display_mode"`
	RightThumbnailMode int                     `json:"right_thumbnail_mode"`
	LeftThumbnailMode  int                     `json:"left_thumbnail_mode"`
	ThumbnailMatchMode int                     `json:"thumbnail_match_mode"`
	SortMode           int                     `json:"sort_mode"`
	Items              []RetroarchPlaylistItem `json:"items"`
}

type RetroarchGMCVars struct {
	CorePath string `json:"CORE_PATH"`
	RomPath  string `json:"ROM_PATH"`
}

type RetroarchGMC struct {
	Cmd  string           `json:"cmd"`
	Vars RetroarchGMCVars `json:"vars"`
}

func loadRetroarchPlaylist(path string) (*RetroarchPlaylist, error) {
	var playlist RetroarchPlaylist
	dat, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dat, &playlist)
	if err != nil {
		return nil, err
	}
	return &playlist, nil
}

func scanLPLFiles(playlistDir string) ([]string, error) {
	return filepath.Glob(playlistDir + "/*.lpl")
}

func GenerateRetroarchGMCs(config command.GMCConfigRetroarchGenerator) {
	playlistDir := filepath.Clean(config.PlaylistsDir)
	playlistPaths, err := scanLPLFiles(playlistDir)
	if err != nil {
		log.Fatal(err)
	}
	gmcPath := filepath.Join("Groovy", "RetroArch")

	for _, playlistPath := range playlistPaths {
		fmt.Println("scanning", playlistPath)
		playlist, err := loadRetroarchPlaylist(playlistPath)
		playlistName, ok := GetBaseFilename(playlistPath)
		if !ok {
			playlistName = filepath.Base(playlistPath)
		}
		if err != nil {
			fmt.Println("skipping playlist: ", err)
			continue
		}
		if playlist.DefaultCorePath == "" {
			fmt.Println("playlist has no default_core_path, skipping")
			continue
		}

		for _, item := range playlist.Items {
			gmcCorePath := playlist.DefaultCorePath
			if item.CorePath != "" && item.CorePath != "DETECT" {
				gmcCorePath = item.CorePath
			}

			gmc := RetroarchGMC{
				Cmd: "retroarch",
				Vars: RetroarchGMCVars{
					CorePath: gmcCorePath,
					RomPath:  item.Path,
				},
			}
			jsonBytes, err := json.Marshal(gmc)
			if err != nil {
				fmt.Println("error:", err)
			}
			WriteGMCtoDir(filepath.Join(gmcPath, playlistName), item.Label, jsonBytes)
		}
	}
}
