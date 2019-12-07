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
	"time"

	"github.com/MckayJT/murmur-cli/internal/MurmurRPC"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func CreateApp() *cli.App {
	a := cli.NewApp()
	a.HelpName = "murmur-cli"
	a.Usage = "manage murmur using gRPC"
	a.Description = "manage murmur using gRPC"
	a.UseShortOptionHandling = false
	a.Authors = []*cli.Author{&cli.Author{
		Name:  "John McKay",
		Email: "",
	},
	}
	configPaths := GetDefaultConfigPaths()
	paths := strings.Join(configPaths, string(os.PathListSeparator))
	a.Flags = []cli.Flag{
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
	a.Commands = append(a.Commands, commands...)
	a.Before = finishSetup
	cli.OsExiter = exitCleanly
	return a
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
	dCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
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
	ctx.App.Metadata["grpcConnection"] = conn
	ctx.App.Metadata["grpcClient"] = clien
	return nil
}

func CloseConnection() {
	if app != nil {
		conn, ok := app.Metadata["grpcConnection"]
		if ok {
			conn.(*grpc.ClientConn).Close()
		}
	}
}

func exitCleanly(code int) {
	CloseConnection()
	os.Exit(code)
}
