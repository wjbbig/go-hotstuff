package basic

import (
	"context"
	pb "go-hotstuff/proto"
)

type BasicHotStuffService struct {}

func (basic *BasicHotStuffService) Broadcast(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (basic *BasicHotStuffService) SendVote(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

