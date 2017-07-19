package main

import "rpc"


func main() {
	//q := queue.New()
	//fmt.Println(q)

	go rpc.AcceptRpc()
	go rpc.SendRpcClient()

	for true{

	}
}
