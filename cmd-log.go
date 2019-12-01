package main

import (
	mr "github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd := &cli.Command{
		Name:     "log",
		Usage:    "manage log files",
		HideHelp: true,
	}

	subs := []*cli.Command{
		&cli.Command{
			Name:      "query",
			Usage:     "query log files",
			ArgsUsage: "<server id> [<min> <max>]",
			Action:    doLogQuery,
		},
	}

	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doLogQuery(ctx *cli.Context) error {
	var funcs []ProcessArgFunc = []ProcessArgFunc{MustServer}
	if ctx.NArg() > 1 {
		funcs = append(funcs, MustUint32, MustUint32)
	}
	client, mCtx, args, err := ProcessArguments(ctx, funcs...)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	query := &mr.Log_Query{Server: args[0].Server()}
	if ctx.NArg() > 1 {
		query.Min = args[1].u32_p()
		query.Max = args[2].u32_p()
	}
	resp, err := client.LogQuery(mCtx, query)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}
