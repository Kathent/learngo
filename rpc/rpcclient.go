package rpc

import (
	"net/rpc"
	"fmt"
	"learngo/utils"
)

const (
	addr = "127.0.0.1"
    port = "1234"
)

//
func SendRpcClient() error{
	var address string = addr + ":" + port

	client, err := rpc.DialHTTP("tcp", address)

	utils.CheckError(err)

	args := "Hello Server"

	var result RpcObject

	err = client.Call("RpcObject.SayHello", args, &result)

	utils.CheckError(err)

	fmt.Println(result.Result)

	return err
}
