package main

import (
	mr "github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd := &cli.Command{
		Name:     "tree",
		Usage:    "",
		HideHelp: true,
	}
	subs := []*cli.Command{
		&cli.Command{
			Name:      "query",
			Usage:     "get tree view of server",
			ArgsUsage: "<server id>",
			Action:    doTreeQuery,
		},
	}
	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doTreeQuery(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.TreeQuery(mCtx, &mr.Tree_Query{Server: args[0].Server()})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}
