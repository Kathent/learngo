package main

import (
	"learngo/go-micro-learn"
	//"learngo/learn_etcd"
	"time"
)


func main() {
	//q := queue.New()
	//fmt.Println(q)

	//go rpc.AcceptRpc()
	//go rpc.SendRpcClient()
	//learn_etcd.LearnEtcd()
	go go_micro_learn.StartServer()
	time.Sleep(time.Second * 4)
	go_micro_learn.StartClient()
	//file_read_analysis.LoadFiles("load_file_dir", "2017-10-02 10:00:00", "2017-10-03 14:00:00")
}
