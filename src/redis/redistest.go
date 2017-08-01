package main

import (
	"github.com/go-redis/redis"
	"fmt"
	"utils"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		/*Password:"12345",*/
		DB:0,
	})

	defer client.Close()


	//ping
	pong, err := client.Ping().Result()
	utils.CheckError(err)
	fmt.Println(pong, err)

	err = client.Set("name", "xl", 0).Err()
	utils.CheckError(err)

	_, err = client.Get("xl").Result()
	if err == redis.Nil{
		fmt.Println("not exist")
	}else {
		utils.CheckError(err)
	}

	subscribe := client.Subscribe("aaa")
	channel := subscribe.Channel()
	dealMessage(channel)


	//subKeyChange := client.Subscribe("__keyspace@0__:lkj")
	subKeyChange2 := client.Subscribe("__keyevent@0__:hset")
	//dealMessage(subKeyChange.Channel())
	dealMessage(subKeyChange2.Channel())
	client.HMSet("lkj", map[string]interface{}{"aa":1, "bb":2})

	client.Publish("aaa", "bbbbb")
	for true{

	}
}

func dealMessage(messages <-chan *redis.Message) {
	go func() {
		for {
			message := <- messages
			fmt.Println("Recieved a message ", message)
		}
	}()
}
