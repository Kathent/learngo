package utils

import (
	"reflect"
	"fmt")

func CheckError(i error) {
	if i != nil {
		panic(fmt.Sprintf("fatal error.. %v", i))
	}
}

type Any interface {}
type Func func(any Any)

func FuncNoError(p interface{}, any Any){
	defer func() {
		if error := recover(); error != nil{
			fmt.Println("error is :", error)
		}
	}()

	pp := reflect.ValueOf(any)
	value := []reflect.Value{pp}
	reflect.ValueOf(p).Call(value)
}