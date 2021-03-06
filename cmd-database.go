package main

import (
	"github.com/MckayJT/murmur-cli/internal/MurmurRPC"
)

func init() {
	cmd := root.Add("database")

	cmd.Add("query", func(args Args) {
		server := args.MustServer(0)
		query := &MurmurRPC.DatabaseUser_Query{
			Server: server,
		}
		if filter, ok := args.String(1); ok {
			query.Filter = &filter
		}
		Output(client.DatabaseUserQuery(ctx, query))
	})

	cmd.Add("get", func(args Args) {
		server := args.MustServer(0)
		id := args.MustUint32(1)
		Output(client.DatabaseUserGet(ctx, &MurmurRPC.DatabaseUser{
			Server: server,
			Id:     &id,
		}))
	})

	cmd.Add("add", func(args Args) {
		server := args.MustServer(0)
		name := args.MustString(1)
		password := args.MustString(2)
		Output(client.DatabaseUserRegister(ctx, &MurmurRPC.DatabaseUser{
			Server:   server,
			Name:     &name,
			Password: &password,
		}))
	})
}
