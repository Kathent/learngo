package main

import (
	//"rpc"
	//"learngo/learn_etcd"
	"learngo/go-micro-learn"
)


func main() {
	//q := queue.New()
	//fmt.Println(q)

	//go rpc.AcceptRpc()
	//go rpc.SendRpcClient()
	//learn_etcd.LearnEtcd()
	go_micro_learn.StartServer()
}
