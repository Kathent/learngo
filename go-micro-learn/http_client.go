package go_micro_learn

import (
	"learngo/proto1"
	"google.golang.org/grpc"
	"context"
)

func SendMessage() {
	conn, _ := grpc.Dial("localhost:8888")
	client := proto.NewGreeterSClient(conn)
	for {
		client.Hello(context.Background(), &proto.HelloRequestS{Name:"aaa"})
	}
}
