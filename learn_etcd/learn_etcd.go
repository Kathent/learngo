package learn_etcd

import (
	"time"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	//"log"
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
	err := LoadClient("192.168.1.4:2379")
	if err != nil{
		panic(err)
	}

	key := "etcd_learn_go"
	go func() {
		Watch(key)
	}()

	delete, err := cl.Delete(context.Background(), "logic", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}

	fmt.Println(delete)
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
	//
	//grant, grantErr := cl.Grant(context.Background(), 4)
	//if grantErr != nil {
	//	panic(grantErr)
	//}
	//
	//resp, respErr := cl.Put(context.Background(), "ttlKey", "ttlVal", clientv3.WithLease(grant.ID))
	//if respErr != nil {
	//	panic(respErr)
	//}
	//
	//log.Println(fmt.Sprintf("put resp %v", resp))
	//
	//get, getErr := cl.Get(context.Background(), "ttlKey")
	//if getErr != nil {
	//	panic(getErr)
	//}
	//
	//log.Println(fmt.Sprintf("get resp %v", get))
	//
	//_, onceErr := cl.KeepAliveOnce(context.Background(), grant.ID)
	//if onceErr != nil {
	//	panic(onceErr)
	//}
	//
	//
	//get, getErr = cl.Get(context.Background(), "ttlKey")
	//if getErr != nil {
	//	panic(getErr)
	//}
	//log.Println(fmt.Sprintf("get resp %v", get))
	//
	//
	//get, getErr = cl.Get(context.Background(), "ttlKey")
	//if getErr != nil {
	//	panic(getErr)
	//}
	//log.Println(fmt.Sprintf("get resp %v", get))
	//
	//alive, aliveErr := cl.KeepAlive(context.Background(), grant.ID)
	//if aliveErr != nil {
	//	panic(aliveErr)
	//}
	//
	//for true {
	//	select {
	//		case r :=<- alive:
	//		log.Println("res", r)
	//	}
	//}
}
