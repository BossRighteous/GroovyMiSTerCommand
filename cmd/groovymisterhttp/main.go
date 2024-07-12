package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BossRighteous/GroovyMiSTerHTTP/pkg/command"
	"github.com/BossRighteous/GroovyMiSTerHTTP/pkg/display"
	"github.com/BossRighteous/GroovyMiSTerHTTP/pkg/server"
)

func main() {
	configPath := "./config.json"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	config, err := command.LoadConfigFromPath(configPath)
	if err != nil {
		log.Fatal("Config could not be loaded from path argument")
	}

	misterHost := config.MisterHost
	cmdChan := make(chan command.GroovyMiSTerCommand)

	cmdr := &command.CommandRunner{
		Config:     config,
		ResultChan: make(chan command.RunResult),
	}

	server.StartUdpClient(misterHost, cmdChan)

	disp := display.NewMiSTerDisplay(misterHost)

	for {
		select {
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
