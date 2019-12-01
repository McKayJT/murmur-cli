package main

import (
	"github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd := &cli.Command{
		Name:      "ban",
		Usage:     "view serever ban list",
		ArgsUsage: "<server>",
		HideHelp:  true,
		Action:    doBanGet,
	}
	commands = append(commands, cmd)
}

func doBanGet(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	ret, err := client.BansGet(mCtx, &MurmurRPC.Ban_Query{Server: args[0].Server()})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(ret)
}
