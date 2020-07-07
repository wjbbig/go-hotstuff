package consensus

import (
	"context"
	pb "github.com/wjbbig/go-hotstuff/proto"
)



type HotStuffService struct {
	hotStuff HotStuff
}

func (basic *HotStuffService) SendMsg(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	basic.hotStuff.GetMsgEntrance() <- in
	return &pb.Empty{}, nil
}

func (basic *HotStuffService) SendRequest(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	basic.hotStuff.GetMsgEntrance() <- in
	return &pb.Empty{}, nil
}


func (basic *HotStuffService) SendReply(ctx context.Context, in *pb.Msg) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (basic *HotStuffService) GetImpl() HotStuff {
	return basic.hotStuff
}

func NewHotStuffService(impl HotStuff) *HotStuffService {
	return &HotStuffService{hotStuff: impl}
}

