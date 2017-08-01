package main

import (
	"reflect"
	"encoding/json"
	"utils"
	"fmt"
	"queuenet"
	"github.com/jinzhu/configor"
)


type FuncRequest struct {
	FuncName string
	Params []Any
}

type Any interface {}

type Func interface {}

type PrintInst struct {

}

func (PrintInst)PrintSomething(a string, b string) string{
	return a + b
}

type MethodInvocation struct {
	method reflect.Value
}

type Printer interface {
	PrintSomething(a string, b string) string
}

func NewClient() func() Printer{
	var instance Printer
	return func() Printer {
		if instance == nil{
			return new(PrintInst)
		}
		return instance
	}
}

var Config = struct {
	APPName string `default:"app name"`

	DB struct {
		Name     string
		User     string `default:"root"`
		Password string `required:"true" env:"DBPassword"`
		Port     uint   `default:"3306"`
	}

	Contacts []struct {
		Name  string
		Email string `required:"true"`
	}
}{}

func main() {
	//testReflect()
	//testTcp()
	//test11()

	configor.Load(&Config, "config.yml")
	fmt.Printf("config: %#v", Config)
}

func test11() {
	client := NewClient()()
	fmt.Printf("%v \n", client)

	client2 := NewClient()()
	fmt.Printf("%v \n", client2)

	fmt.Println(client == client2)

	fmt.Printf("%T, %#v %p %p\n", client, client, client, client2)
	var fn  func(a int, b int) int = func(a int, b int) int {
		return 1
	}

	makeFunc := reflect.MakeFunc(reflect.TypeOf(&fn).Elem(), func(args []reflect.Value) ([]reflect.Value) {
		return []reflect.Value{reflect.ValueOf(2)}
	})

	fmt.Println(makeFunc.Call([]reflect.Value{reflect.ValueOf(0), reflect.ValueOf(0)})[0])
}

func testTcp() {
	server := queuenet.NewServer("127.0.0.1:2222")
	client := queuenet.NewClient("127.0.0.1:2222")
	go server.StartServer()
	go client.StartClient()

	for true {

	}
}
func testReflect() {
	register(&PrintInst{})
	var params = []Any{"a", "b"}
	var fc = FuncRequest{
		FuncName: "PrintSomething",
		Params:   params,
	}
	var fcCreated FuncRequest
	marshal, error := json.Marshal(fc)
	utils.CheckError(error)
	json.Unmarshal(marshal, &fcCreated)
	fmt.Println(fcCreated)
	invocation := invocations[fcCreated.FuncName]
	var values = make([]reflect.Value, len(fcCreated.Params))
	for ind, tmp := range fcCreated.Params {
		values[ind] = reflect.ValueOf(tmp)
		fmt.Println(values[ind])
	}
	results := invocation.method.Call(values)
	fmt.Println(results)
}

var invocations = make(map[string]MethodInvocation)

func register(inst interface{}) {
	instValue := reflect.ValueOf(inst).Elem()
	instType := reflect.TypeOf(inst).Elem()
	for i := 0; i < instValue.NumMethod(); i++ {
		var method reflect.Value = instValue.Method(i)
		var typeMethod reflect.Method = instType.Method(i)
		var invocation MethodInvocation = MethodInvocation{
			method: method,
		}

		invocations[typeMethod.Name] = invocation
	}
}
