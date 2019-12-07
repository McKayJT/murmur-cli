// +build !generate

package main

import (
	"github.com/urfave/cli/v2"
	"os"
	"context"
)

var (
	commands []*cli.Command
	app *cli.App
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		defer CloseConnection()
		defer cancel()
		interrupts := CreateInterruptChannel()
		<-interrupts
	}()
	app = CreateApp()
	defer CloseConnection()
	app.RunContext(ctx, os.Args)
}
