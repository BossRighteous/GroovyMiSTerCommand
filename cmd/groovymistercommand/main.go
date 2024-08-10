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
		log.Fatal("Config could not be loaded from path argument")
	}

	if runType == "server" {
		runServer(config)
	} else if runType == "generate:retroarch" {
		runGenerateRetroArch(config)
	} else if runType == "generate:mame" {
		runGenerateMAME(config)
	} else if runType == "generate:directory" {
		runGenerateMednafen(config)
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

func runGenerateMednafen(_ *command.GMCConfig) {
	fmt.Println("generate:mednafen not yet implemented")
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

	disp := display.NewMiSTerDisplay(misterHost)

	for {
		select {
		case <-beaconTicker.C:
			if !cmdr.IsRunning() {
				server.SendBeacon()
			}
		case res := <-cmdr.ResultChan:
			fmt.Println(res)
			if res.BlitMessage {
				if len(res.MessageLines) > 0 {
					disp.BlitText(res.MessageLines)
				} else if res.Message != "" {
					disp.BlitText(display.ReflowText(res.Message))
				} else {
					disp.BlitText([]string{"Process Completed"})
				}
			}
		case cmd := <-cmdChan:
			disp.SafeClose()
			fmt.Println(cmd.Raw)
			res := cmdr.Run(cmd)
			if res.BlitMessage {
				if len(res.MessageLines) > 0 {
					disp.BlitText(res.MessageLines)
				} else {
					disp.BlitText(display.ReflowText(res.Message))
				}
			}
		}
	}
}
