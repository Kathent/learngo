package go_micro_learn

import (
	"github.com/micro/go-micro"
	"learngo/proto"
	"context"
	"fmt"
	"log"
	"github.com/micro/go-plugins/registry/etcdv3"
	registry2 "github.com/micro/go-micro/registry"
	"github.com/micro/go-plugins/client/grpc"
)

func StartClient(){
	registry := etcdv3.NewRegistry(registry2.Addrs(ETCD_ADDR))
	services, err := registry.GetService("hello")

	if err != nil {
		panic(err)
	}

	for _, s := range services {
		log.Println(fmt.Sprintf("service:%+v len:%d", s, len(s.Nodes)))
		for _, node := range s.Nodes{
			fmt.Println(fmt.Sprintf("node :%v", node))
		}

		for _, ep := range s.Endpoints {
			fmt.Println(fmt.Sprintf("ep :%v", ep))
		}

		for key, md := range s.Metadata{
			fmt.Println(fmt.Sprintf("k:%s,v:%s", key, md))
		}
	}

	service := micro.NewService(micro.Client(grpc.NewClient()), micro.Name("hello"), micro.Registry(registry))

	greeter := proto.NewGreeterClient("hello", service.Client())

	response, err := greeter.Hello(context.Background(), &proto.HelloRequest{Name: fmt.Sprintf("service:%d", &service)})
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("res is %+v", response))
}
