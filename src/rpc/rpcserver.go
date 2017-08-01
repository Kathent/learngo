package rpc

import (
	"net/rpc"
	"net/http"
	"utils"
	//"reflect"
	//"log"
	//"fmt"
	//"fmt"
)

func AcceptRpc() {
	q := new(RpcObject)

	rpc.Register(q)

	rpc.HandleHTTP()

	err := http.ListenAndServe(":1234", nil)

	utils.CheckError(err)
}

//func Register(object *RpcObject) {
//	typeOf := reflect.ValueOf(object)
//	for i := 0; i < typeOf.NumMethod(); i++ {
//		method := typeOf.Method(i)
//		fmt.Println(method, method.Type())
//		makeFunc := reflect.MakeFunc(method.Type(), func(args []reflect.Value) (results []reflect.Value) {
//			log.Println(args, results)
//			return results
//		})
//
//		method.Set(makeFunc)
//	}
//	rpc.Register(object)
//}
