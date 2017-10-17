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
	"github.com/micro/go-plugins/client/http"
	client2 "github.com/micro/go-micro/client"
	"time"
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

	//client := grpc.NewClient()
	client := http.NewClient(client2.ContentType("application/json"))
	service := micro.NewService(micro.Client(client), micro.Name("hello"), micro.Registry(registry))

	//greeter := proto.NewGreeterClient("hello", service.Client())

	//response, err := greeter.Hello(context.Background(), &proto.HelloRequest{Name: fmt.Sprintf("service:%d", &service)})
	//if err != nil {
	//	panic(err)
	//}
	var res string
	req := client.NewRequest("hello", "/hello", `{"1":"2"}`, )
	err = service.Client().Call(context.TODO(), req, &res)
	if err == nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("res is %+v len:%d", res, len(res)))
}

func StartGrcpClient() {
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

	//service := micro.NewService(micro.Client(grpc.NewClient()), micro.Name("hello"), micro.Registry(registry))

	cl := grpc.NewClient(client2.Registry(registry))
	greeter := proto.NewGreeterClient("hello", cl)

	response, err := greeter.Hello(context.Background(), &proto.HelloRequest{Name: fmt.Sprintf("service:%d", &cl)})
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("res is %+v", response))
}

type TimerObj struct {
	ModuleName string `json:"module_name"`//模块名称
	Expire int64 `json:"expire"`//过期时间
	CallBack string `json:"call_back"`//回调方法
}

type JsonObj struct {
	Code string `json:"code"`
	Msg string `json:"msg"`
}

func StartGinClient(){
	registry := etcdv3.NewRegistry(registry2.Addrs(ETCD_ADDR))

	//client := grpc.NewClient()
	client := http.NewClient(client2.ContentType("application/json"))
	service := micro.NewService(micro.Client(client), micro.Name("hello"), micro.Registry(registry))

	obj := TimerObj{
		ModuleName: "test",
		Expire: time.Now().Unix() + 100,
		CallBack: "/hello",
	}
	req := client.NewRequest("TIMER", "/hello", obj)

	res := JsonObj{}
	err := service.Client().Call(context.TODO(), req, &res)
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Sprintf("res is %+v", res))
}
