package go_micro_learn

import (
	"github.com/micro/go-micro"
	"github.com/micro/go-plugins/registry/etcd"
	registry2 "github.com/micro/go-micro/registry"
)
const(
	ETCD_ADDR = "192.168.96.140:2379"
)

func StartServer(){
	registry := etcd.NewRegistry(registry2.Addrs(ETCD_ADDR))

	service := micro.NewService(micro.Name("hello"),
		micro.Version("latest"),
		micro.Registry(registry))

	service.Init()
}
