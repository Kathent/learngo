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
	"github.com/micro/go-plugins/server/http"
	http2 "net/http"
	"io/ioutil"
	"github.com/gin-gonic/gin"
)
const(
	ETCD_ADDR = "192.168.1.4:2379"
	GRPC_ADDR = "localhost:7777"
)

type Greeter struct{
	Charactor string `json:"charactor"`
}

func (g *Greeter)Hello(ctx context.Context, r *proto.HelloRequest, s *proto.HelloResponse) error{
	s.Greeting = fmt.Sprintf("Hello Mr.%s", r.Name)
	return nil
}

func StartServer(){
	registry := etcdv3.NewRegistry(rg.Addrs(ETCD_ADDR))
	//s := grpc.NewServer(server.Address(GRPC_ADDR), server.Name("hello"),
	//	server.Version("latest"),)
	mux := http2.NewServeMux()
	mux.HandleFunc("/hello", func(writer http2.ResponseWriter, request *http2.Request) {
		bytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			panic(err)
		}
		log.Println(fmt.Sprintf("param:%s", string(bytes)))
		writer.Write([]byte(`{"2":"3"}`))
	})
	s := http.NewServer(server.Address(GRPC_ADDR), server.Name("hello"),
		server.Version("latest"))
	s.Handle(s.NewHandler(mux))

	service := micro.NewService(
		micro.Server(s),
		micro.Registry(registry),
		micro.RegisterTTL(time.Second * 10))

	service.Init()

	//var h proto.GreeterHandler = new(Greeter)

	//proto.RegisterGreeterHandler(service.Server(), h)

	log.Fatal(service.Run())
}

func StartGrpcServer(){
	registry := etcdv3.NewRegistry(rg.Addrs(ETCD_ADDR))
	s := grpc.NewServer(server.Address(GRPC_ADDR), server.Name("hello"),
		server.Version("latest"),)

	service := micro.NewService(
		micro.Server(s),
		micro.Registry(registry),
		micro.RegisterTTL(time.Second * 20))

	//service.Init()

	var h proto.GreeterHandler = new(Greeter)

	proto.RegisterGreeterHandler(service.Server(), h)

	log.Fatal(service.Run())
}

func StartGinServer(){
	registry := etcdv3.NewRegistry(rg.Addrs(ETCD_ADDR))
	//s := grpc.NewServer(server.Address(GRPC_ADDR), server.Name("hello"),
	//	server.Version("latest"),)
	//mux := http2.NewServeMux()
	//mux.HandleFunc("/hello", func(writer http2.ResponseWriter, request *http2.Request) {
	//	bytes, err := ioutil.ReadAll(request.Body)
	//	if err != nil {
	//		panic(err)
	//	}
	//	log.Println(fmt.Sprintf("param:%s", string(bytes)))
	//	writer.Write([]byte(`{"2":"3"}`))
	//})

	g := gin.Default()
	g.POST("/hello", func(c *gin.Context) {
		c.JSON(http2.StatusOK, Greeter{Charactor: "aaaa"})
	})

	s := http.NewServer(server.Address(GRPC_ADDR), server.Name("hello"),
		server.Version("latest"))
	s.Handle(s.NewHandler(g))

	service := micro.NewService(
		micro.Server(s),
		micro.Registry(registry),
		micro.RegisterTTL(time.Second * 10))

	service.Init()

	//var h proto.GreeterHandler = new(Greeter)

	//proto.RegisterGreeterHandler(service.Server(), h)

	log.Fatal(service.Run())
}