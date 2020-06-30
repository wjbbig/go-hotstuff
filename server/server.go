package main

import (
	"go-hotstuff/hotstuff/basic"
	"go-hotstuff/proto"
	"google.golang.org/grpc"
	"net"
)

func main() {
	rpcServer := grpc.NewServer()
	proto.RegisterBasicHotStuffServer(rpcServer, new(basic.BasicHotStuffService))
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	rpcServer.Serve(listen)
}


