package basic

import (
	pb "github.com/wjbbig/go-hotstuff/proto"
	"os"
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
	defer stuff.SafeExit()
	time.Sleep(time.Second*2)
}

func TestRemoveDBFile(t *testing.T) {
	err := os.RemoveAll("/opt/hotstuff/dbfile")
	if err != nil {
		t.Fatal(err)
	}
}
