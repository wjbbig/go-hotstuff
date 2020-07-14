package main

import (
	"context"
	go_hotstuff "github.com/wjbbig/go-hotstuff"
	"github.com/wjbbig/go-hotstuff/config"
	"github.com/wjbbig/go-hotstuff/logging"
	pb "github.com/wjbbig/go-hotstuff/proto"
	"google.golang.org/grpc"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var logger = logging.GetLogger()

const timeoutDuration = time.Second * 2

type command string

type HotStuffClient struct {
	consensusResult        map[command]reply
	hotStuffConfig         config.HotStuffConfig
	requestTimeout         *go_hotstuff.Timer
	requestTimeoutDuration time.Duration
	replyChan              chan *pb.Msg
	cancelFunc             context.CancelFunc
	mut                    sync.Mutex
}

type reply struct {
	result string
	count  int
}

func NewHotStuffClient() *HotStuffClient {
	client := &HotStuffClient{
		consensusResult:        make(map[command]reply),
		hotStuffConfig:         *config.NewHotStuffConfig(),
		requestTimeoutDuration: timeoutDuration,
		requestTimeout:         go_hotstuff.NewTimer(timeoutDuration),
		replyChan:              make(chan *pb.Msg),
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	client.cancelFunc = cancelFunc
	client.requestTimeout.Init()
	client.requestTimeout.Stop()
	go client.receiveReply(ctx)
	return client
}

func (client *HotStuffClient) getResults(cmd command) (re reply, b bool) {
	client.mut.Lock()
	defer client.mut.Unlock()
	re, b = client.consensusResult[cmd]
	return
}

func (client *HotStuffClient) setResult(cmd command, re reply) {
	client.mut.Lock()
	defer client.mut.Unlock()
	client.consensusResult[cmd] = re
}

func (client *HotStuffClient) receiveReply(ctx context.Context) {
	for {
		select {
		case msg := <-client.replyChan:
			replyMsg := msg.GetReply()
			cmd := command(replyMsg.Command)
			if re, ok := client.getResults(cmd); ok {
				if re.result == replyMsg.Result {
					re.count++
					if re.count == client.hotStuffConfig.F+1 {
						logger.WithField("cmd", cmd).Info("consensus success.")
					}
					client.setResult(cmd, re)
				}
			} else {
				re := reply{
					result: re.result,
					count:  1,
				}
				client.setResult(cmd, re)
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	client := pb.NewBasicHotStuffClient(conn)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 10000; i++ {
		time.Sleep(time.Millisecond * 200)
		logger.Info("Send request~~~~")
		_, err = client.SendRequest(context.Background(), &pb.Msg{Payload: &pb.Msg_Request{Request: &pb.Request{
			Cmd:           strconv.Itoa(rand.Intn(100)) + "," + strconv.Itoa(rand.Intn(100)),
			ClientAddress: "localhost:9999",
		}}})
	}
}
