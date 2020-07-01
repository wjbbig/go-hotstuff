package basic

import (
	"context"
	pb "go-hotstuff/proto"
)



type BasicHotStuffService struct {
	BasicHotStuff *BasicHotStuff
}

func (basic *BasicHotStuffService) Broadcast(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	basic.BasicHotStuff.MsgEntrance <- in
	return &pb.Empty{}, nil
}

func (basic *BasicHotStuffService) SendVote(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	basic.BasicHotStuff.MsgEntrance <- in
	return &pb.Empty{}, nil
}

func (basic *BasicHotStuffService) SendRequest(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	basic.BasicHotStuff.MsgEntrance <- in
	return &pb.Empty{}, nil
}


func (basic *BasicHotStuffService) SendReply(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func NewBasicHotStuffService(id int, networkType string) *BasicHotStuffService {
	// TODO FACTORY
	return &BasicHotStuffService{BasicHotStuff:NewBasicHotStuff(id)}
}

