package go_micro_learn

import (
	"github.com/micro/go-micro"
	rg "github.com/micro/go-micro/registry"
	proto "learngo/proto"
	"golang.org/x/net/context"
	"fmt"
	"log"
	"github.com/micro/go-plugins/registry/etcdv3"
	"time"
	"github.com/micro/go-plugins/server/grpc"
	"github.com/micro/go-micro/server"
)
const(
	ETCD_ADDR = "192.168.1.4:2379"
	GRPC_ADDR = ":7777"
)

type Greeter struct{}

func (g *Greeter)Hello(ctx context.Context, r *proto.HelloRequest, s *proto.HelloResponse) error{
	s.Greeting = fmt.Sprintf("Hello Mr.%s", r.Name)
	return nil
}

func StartServer(){
	registry := etcdv3.NewRegistry(rg.Addrs(ETCD_ADDR))
	s := grpc.NewServer(server.Address(GRPC_ADDR), server.Name("hello"),
		server.Version("latest"),)

	service := micro.NewService(
		micro.Server(s),
		micro.Registry(registry),
		micro.RegisterTTL(time.Second * 2000))

	//service.Init()

	var h proto.GreeterHandler = new(Greeter)

	proto.RegisterGreeterHandler(service.Server(), h)

	log.Fatal(service.Run())
}
