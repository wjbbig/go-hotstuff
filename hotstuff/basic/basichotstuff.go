package basic

import (
	"go-hotstuff/hotstuff"
	pb "go-hotstuff/proto"
)

type BasicHotStuff struct {
	hotstuff.HotStuffImpl
	MsgEntrance chan *pb.Msg   // receive msg
	prepareQC   *pb.QuorumCert // highQC
	preCommitQC *pb.QuorumCert // lockQC
	commitQC    *pb.QuorumCert
}

func NewBasicHotStuff() *BasicHotStuff {
	msgEntrance := make(chan *pb.Msg)
	bhs := &BasicHotStuff{
		MsgEntrance: msgEntrance,
		prepareQC:   nil,
		preCommitQC: nil,
		commitQC:    nil,
	}
	bhs.View = hotstuff.NewView(1, 1)
	go bhs.receiveMsg()

	return bhs
}

// receiveMsg receive msg from msg channel
func (bhs *BasicHotStuff) receiveMsg() {
	for {
		select {
		case msg := <-bhs.MsgEntrance:
			bhs.handleMsg(msg)
		case <-bhs.TimeChan.C:
			// TODO: timout, goto next view with highQC
		}
	}
}

// handleMsg handle different msg with different way
func (bhs *BasicHotStuff) handleMsg(msg *pb.Msg) {
	switch msg {
	}
}
