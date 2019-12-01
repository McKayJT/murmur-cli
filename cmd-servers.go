package main

import (
	"io"

	mr "github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd := &cli.Command{
		Name:     "server",
		Usage:    "Manage servers",
		HideHelp: true,
	}
	subs := []*cli.Command{
		&cli.Command{
			Name:   "create",
			Usage:  "create new server",
			Action: doCreateServer,
		},
		&cli.Command{
			Name:   "query",
			Usage:  "query servers",
			Action: doQueryServer,
		},
		&cli.Command{
			Name:      "get",
			Usage:     "get server information",
			ArgsUsage: "<server id>",
			Action:    doGetServer,
		},
		&cli.Command{
			Name:      "start",
			Usage:     "start server",
			ArgsUsage: "<server id>",
			Action:    doStartServer,
		},
		&cli.Command{
			Name:      "stop",
			Usage:     "stop server",
			ArgsUsage: "<server id>",
			Action:    doStopServer,
		},
		&cli.Command{
			Name:      "remove",
			Usage:     "remove server",
			ArgsUsage: "<server id>",
			Action:    doRemoveServer,
		},
		&cli.Command{
			Name:      "events",
			Usage:     "listen for server events",
			ArgsUsage: "<server id>",
			Action:    doEventsServer,
		},
	}
	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doCreateServer(ctx *cli.Context) error {
	client, mCtx, _, err := ProcessArguments(ctx)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.ServerCreate(mCtx, RPCVoid)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doQueryServer(ctx *cli.Context) error {
	client, mCtx, _, err := ProcessArguments(ctx)
	if err != nil {
		return NewUsageError(ctx, err)
	}

	resp, err := client.ServerQuery(mCtx, &mr.Server_Query{})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doGetServer(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.ServerGet(mCtx, args[0].Server())
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doStartServer(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}

	resp, err := client.ServerStart(mCtx, args[0].Server())
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doStopServer(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.ServerStop(mCtx, args[0].Server())
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doRemoveServer(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.ServerRemove(mCtx, args[0].Server())
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doEventsServer(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	stream, err := client.ServerEvents(mCtx, args[0].Server())
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return cli.NewExitError(err, 1)
			}
			return nil
		}
		err = Output(msg)
		if err != nil {
			return cli.NewExitError(err, 1)
		}
	}
}
