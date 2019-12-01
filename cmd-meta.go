package main

import (
	"github.com/urfave/cli/v2"
	"io"
)

func init() {
	cmd := &cli.Command{
		Name:      "meta",
		Usage:     "[uptime|version|events]",
		ArgsUsage: "Gets metadata from murmur",
		HideHelp:  true,
	}

	subs := []*cli.Command{
		&cli.Command{
			Name:   "uptime",
			Usage:  "gets server uptime",
			Action: doUptime,
		},
		&cli.Command{
			Name:   "version",
			Usage:  "gets murmur version",
			Action: doVersion,
		},
		&cli.Command{
			Name:   "events",
			Usage:  "listen for meta events",
			Action: doEvents,
		},
	}
	cmd.Subcommands = subs

	commands = append(commands, cmd)
}

func doUptime(ctx *cli.Context) error {
	client, mCtx, _, err := ProcessArguments(ctx)
	resp, err := client.GetUptime(mCtx, RPCVoid)
	if err != nil {
		return err
	}
	return Output(resp)
}

func doVersion(ctx *cli.Context) error {
	client, mCtx, _, err := ProcessArguments(ctx)
	resp, err := client.GetVersion(mCtx, RPCVoid)
	if err != nil {
		return err
	}
	return Output(resp)
}

func doEvents(ctx *cli.Context) error {
	client, mCtx, _, err := ProcessArguments(ctx)
	stream, err := client.Events(mCtx, RPCVoid)
	if err != nil {
		return err
	}
	for {
		msg, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		err = Output(msg)
		if err != nil {
			return err
		}
	}
	return nil
}
