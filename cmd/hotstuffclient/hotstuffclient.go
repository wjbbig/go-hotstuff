package main

import (
	"context"
	pb "github.com/wjbbig/go-hotstuff/proto"
)

type hotStuffGRPCClient struct {
	client *HotStuffClient
}

func (gc *hotStuffGRPCClient) SendMsg(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	gc.client.replyChan <- in
	return &pb.Empty{}, nil
}

func (gc *hotStuffGRPCClient) SendRequest(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (gc *hotStuffGRPCClient) SendReply(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}
