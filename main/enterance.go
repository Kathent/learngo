package main

import (
	//"rpc"
	"learngo/learn_etcd"
)


func main() {
	//q := queue.New()
	//fmt.Println(q)

	//go rpc.AcceptRpc()
	//go rpc.SendRpcClient()
	go learn_etcd.LearnEtcd()

	for true{

	}
}
