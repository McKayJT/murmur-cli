package main

import (
	mr "github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd := &cli.Command{
		Name:     "user",
		Usage:    "manage connected users",
		HideHelp: true,
	}
	subs := []*cli.Command{
		&cli.Command{
			Name:      "query",
			Usage:     "get active users",
			ArgsUsage: "<server id>",
			Action:    doUserQuery,
		},
		&cli.Command{
			Name:      "get",
			Usage:     "get info for active user",
			ArgsUsage: "<server id> <session id>",
			Action:    doUserGet,
		},
		&cli.Command{
			Name:      "kick",
			Usage:     "kick user from server",
			ArgsUsage: "<server id> <session id> [reason]",
			Action:    doUserKick,
		},
	}
	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doUserQuery(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.UserQuery(mCtx, &mr.User_Query{Server: args[0].Server()})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doUserGet(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustUint32)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.UserGet(mCtx, &mr.User{
		Server:  args[0].Server(),
		Session: args[1].u32_p(),
	})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doUserKick(ctx *cli.Context) error {
	funcs := []ProcessArgFunc{MustServer, MustUint32}
	if ctx.NArg() > 2 {
		funcs = append(funcs, MustString)
	}
	client, mCtx, args, err := ProcessArguments(ctx, funcs...)
	if err != nil {
		return NewUsageError(ctx, err)
	}

	kick := &mr.User_Kick{
		Server: args[0].Server(),
		User:   &mr.User{Session: args[1].u32_p()},
	}
	if len(args) > 2 {
		kick.Reason = args[2].s_p()
	}
	resp, err := client.UserKick(mCtx, kick)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}
