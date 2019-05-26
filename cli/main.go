package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/chzyer/readline"
	"github.com/urfave/cli"
	"github.com/winded/tyomaa/cli/commands"
	"github.com/winded/tyomaa/cli/settings"
	"github.com/winded/tyomaa/shared/api/client"
	"github.com/winded/tyomaa/shared/api/client/standard_http_client"
)

func saveSettings(s settings.Settings) {
	settings.Save(s)
}

func main() {
	debug := false
	if idbg, err := strconv.Atoi(os.Getenv("DEBUG")); err != nil && idbg > 0 {
		debug = true
	}

	cliSettings, err := settings.Load()
	if err != nil {
		panic(err)
	}

	apiClient := client.NewApiClient(standard_http_client.New(time.Second*10), cliSettings.Api)
	inputReader, err := readline.New("")
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "tyomaa-cli"
	app.Usage = "CLI interface for tyomaa"
	app.Commands = commands.GetCommands(cliSettings, apiClient, inputReader, saveSettings)

	err = app.Run(os.Args)
	if err != nil {
		if debug {
			log.Fatal(err)
		} else {
			fmt.Println("ERROR: " + err.Error())
		}
	}
}
