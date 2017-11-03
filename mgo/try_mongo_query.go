package mgo

import (
	"net"
	"time"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

type(
	header struct {
		length int32
		requestId int32
		responseTo int32
		opCode int32
	}

	opQuery struct {
		header header
		flags int32
		cName string
		numberSkip int32
		numberReturn int32
		query string
		returnFieldSelector string
	}
)

func TryMongoDial() {
	tcpAddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:27017")
	if err != nil {
		panic(err)
	}

	localAddr, err := net.ResolveTCPAddr("tcp", "localhost:27016")
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp", localAddr, tcpAddr)
	if err != nil {
		panic(err)
	}

	//n, err := conn.Write([]byte{1, 1, 1, 1, 1, 1})
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Println(n)
	//
	//bt := make([]byte, 8)
	//for {
	//	conn.SetReadDeadline(time.Now().Add(time.Second * 2))
	//	n, err := conn.Read(bt)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	fmt.Println(fmt.Sprintf("read %s", string(bt[:n])))
	//}
	err = query(conn)
	if err != nil {
		panic(err)
	}

	bt := make([]byte, 256)
	for {
		conn.SetReadDeadline(time.Now().Add(time.Second * 2))
		n, err := conn.Read(bt)
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("read %s", string(bt[:n])))
	}
}

func query(conn net.Conn) error {
	bts := getQuery()
	_, err := conn.Write(bts)
	return err
}

func getQuery() []byte{
	queryBson := bson.M{
		"code": "123",
	}
	out, err := bson.Marshal(queryBson)
	if err != nil {
		panic(err)
	}

	queryString := string(out)
	out, err = bson.Marshal(bson.M{})

	var query = opQuery{
		header:header{
			requestId:1,
			responseTo:0,
			opCode:2004,
		},
		cName:"local.student",
		numberReturn:1,
		query: queryString,
		returnFieldSelector: string(out),
	}

	bts := writeQuery(query)
	return bts
}

func writeQuery(query opQuery) []byte {
	res := make([]byte, 0)
	res = addInt32(res, query.header.length)
	res = addInt32(res, query.header.requestId)
	res = addInt32(res, query.header.responseTo)
	res = addInt32(res, query.header.opCode)

	res = addInt32(res, query.flags)
	res = addString(res, query.cName)
	res = addInt32(res, query.numberSkip)
	res = addInt32(res, query.numberReturn)
	res = addString(res, query.query)
	res = addString(res, query.returnFieldSelector)
	length := len(res)
	setInt32(res, 0, int32(length))
	return res
}

func setInt32(res []byte, start int, length int32) {
	res[start] = byte(length)
	res[start+1] = byte(length >> 8)
	res[start+2] = byte(length >> 16)
	res[start+3] = byte(length >> 24)
}

func addInt32(bytes []byte, i int32) []byte{
	return append(bytes, byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24))
}

func addString(bytes []byte, s string) []byte{
	bytes = append(bytes, []byte(s)...)
	return append(bytes, byte(0))
}
