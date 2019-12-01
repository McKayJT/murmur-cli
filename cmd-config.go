package main

import (
	mr "github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd := &cli.Command{
		Name:      "config",
		Usage:     "change murmur configuration options",
		ArgsUsage: "[get|get-field|set-field]",
		HideHelp:  true,
	}
	subs := []*cli.Command{
		&cli.Command{
			Name:      "get",
			Usage:     "get server configuration options",
			ArgsUsage: "<server id>",
			Action:    doGet,
		},
		&cli.Command{
			Name:      "get-field",
			Usage:     "get value for configuration option",
			ArgsUsage: "<serverId> <key>",
			Action:    doGetField,
		},
		&cli.Command{
			Name:      "set-field",
			Usage:     "set configuration option",
			ArgsUsage: "<server id> <key> <value>",
			Action:    doSetField,
		},
		&cli.Command{
			Name:      "get-default",
			Usage:     "get default configuration",
			ArgsUsage: "",
			Action:    doGetDefault,
		},
	}
	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doGet(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.ConfigGet(mCtx, args[0].Server())

	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doGetField(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustString)
	if err != nil {
		return NewUsageError(ctx, err)
	}

	resp, err := client.ConfigGetField(mCtx, &mr.Config_Field{
		Server: args[0].Server(),
		Key:    args[1].s_p(),
	})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doSetField(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustString, MustString)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.ConfigSetField(mCtx, &mr.Config_Field{
		Server: args[0].Server(),
		Key:    args[1].s_p(),
		Value:  args[2].s_p(),
	})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doGetDefault(ctx *cli.Context) error {
	client, mCtx, _, err := ProcessArguments(ctx)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.ConfigGetDefault(mCtx, RPCVoid)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}
