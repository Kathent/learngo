package mgo

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func LearnMgo() {
	session, err := mgo.Dial("192.168.96.140")
	if err != nil {
		panic(err)
	}
	session.SetSafe(&mgo.Safe{})

	collection := session.DB("local").C("student")
	fmt.Println(fmt.Sprintf("%+v", collection))

	type tmpS struct {
		Code string
		Msg string
	}

	type tmpSS struct {
		TotalPrice int32 `bson:"totalPrice"`
		AverageQuantity int32 `bson:"averageQuantity"`
		Count int32 `bson:"count"`
	}

	type tmpSSS struct {
		Code string
		Msg string
		Count int32
	}

	var tmpssss = &tmpS{
		Code: "123",
		Msg:  "234",
	}
	//insertErr := collection.Insert(tmpssss)
	//
	//if insertErr != nil {
	//	panic(insertErr)
	//}

	var newTmps tmpS
	err = collection.Find(bson.M{"code": "123"}).One(&newTmps)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("%+v, %+v", newTmps, tmpssss))

	var result []tmpS
	err = collection.Find(bson.M{"code": bson.M{"$ne": "234"}}).All(&result)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("result:%+v", result))

	err = collection.Find(bson.M{"$and": []bson.M{{"code": "123"}, {"msg":"234"}}}).All(&result)
	fmt.Println(fmt.Sprintf("result:%+v", result))
	bulk := collection.Bulk()
	//bulk.Update(bson.M{"code": "123"}, bson.M{"$push": bson.M{"push": "111"}})
	//bulk.Update(bson.M{"code": "123"}, bson.M{"$set": bson.M{"push": "222"}})
	//bulk.Update(bson.M{"code": "123"}, bson.M{"$pull": bson.M{"push": "111"}})
	//bulk.Remove(bson.M{"code":  "123", "push":"222"})
	bulk.UpdateAll(bson.M{}, bson.M{"$unset": bson.M{"id_": 1}})
	bulkResult, err := bulk.Run()
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("bulk result..%+v", bulkResult))


	pipe := collection.Pipe([]bson.M{{"$match": bson.M{"code": bson.M{"$gt":"111"}}},
	  								 {"$group":	bson.M{"_id": "$code", "count": bson.M{"$sum": 1}, "max": bson.M{"$max":"$code"}}},
	  								 {"$sort": bson.M{"msg": -1}},
	  								 {"$limit": 3},
	  								 })
	var arr []tmpSSS
	allErr := pipe.All(&arr)
	if allErr != nil {
		panic(allErr)
	}

	fmt.Println(fmt.Sprintf("pipe result:%+v", arr))

	i := session.DB("local").C("123").Pipe([]bson.M{
		{"$group": bson.M{
			"_id":             bson.M{"month": bson.M{"$month": "$date"}, "day": bson.M{"$dayOfMonth": "$date"}, "year": bson.M{"$year": "$date"}},
			"totalPrice":      bson.M{"$sum": bson.M{"$multiply": []string{"$price", "$quantity"}}},
			"averageQuantity": bson.M{"$avg": "$quantity"},
			"count":           bson.M{"$sum": 1},
		}},
	})

	var sssss = make([]tmpSS, 0)
	all := i.All(&sssss)
	if all != nil {
		panic(all)
	}

	fmt.Println(fmt.Sprintf("%+v", sssss))

	//js := `function() {emit(this.code, this.msg)},
    //function(key, values) {return values.join("")},
    //{
    //   query: {code: "123"},
    //   out: "result",
    //}`

    var res []map[string]interface{}
	reduce := &mgo.MapReduce{
		Map:    "function() {emit(this.code, this.msg)}",
		Reduce: "function(key, values) {return values.join(\"123\")}",
	}

	info, err := collection.Find(nil).MapReduce(reduce, &res)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("res:%+v", info))
}
