package main

import (
	"reflect"
)

type TestStruct struct {
	c int
	d string
}

type preSpring struct {
	a int
	b string
	e TestStruct `autowire:`
}

type BeanDefinition struct {
	BeanType reflect.Type
}

func main() {
	register(preSpring{})
}

func register(any interface{}) {
	//typeOf := reflect.TypeOf(any)
	//BeanDefinition{BeanType:typeOf}
}
