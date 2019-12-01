package main

import (
	mr "github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/urfave/cli/v2"
	"io"
)

func init() {
	cmd := &cli.Command{
		Name:     "contextaction",
		Usage:    "do actions based on other actions",
		HideHelp: true,
	}
	subs := []*cli.Command{
		&cli.Command{
			Name:      "add",
			Usage:     "add new contextaction",
			ArgsUsage: "<server id> <actionmask> <action> <text> <session id>",
			Action:    doAddAction,
		},
		&cli.Command{
			Name:      "remove",
			Usage:     "remove contextaction",
			ArgsUsage: "<server id> <action> [session]",
			Action:    doRemoveAction,
		},
		&cli.Command{
			Name:      "events",
			Usage:     "listen for contextactions",
			ArgsUsage: "<server id> <action>",
			Action:    doEventAction,
		},
	}

	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doAddAction(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustBitmask, MustString, MustString, MustUint32)
	if err != nil {
		return NewUsageError(ctx, err)
	}

	resp, err := client.ContextActionAdd(mCtx, &mr.ContextAction{
		Server:  args[0].Server(),
		Context: args[1].u32_p(),
		Action:  args[2].s_p(),
		Text:    args[3].s_p(),
		User: &mr.User{
			Session: args[4].u32_p(),
		},
	})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doRemoveAction(ctx *cli.Context) error {
	funcs := []ProcessArgFunc{MustServer, MustString}
	if ctx.NArg() > 2 {
		funcs = append(funcs, MustUint32)
	}
	client, mCtx, args, err := ProcessArguments(ctx, funcs...)
	if err != nil {
		return NewUsageError(ctx, err)
	}

	contextAction := &mr.ContextAction{
		Server: args[0].Server(),
		Action: args[1].s_p(),
	}
	if len(args) > 2 {
		contextAction.User = &mr.User{
			Session: args[2].u32_p(),
		}
	}
	resp, err := client.ContextActionRemove(mCtx, contextAction)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doEventAction(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustString)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	stream, err := client.ContextActionEvents(mCtx, &mr.ContextAction{
		Server: args[0].Server(),
		Action: args[1].s_p(),
	})
	if err != nil {
		return cli.NewExitError(err, 2)
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return cli.NewExitError(err, 2)
			}
			return nil
		}
		err = Output(msg)
		if err != nil {
			return cli.NewExitError(err, 2)
		}
	}
}
