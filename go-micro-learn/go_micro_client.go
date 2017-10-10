package go_micro_learn

import (
	"github.com/micro/go-micro"
	"learngo/proto"
	"context"
	"fmt"
	"log"
)

func StartClient(){
	service := micro.NewService(micro.Name("new client"))

	greeter := proto.NewGreeterClient("greeter", service.Client())

	response, err := greeter.Hello(context.Background(), &proto.HelloRequest{Name: fmt.Sprintf("service:%d", &service)})
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("res is %+v", response))
}
