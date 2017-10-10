package go_micro_learn

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/etcd"
	registry2 "github.com/micro/go-micro/registry"
	proto "learngo/proto"
	"golang.org/x/net/context"
	"fmt"
	"log"
)
const(
	ETCD_ADDR = "192.168.96.140:2379"
)

type Greeter struct{}

func (g *Greeter)Hello(ctx context.Context, r *proto.HelloRequest, s *proto.HelloResponse) error{
	s.Greeting = fmt.Sprintf("Hello Mr.%s", r.Name)
	return nil
}

func StartServer(){
	registry := etcd.NewRegistry(registry2.Addrs(ETCD_ADDR))

	service := micro.NewService(micro.Name("hello"),
		micro.Version("latest"),
		micro.Registry(registry))

	service.Init()

	var h proto.GreeterHandler = new(Greeter)

	proto.RegisterGreeterHandler(service.Server(), h)

	log.Fatal(service.Run())
}
