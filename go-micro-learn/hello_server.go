package go_micro_learn

import (
	"learngo/proto1"
	"google.golang.org/grpc"
	"net"
	"golang.org/x/net/context"
	"github.com/alecthomas/log4go"
)

type GreeterServerImpl struct {
	proto.GreeterSServer
}

var helloServer GreeterServerImpl

func (GreeterServerImpl) Hello(ctx context.Context, req *proto.HelloRequestS) (*proto.HelloResponseS, error){
	var resp proto.HelloResponseS
	resp.Greeting = req.Name
	log4go.Info("name:%s", req.Name)
	return &resp, nil
}

func StartHelloServer() {
	s := grpc.NewServer()
	proto.RegisterGreeterSServer(s, helloServer)

	listener, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		log4go.Info("StartHelloServer err:%v", err)
	}
	s.Serve(listener)
}
