package mgo

import (
	"testing"
	"time"
	"gopkg.in/mgo.v2"
)

func TestInsert(t *testing.T) {
	Insert()
}

func TestDeleteAll(t *testing.T) {
	DeleteAll()
}

func TestInsertNew(t *testing.T) {
	start := time.Now().Unix()
	InsertNew(10000000)
	end := time.Now().Unix()

	t.Logf("consume time:%d", end - start)
}

func BenchmarkQueryMongo(b *testing.B) {
	session, err := mgo.Dial("192.168.96.224:27017")
	if err != nil {
		panic(err)
	}
	db := session.DB("sharding_test")
	c := db.C("sharding_collection")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			QueryMongo(c, 1)
		}
	})
}
