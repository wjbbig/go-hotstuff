package main

import (
	"flag"
	"github.com/wjbbig/go-hotstuff/hotstuff/basic"
	"github.com/wjbbig/go-hotstuff/logging"
	"github.com/wjbbig/go-hotstuff/proto"
	"google.golang.org/grpc"
	"net"
	"strings"
)

var id int
var networkType string
var logger = logging.GetLogger()
func init() {
	flag.IntVar(&id, "id", 0, "node id")
	flag.StringVar(&networkType, "type", "basic", "which type of network you want to create.  basic/chained/event-driven")
}

func main() {
	flag.Parse()
	if id <= 0 {
		flag.Usage()
		return
	}
	rpcServer := grpc.NewServer()
	// TODO should use factory method to create one of hotstuff types
	basicHotStuffService := basic.NewBasicHotStuffService(id, networkType)
	proto.RegisterBasicHotStuffServer(rpcServer, basicHotStuffService)

	info := basicHotStuffService.BasicHotStuff.GetSelfInfo()
	port := info.Address[strings.Index(info.Address,":"):]

	logger.Infof("[HOTSTUFF] Server start at port%s", port)
	listen, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	rpcServer.Serve(listen)
}


