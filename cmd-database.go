package main

import (
	mr "github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd := &cli.Command{
		Name:     "database",
		Usage:    "manage user database",
		HideHelp: true,
	}

	subs := []*cli.Command{
		&cli.Command{
			Name:      "query",
			Usage:     "query user database",
			ArgsUsage: "<server id> [filter]",
			Action:    doDBQuery,
		},
		&cli.Command{
			Name:      "get",
			Usage:     "Get info on one user",
			ArgsUsage: "<serever id> <user id>",
			Action:    doDBGet,
		},
		&cli.Command{
			Name:      "add",
			Usage:     "add user to database",
			ArgsUsage: "<server id> <username> <password>",
			Action:    doDBAdd,
		},
	}
	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doDBQuery(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	query := &mr.DatabaseUser_Query{
		Server: args[0].Server(),
	}
	filter := ctx.Args().Get(1)
	if len(filter) != 0 {
		query.Filter = &filter
	}
	resp, err := client.DatabaseUserQuery(mCtx, query)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doDBGet(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustUint32)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	resp, err := client.DatabaseUserGet(mCtx, &mr.DatabaseUser{
		Server: args[0].Server(),
		Id:     args[1].u32_p(),
	})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doDBAdd(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustString, MustString)
	if err != nil {
		return NewUsageError(ctx, err)
	}

	resp, err := client.DatabaseUserRegister(mCtx, &mr.DatabaseUser{
		Server:   args[0].Server(),
		Name:     args[1].s_p(),
		Password: args[2].s_p(),
	})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}
