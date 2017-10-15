package main

import (
	//"rpc"
	//"learngo/learn_etcd"
	"learngo/file_read_analysis"
)


func main() {
	//q := queue.New()
	//fmt.Println(q)

	//go rpc.AcceptRpc()
	//go rpc.SendRpcClient()
	//learn_etcd.LearnEtcd()
	//go_micro_learn.StartServer()
	file_read_analysis.LoadFiles("load_file_dir", "2017-10-02 10:00:00", "2017-10-03 14:00:00")
}
