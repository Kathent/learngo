package rpc


type RpcObject struct{
	Result interface{}
}



func (*RpcObject) SayHello(params string, result *RpcObject) error{
	result.Result = params
	return nil
}
