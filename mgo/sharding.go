package mgo

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"fmt"
)

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
