package learn_etcd

import (
	"time"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
)

var cl *clientv3.Client

func GetValue(key string) ([]interface{}, error){
	resp, err := cl.Get(context.Background(), key)
	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("resp %v", resp))

	res := []interface{}{}
	for index, val := range resp.Kvs {
		fmt.Println(fmt.Sprintf("key %d, val %v", index, val))
		if string(val.Key) == key {
			res = append(res, val.Value)
		}
	}

	return res, nil
}

func Put(key string, val string) (interface{}, error){
	put, err := cl.Put(context.Background(), key, val)
	if err != nil {
		return nil, err
	}

	if put != nil && put.PrevKv != nil{
		return put.PrevKv.Value, nil
	}

	return nil, nil
}

func Watch(key string) {
	watcher := cl.Watch(context.Background(), key)

	for resp := range watcher{
		fmt.Println(fmt.Sprintf("resp is.... %v", resp))
	}
}

func LoadClient(addr string) error{
	cli, error := clientv3.New(clientv3.Config{
		Endpoints: []string{addr}, DialTimeout: time.Second * 10})
	if error != nil {
		return error
	}

	cl = cli

	return nil
}

func LearnEtcd(){
	err := LoadClient("192.168.96.140:2379")
	if err != nil{
		panic(err)
	}

	key := "etcd_learn_go"
	go func() {
		Watch(key)
	}()

	put, err := Put(key, "start")
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("pre put val is ..%v", put))

	val, valErr := GetValue(key)
	if valErr != nil{
		panic(valErr)
	}

	fmt.Println(fmt.Sprintf("get value is...%v",val))
}
