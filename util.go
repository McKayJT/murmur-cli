package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/MckayJT/murmur-cli/internal/MurmurRPC"
	"github.com/golang/protobuf/proto"
	"github.com/urfave/cli/v2"
)

var RPCVoid *MurmurRPC.Void = &MurmurRPC.Void{}

func NewUsageError(ctx *cli.Context, msg error) error {
	var out strings.Builder
	c := ctx.Command
	if len(c.HelpName) == 0 {
		c.HelpName = fmt.Sprintf("%s %s", ctx.App.Name, c.FullName())
	}
	cli.HelpPrinterCustom(&out, cli.CommandHelpTemplate, c, nil)
	fmt.Fprintf(&out, "\n%v\n", msg)
	return cli.NewExitError(out.String(), 1)
}

type RPCArg = interface {
	Server() *MurmurRPC.Server
	u32() uint32
	u32_p() *uint32
	i32() int32
	i32_p() *int32
	s() string
	s_p() *string
}

type rpcarg struct{ d interface{} }

func (a *rpcarg) Server() *MurmurRPC.Server {
	return a.d.(*MurmurRPC.Server)
}

func (a *rpcarg) u32() uint32 {
	return a.d.(uint32)
}

func (a *rpcarg) u32_p() *uint32 {
	ret := a.d.(uint32)
	return &ret
}

func (a *rpcarg) i32() int32 {
	return a.d.(int32)
}

func (a *rpcarg) i32_p() *int32 {
	ret := a.d.(int32)
	return &ret
}

func (a *rpcarg) s() string {
	return a.d.(string)
}

func (a *rpcarg) s_p() *string {
	ret := a.d.(string)
	return &ret
}

func newRPCArg(data interface{}) RPCArg {
	return &rpcarg{d: data}
}

type ProcessArgFunc func(string) (RPCArg, error)

func Output(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent(" ", "   ")
	if err := encoder.Encode(data); err != nil {
		return cli.Exit(err, 1)
	}
	return nil
}

func MustServer(arg string) (RPCArg, error) {
	if len(arg) <= 0 {
		return nil, errors.New("missing server ID argument")
	}
	id, err := strconv.Atoi(arg)
	if err != nil {
		return nil, errors.New("invalid server ID")
	}
	return newRPCArg(&MurmurRPC.Server{Id: proto.Uint32(uint32(id))}), nil
}

func MustString(arg string) (RPCArg, error) {
	if len(arg) == 0 {
		return nil, errors.New("missing string value")
	}
	return newRPCArg(arg), nil
}

func MustUint32(arg string) (RPCArg, error) {
	if len(arg) == 0 {
		return nil, errors.New("missing uint32 value")
	}
	n, err := strconv.ParseUint(arg, 10, 32)
	if err != nil {
		return nil, err
	}
	return newRPCArg(uint32(n)), nil
}

//func (a Args) MustBitmask(i int, values map[string]int32, allowEmpty bool) int32 {
func MustBitmask(arg string) (RPCArg, error) {
	var val int32
	for _, item := range strings.Split(arg, ",") {
		itemVal, ok := MurmurRPC.ContextAction_Context_value[item]
		if !ok {
			return nil, errors.New("invalid bitmask value")
		}
		val |= itemVal
	}
	if val == 0 {
		return nil, errors.New("empty bitmask value")
	}
	return newRPCArg(uint32(val)), nil
}

func ProcessArguments(ctx *cli.Context, funcs ...ProcessArgFunc) (MurmurRPC.V1Client, context.Context, []RPCArg, error) {
	var rets []RPCArg
	client := ctx.App.Metadata["grpcClient"].(MurmurRPC.V1Client)
	if ctx.NArg() < len(funcs) {
		return nil, nil, nil, fmt.Errorf("Invalid argument count: expected %d, got %d", len(funcs), ctx.NArg())
	}
	for i, f := range funcs {
		ret, err := f(ctx.Args().Get(i))
		if err != nil {
			return nil, nil, nil, err
		}
		rets = append(rets, ret)
	}
	return client, context.Background(), rets, nil
}

func GetDefaultConfigPaths() []string {
	var paths []string
	XDG_CONFIG_HOME, err := os.UserConfigDir()
	if err == nil {
		p, err := filepath.Abs(filepath.Join(XDG_CONFIG_HOME, "murmur-cli", "murmur-cli.toml"))
		if err == nil {
			paths = append(paths, p)
		}
	}

	HOME, err := os.UserHomeDir()
	if err == nil {
		p, err := filepath.Abs(filepath.Join(HOME, ".murmur-cli.toml"))
		if err == nil {
			paths = append(paths, p)
		}
	}
	return paths
}
