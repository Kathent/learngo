package rpc

import (
	"net/rpc"
	"net/http"
	"utils"
)

func AcceptRpc() {
	q := new(RpcObject)

	rpc.Register(q)

	rpc.HandleHTTP()

	err := http.ListenAndServe(":1234", nil)

	utils.CheckError(err)
}
