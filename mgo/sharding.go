package mgo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"fmt"
	"time"
	"math/rand"
)

func InsertNew(size int) {
	session, err := mgo.Dial("192.168.96.224:27017")
	if err != nil {
		panic(err)
	}
	db := session.DB("sharding_test")
	c := db.C("sharding_collection")
	index := 0
	bulk := c.Bulk()

	vccId := []int{782, 10010, 332, 445}
	for i := 0; i < size; i++ {
		unix := time.Now().Unix()
		bulk.Insert(bson.M{"vcc_id": vccId[rand.Intn(len(vccId))], "start_time": unix, "rand_val": rand.Int63n(unix)})
		index++
		if index >= 10000 {
			index = 0
			result, err := bulk.Run()
			if err != nil {
				panic(err)
			}

			fmt.Println(result)
			bulk = c.Bulk()
		}else if i >= size {
			index = 0
			result, err := bulk.Run()
			if err != nil {
				panic(err)
			}

			fmt.Println(result)
		}
	}
}

func  QueryMongo(collection *mgo.Collection, size int) {
	st := struct {
		VccId int `bson:"vcc_id"`
		StartTime int64 `bson:"start_time"`
		RandVal int64 `bson:"rand_val"`
	}{}

	for i := 0; i < size; i++ {
		randN := rand.Int63n(time.Now().Unix())
		find := collection.Find(bson.M{"vcc_id": rand.Intn(4), "start_time": randN,
			"rand_val": randN}).One(&st)
		if find != mgo.ErrNotFound {
			panic(find)
		}
	}
}

func Insert() {
	session, err := mgo.Dial("192.168.96.204:30031")
	if err != nil {
		panic(err)
	}
	db := session.DB("sharding_test")
	c := db.C("sharding_test_c")
	index := 0
	bulk := c.Bulk()
	for i := 100000; i < 200000; i++ {
		bulk.Insert(bson.M{"user_id": i, "name": strconv.Itoa(i)})
		index++
		if index >= 100 {
			index = 0
			result, err := bulk.Run()
			if err != nil {
				panic(err)
			}

			fmt.Println(result)
			bulk = c.Bulk()
		}else if i == 100000 {
			index = 0
			result, err := bulk.Run()
			if err != nil {
				panic(err)
			}

			fmt.Println(result)
		}
	}
}

func DeleteAll()  {
	session, err := mgo.Dial("192.168.96.204:30031")
	if err != nil {
		panic(err)
	}
	db := session.DB("sharding_test")
	c := db.C("sharding_test_c")
	info, err := c.RemoveAll(bson.M{})
	if err != nil {
		panic(err)
	}

	fmt.Println(info)
}
