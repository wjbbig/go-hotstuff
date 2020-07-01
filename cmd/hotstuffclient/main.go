package main

import (
	"context"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewBasicHotStuffClient(conn)

	_, err = client.SendRequest(context.Background(), &pb.Msg{Payload: &pb.Msg_Request{Request: &pb.Request{
		Cmd:           "1+2",
		ClientAddress: "localhost:9999",
	}}})
	if err != nil {
		panic(err)
	}
}
