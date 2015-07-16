package main

import (
	"github.com/layeh/murmur-cli/MurmurRPC"
)

func init() {
	cmd := root.Add("channel")

	cmd.Add("query", func(args Args) {
		server := args.MustServer(0)
		Output(client.QueryChannels(ctx, &MurmurRPC.Channel_Query{
			Server: server,
		}))
	})

	cmd.Add("get", func(args Args) {
		server := args.MustServer(0)
		id := args.MustUint32(1)
		Output(client.GetChannel(ctx, &MurmurRPC.Channel{
			Server: server,
			Id:     &id,
		}))
	})

	cmd.Add("add", func(args Args) {
		server := args.MustServer(0)
		id := args.MustUint32(1)
		name := args.MustString(2)
		Output(client.AddChannel(ctx, &MurmurRPC.Channel{
			Server: server,
			Parent: &MurmurRPC.Channel{
				Id: &id,
			},
			Name: &name,
		}))
	})

	cmd.Add("remove", func(args Args) {
		server := args.MustServer(0)
		id := args.MustUint32(1)
		Output(client.RemoveChannel(ctx, &MurmurRPC.Channel{
			Server: server,
			Id:     &id,
		}))
	})
}
