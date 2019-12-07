package main

import (
	"github.com/urfave/cli/v2"
	"os"
)

var (
	commands []*cli.Command
	app *cli.App
)

func main() {
	app = CreateApp()
	defer CloseConnection()
	app.Run(os.Args)
}
