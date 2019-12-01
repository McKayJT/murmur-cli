package main

import mr "github.com/MckayJT/murmur-cli/internal/MurmurRPC"
import "github.com/urfave/cli/v2"

func init() {
	cmd := &cli.Command{
		Name:                   "acl",
		Usage:                  "manage and view acl permissions",
		ArgsUsage:              "[get|get-effective-permissions]",
		HideHelp:               true,
		UseShortOptionHandling: false,
	}

	subs := []*cli.Command{
		&cli.Command{
			Name:      "get",
			Usage:     "get ACL list for a channel",
			ArgsUsage: "<server id> <channel id>",
			Action:    doACLGet,
		},
		&cli.Command{
			Name:      "get-effective-permissions",
			Usage:     "get effective permissions for a user",
			ArgsUsage: "<server id> <sessions> <channel id>",
			Action:    doEffectiveACL,
		},
	}
	cmd.Subcommands = subs
	commands = append(commands, cmd)
}

func doACLGet(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustUint32)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	ret, err := client.ACLGet(mCtx, &mr.Channel{
		Server: args[0].Server(),
		Id:     args[1].u32_p(),
	})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(ret)
}

func doEffectiveACL(ctx *cli.Context) error {
	client, mCtx, args, err := ProcessArguments(ctx, MustServer, MustUint32, MustUint32)
	if err != nil {
		return NewUsageError(ctx, err)
	}
	ret, err := client.ACLGetEffectivePermissions(mCtx,
		&mr.ACL_Query{
			Server: args[0].Server(),
			User: &mr.User{
				Session: args[1].u32_p(),
			},
			Channel: &mr.Channel{
				Id: args[2].u32_p(),
			},
		})
	if err != nil {
		return cli.NewExitError(err, 1)
	}
	return Output(ret)
}
