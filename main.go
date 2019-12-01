package main

import (
	"context"
	"crypto/tls"
	_ "encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/MckayJT/murmur-cli/internal/MurmurRPC"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	commands []*cli.Command
)

/*

const usage = `murmur-cli provides an interface to a grpc-enabled murmur server.
usage: murmur-cli [flags] [command... [arguments...]]

Flags:
  -address="127.0.0.1:50051"   address and port of murmur's grpc endpoint
                                (can also be set via $MURMUR_ADDRESS).
  -timeout="10s"               duration to wait for connection.
  -template=""                 Go text/template template to use when outputing
                                data. By default, JSON objects are printed.
  -hostoverride=""             Expect this host name from the server
  -cert=""                     Client certificate (pem format)
  -key=""                      Client certificate key (pem format), unencrypted
  -insecure=false              Disable TLS encryption.
  -help                        Print command list.
`

const usageCommands = `
Commands:
  acl get <server id> <channel id>
  acl get-effective-permissions <server id> <session> <channel id>

  ban get <server id>

  channel query <server id>
  channel get <server id> <channel id>
  channel add <server id> <parent channel id> <name>
  channel remove <server id> <channel id>

  config get <server id>
  config get-field <server id> <key>
  config set-field <server id> <key> <value>
  config get-defaults

  contextaction add <server id> <context> <action> <text> <session>
    Context is a comma seperated list of the following:
      Server
      Channel
      User
  contextaction remove <server id> <action> [session]
  contextaction events <server id> <action>

  database query <server id> [filter]
  database get <server id> <user id>
  database add <server id> <user id> <password>

  log query <server id> (<min> <max>)

  meta uptime
  meta version
  meta events

  server create
  server query
  server get <server id>
  server start <server id>
  server stop <server id>
  server remove <server id>
  server events <server id>

  textmessage send <server id> [sender:<session>] [targets...] <text>
    Valid targets:
      user:<session>
      channel:<id>
      tree:<id>
  textmessage filter <server id> <program> [args...]

  tree query <server id>

  user query <server id>
  user get <server id> <session>
  user kick <server id> <session> [reason]
`

*/
var outputTemplate *template.Template
var app *cli.App

func main() {
	app = cli.NewApp()
	app.HelpName = "murmur-cli"
	app.Usage = "manage murmur using gRPC"
	app.Description = "manage murmur using gRPC"
	app.UseShortOptionHandling = false
	app.Authors = []*cli.Author{&cli.Author{
		Name:  "John McKay",
		Email: "",
	},
	}
	configPaths := GetDefaultConfigPaths()
	paths := strings.Join(configPaths, string(os.PathListSeparator))
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "noconfig",
			Usage: "disable configuration file",
			Value: false,
		},
		&cli.PathFlag{
			Name:        "config",
			Usage:       "Load configuration from `FILE`",
			TakesFile:   true,
			DefaultText: paths,
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "address",
			Aliases:     []string{"a"},
			Usage:       "address and port to connect to server",
			EnvVars:     []string{"MURMUR_ADDRESS"},
			DefaultText: "127.0.0.1:50051",
			Value:       "127.0.0.1:50051",
			TakesFile:   false,
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:        "timeout",
			Aliases:     []string{"t"},
			Usage:       "Timeout when connecting to server",
			DefaultText: "10s",
			Value:       time.Second * 10,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "insecure",
			Usage:       "Do not attempt TLS connection",
			DefaultText: "false",
			Value:       false,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "hostoverride",
			Usage:       "Override server certificate hostname check",
			DefaultText: "none",
			Value:       "",
			TakesFile:   false,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "cert",
			Aliases:     []string{"c"},
			Usage:       "PEM client TLS certificate to use",
			DefaultText: "none",
			Value:       "",
			TakesFile:   false,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "key",
			Aliases:     []string{"k"},
			Usage:       "Unencrypted PEM key for client certificate",
			DefaultText: "none",
			Value:       "",
			TakesFile:   false,
		}),
	}
	app.Commands = append(app.Commands, commands...)
	app.Before = finishSetup
	app.Run(os.Args)
}

func finishSetup(ctx *cli.Context) error {
	err := loadSettings(ctx)
	if err != nil {
		return err
	}
	err = setupConnection(ctx)
	if err != nil {
		return err
	}
	return nil
}

func loadSettings(ctx *cli.Context) error {
	var paths []string
	if ctx.Bool("noconfig") {
		return nil
	}

	if len(ctx.String("config")) != 0 {
		paths = append(paths, ctx.String("config"))
	} else {
		paths = GetDefaultConfigPaths()
	}
	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil || !stat.Mode().IsRegular() {
			continue
		}
		configSource, err := altsrc.NewTomlSourceFromFile(path)
		if err != nil {
			continue
		}
		err = altsrc.ApplyInputSourceValues(ctx, configSource, ctx.App.Flags)
		if err != nil {
			continue
		}
		ctx.App.Metadata["configSource"] = &configSource
		return nil
	}
	return nil
}

func setupConnection(ctx *cli.Context) error {
	address := ctx.String("address")
	timeout := ctx.Duration("timeout")
	insecure := ctx.Bool("insecure")
	hostoverride := ctx.String("hostoverride")
	cert := ctx.String("cert")
	key := ctx.String("key")

	// grpc connection
	dCtx, _ := context.WithTimeout(context.Background(), timeout)
	opts := []grpc.DialOption{
		grpc.WithBlock(),
	}
	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		var tlsConfig tls.Config
		if cert != "" && key != "" {
			cert, err := tls.LoadX509KeyPair(cert, key)
			if err != nil {
				fmt.Printf("Error loading certificate: %v\n", err)
				os.Exit(1)
			}
			tlsConfig.Certificates = append(tlsConfig.Certificates, cert)
		}
		creds := credentials.NewTLS(&tlsConfig)
		if hostoverride != "" {
			err := creds.OverrideServerName(hostoverride)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	}
	conn, err := grpc.DialContext(dCtx, address, opts...)
	if err != nil {
		fmt.Printf("Error connecting: %v\n", err)
		fmt.Printf("timeout: %v\n", timeout)
		os.Exit(1)
	}

	clien := MurmurRPC.NewV1Client(conn)
	ctx.App.Metadata["grpcConnection"] = &conn
	ctx.App.Metadata["grpcClient"] = clien
	return nil
	//defer conn.Close()
	/*
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if ok {
					jsonErr := struct {
						Error string `json:"error"`
					}{
						Error: err.Error(),
					}
					json.NewEncoder(os.Stderr).Encode(&jsonErr)
					os.Exit(3)
				}
			}
		}()

		if root.Do() != nil {
			flag.Usage()
			os.Exit(1)
		}
	*/
}
