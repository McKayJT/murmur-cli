package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"

	mr "github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	pb "github.com/golang/protobuf/proto"
	"github.com/urfave/cli/v2"
)

func init() {

	cmd := &cli.Command{
		Name:            "textmessage",
		Usage:           "send or filter text messages",
		SkipFlagParsing: false,
		HideHelp:        true,
	}

	subs := []*cli.Command{
		&cli.Command{
			Name:      "send",
			Usage:     "send text messages",
			ArgsUsage: "<server id> <textual message>",
			Flags: []cli.Flag{
				&cli.UintFlag{
					Name:     "sender",
					Usage:    "sends from user id",
					Required: false,
				},
				&cli.IntSliceFlag{
					Name:     "user",
					Usage:    "list of users to send to",
					Required: false,
				},
				&cli.IntSliceFlag{
					Name:     "channel",
					Usage:    "list of channels to send to",
					Required: false,
				},
				&cli.IntSliceFlag{
					Name:     "tree",
					Usage:    "list of trees to blaze with",
					Required: false,
				},
			},
			Action: doSendText,
		},
		&cli.Command{
			Name:      "filter",
			Usage:     "run an arbitrary command(????) to filter messages",
			ArgsUsage: "<server id> <executable> [argument list]",
			Action:    doFilterText,
		},
	}
	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doSendText(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	messages := ctx.Args().Tail()
	if messages == nil {
		return NewUsageError(ctx, errors.New("You need to specify the message"))
	}
	message := strings.Join(messages, " ")
	tm := &mr.TextMessage{
		Server: args[0].Server(),
		Text:   &message,
	}
	if ctx.IsSet("sender") {
		session := uint32(ctx.Uint("sender"))
		tm.Actor = &mr.User{
			Server:  args[0].Server(),
			Session: pb.Uint32(session),
		}
	}
	if ctx.IsSet("user") {
		users := ctx.IntSlice("user")
		for _, u := range users {
			if u < 0 {
				return NewUsageError(ctx, errors.New("Unsigned integers only!"))
			}
			tm.Users = append(tm.Users, &mr.User{
				Server:  args[0].Server(),
				Session: pb.Uint32(uint32(u)),
			})
		}
	}
	if ctx.IsSet("channel") {
		channels := ctx.IntSlice("channel")
		for _, c := range channels {
			if c < 0 {
				return NewUsageError(ctx, errors.New("Unsigned integers only!"))
			}
			tm.Channels = append(tm.Channels, &mr.Channel{
				Server: args[0].Server(),
				Id:     pb.Uint32(uint32(c)),
			})
		}
	}
	if ctx.IsSet("tree") {
		trees := ctx.IntSlice("tree")
		for _, t := range trees {
			if t < 0 {
				return NewUsageError(ctx, errors.New("Unsigned integers only!"))
			}
			tm.Trees = append(tm.Trees, &mr.Channel{
				Server: args[0].Server(),
				Id:     pb.Uint32(uint32(t)),
			})
		}
	}
	resp, err := client.TextMessageSend(mCtx, tm)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(resp)
}

func doFilterText(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustString)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	arguments := ctx.Args().Tail()[1:]

	stream, err := client.TextMessageFilter(mCtx)
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	msg := mr.TextMessage_Filter{
		Server: args[0].Server(),
	}
	if err := stream.Send(&msg); err != nil {
		return cli.NewExitError(err, 1)
	}

	for {
		req, err := stream.Recv()
		if err != nil {
			if err != io.EOF {
				return cli.NewExitError(err, 1)
			}
			return nil
		}
		var resp mr.TextMessage_Filter
		cmd := exec.Command(args[1].s(), arguments...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		pipe, err := cmd.StdinPipe()
		if err == nil {
			encoder := json.NewEncoder(pipe)
			go encoder.Encode(req.Message)
			cmd.Run()
			if cmd.ProcessState != nil {
				if !cmd.ProcessState.Success() {
					resp.Action = mr.TextMessage_Filter_Reject.Enum()
				}
			}
		}
		stream.Send(&resp)
	}
	return nil
}
