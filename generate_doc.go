// +build generate

//go:generate go build -tags generate -o docs/makedoc

package main

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	commands []*cli.Command
	app      *cli.App
)

func main() {
	app = CreateApp()
	app.Before = nil
	app.Name = "murmur-cli"
	app.HelpName = "murmur-cli"
	app.Setup()

	_, err := os.Stat("murmur-cli.fish")
	if err != nil {
		fish, err := os.Create("murmur-cli.fish")
		if err != nil {
			panic("could not create murmur-cli.fish")
		}
		completion, err := app.ToFishCompletion()
		if err != nil {
			panic("Could not generate fish completion!")
		}
		_, err = fish.WriteString(completion)
		if err != nil {
			panic("Could not write fish completion file")
		}
		fish.Close()
	} else {
		fmt.Println("murmur-cli.fish already exists; I will not overwrite")
	}

	_, err = os.Stat("murmur-cli.md")
	if err != nil {
		md, err := os.Create("murmur-cli.md")
		if err != nil {
			panic("could not create murmur-cli.md")
		}
		markdown, err := app.ToMarkdown()
		if err != nil {
			panic("Could not generate markdown!")
		}
		_, err = md.WriteString(markdown)
		if err != nil {
			panic("Could not write markdown file")
		}
		md.Close()

	} else {
		fmt.Println("murmur-cli.md already exists; I will not overwrite")
	}
}
