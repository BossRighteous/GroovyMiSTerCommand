package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BossRighteous/GroovyMiSTerCommand/pkg/command"
	"github.com/BossRighteous/GroovyMiSTerCommand/pkg/display"
	"github.com/BossRighteous/GroovyMiSTerCommand/pkg/generators"
	"github.com/BossRighteous/GroovyMiSTerCommand/pkg/server"
)

func main() {
	runType := "server"

	if len(os.Args) > 1 {
		//arg one is program name
		runType = os.Args[1]
	} else {
		fmt.Println("No explicit subcommand provided, running 'server'")
	}

	configPath := "./config.json"

	config, err := command.LoadConfigFromPath(configPath)
	if err != nil {
		fmt.Println("Config could not be loaded from path argument")
		time.Sleep(time.Second * 3)
		return
	}

	if runType == "server" {
		runServer(config)
	} else if runType == "generate:retroarch" {
		runGenerateRetroArch(config)
	} else if runType == "generate:mame" {
		runGenerateMAME(config)
	} else if runType == "generate:directory" {
		if len(os.Args) <= 2 {
			fmt.Println("Missing directory name")
			return
		}
		runGenerateDirectory(config, os.Args[2])
	} else {
		fmt.Println("Unknown Command, exiting")
	}
}

func runGenerateRetroArch(config *command.GMCConfig) {
	generators.GenerateRetroarchGMCs(config.Generators.Retroarch)
}

func runGenerateMAME(config *command.GMCConfig) {
	generators.GenerateMameGMCs(config.Generators.Mame)
}

func runGenerateDirectory(config *command.GMCConfig, name string) {
	found := false
	for _, dir := range config.Generators.Directories {
		if dir.Name == name {
			found = true
			generators.GenerateDirectoryGMCs(dir)
		}
	}
	if !found {
		log.Fatal("No matching directory config found")
	}
}

func runServer(config *command.GMCConfig) {

	misterHost := config.MisterHost
	cmdChan := make(chan command.GroovyMiSTerCommand)

	cmdr := &command.CommandRunner{
		Config:     config,
		ResultChan: make(chan command.RunResult),
	}

	server := server.StartUdpClient(misterHost, cmdChan)
	freq, _ := time.ParseDuration("2s")
	beaconTicker := time.NewTicker(time.Duration(freq))

	var disp *display.MiSTerDisplay
	if config.DisplayMessages {
		disp = display.NewMiSTerDisplay(misterHost)
	}

	for {
		select {
		case <-beaconTicker.C:
			if !cmdr.IsRunning() {
				server.SendBeacon()
			}
		case res := <-cmdr.ResultChan:
			fmt.Println("Process Result: ", res)
			if res.BlitMessage && disp != nil {
				if len(res.MessageLines) > 0 {
					disp.BlitText(res.MessageLines)
				} else if res.Message != "" {
					disp.BlitText(display.ReflowText(res.Message))
				} else {
					disp.BlitText([]string{"Process Completed"})
				}
			}
		case cmd := <-cmdChan:
			if disp != nil {
				disp.SafeClose()
			}
			//fmt.Println(cmd.Raw)
			res := cmdr.Run(cmd)
			fmt.Println("Process Result: ", res)
			if res.BlitMessage && disp != nil {
				if len(res.MessageLines) > 0 {
					disp.BlitText(res.MessageLines)
				} else {
					disp.BlitText(display.ReflowText(res.Message))
				}
			}
		}
	}
}
