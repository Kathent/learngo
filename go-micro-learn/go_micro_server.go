package go_micro_learn

import "github.com/micro/go-micro"

func StartServer(){
	micro.NewService(micro.Name("hello"),
		micro.Version("latest"))
}
