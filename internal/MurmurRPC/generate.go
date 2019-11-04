// +build ignore

//go:generate env GOBIN=$PWD go get -u github.com/golang/protobuf/protoc-gen-go
//go:generate curl -LO https://raw.githubusercontent.com/mumble-voip/mumble/master/src/murmur/MurmurRPC.proto
//go:generate protoc --plugin=protoc-gen-go --go_out=plugins=grpc:. MurmurRPC.proto
//go:generate rm protoc-gen-go

package MurmurRPC
