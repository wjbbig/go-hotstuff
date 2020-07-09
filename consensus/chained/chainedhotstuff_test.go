package chained

import (
	pb "github.com/wjbbig/go-hotstuff/proto"
	"testing"
	"time"
)

func TestChainedHotStuff(t *testing.T) {
	chained := NewChainedHotStuff(1, nil)
	request := &pb.Msg_Request{Request:&pb.Request{
		Cmd:           "1+2",
		ClientAddress: "localhost:9999",
	}}
	msg := &pb.Msg{Payload:request}
	chained.MsgEntrance <- msg

	time.Sleep(10*time.Second)
}