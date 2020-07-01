package basic

import (
	pb "github.com/wjbbig/go-hotstuff/proto"
	"testing"
	"time"
)

func TestBasicHotStuff_ReceiveMsg(t *testing.T) {
	stuff := NewBasicHotStuff(1)
	request := &pb.Msg_Request{Request:&pb.Request{
		Cmd:           "1+2",
		ClientAddress: "localhost:9999",
	}}
	msg := &pb.Msg{Payload:request}
	stuff.MsgEntrance <- msg
	stuff.MsgEntrance <- msg
	stuff.MsgEntrance <- msg
	defer stuff.Close()
	time.Sleep(time.Second*2)
}
