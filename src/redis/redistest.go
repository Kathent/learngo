package main

import (
	"github.com/go-redis/redis"
	"fmt"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password:"12345",
		DB:0,
	})


	//ping
	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

	err = client.Set("name", "xl", 0).Err()
	if  err != nil{
		panic(err)
	}

	i, err := client.Get("xl").Result()
	if err == redis.Nil{
		fmt.Println("not exist")
	} else if err != nil {
		panic(err)
	}

	fmt.Println("result is :", i)
}
