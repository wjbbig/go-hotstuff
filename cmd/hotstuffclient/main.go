package main

import (
	"context"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"google.golang.org/grpc"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewBasicHotStuffClient(conn)
	rand.Seed(time.Now().UnixNano())
	for {
		time.Sleep(time.Millisecond * 500)
		_, err = client.SendRequest(context.Background(), &pb.Msg{Payload: &pb.Msg_Request{Request: &pb.Request{
			Cmd:           strconv.Itoa(rand.Intn(100)) + "," + strconv.Itoa(rand.Intn(100)),
			ClientAddress: "localhost:9999",
		}}})
	}
}
