package main

import (
	"github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/urfave/cli/v2"
)

func init() {
	cmd := &cli.Command{
		Name:     "channel",
		Usage:    "manage channels",
		HideHelp: true,
	}
	subs := []*cli.Command{
		&cli.Command{
			Name:      "query",
			Usage:     "get list of channels",
			ArgsUsage: "<server>",
			Action:    doChannelQuery,
		},
		&cli.Command{
			Name:      "get",
			Usage:     "get channel information",
			ArgsUsage: "<serverid> <channelid>",
			Action:    doGetChannel,
		},
		&cli.Command{
			Name:      "add",
			UsageText: "add channel to serever",
			ArgsUsage: "<server> <parentid> <name>",
			Action:    doAddChannel,
		},
		&cli.Command{
			Name:      "remove",
			UsageText: "removes channel from server",
			ArgsUsage: "<server> <channelid>",
			Action:    doRemoveChannel,
		},
	}
	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doChannelQuery(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	ret, err := client.ChannelQuery(mCtx, &MurmurRPC.Channel_Query{Server: args[0].Server()})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(ret)
}

func doGetChannel(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustUint32)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	ret, err := client.ChannelGet(mCtx, &MurmurRPC.Channel{
		Server: args[0].Server(),
		Id:     args[1].u32_p(),
	})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(ret)
}

func doAddChannel(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustUint32, MustString)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	ret, err := client.ChannelAdd(mCtx, &MurmurRPC.Channel{
		Server: args[0].Server(),
		Parent: &MurmurRPC.Channel{
			Id: args[1].u32_p(),
		},
		Name: args[2].s_p(),
	})
	if err != nil {
		return err
	}
	return Output(ret)
}

func doRemoveChannel(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustUint32)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	ret, err := client.ChannelRemove(mCtx, &MurmurRPC.Channel{
		Server: args[0].Server(),
		Id:     args[1].u32_p(),
	})
	return Output(ret)
}
