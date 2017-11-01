package mango_learn

import (
	"net"
	"fmt"
)

func LearnMango(){
	conn, err := net.Dial("tcp", "localhost:27017")
	if err != nil {
		panic(err)
	}

	bts := getQuery()
	conn.Write(bts)


	readBt := make([]byte, 1)
	for {
		_, err := conn.Read(readBt)
		if err != nil {
			panic(err)
		}

		fmt.Println(fmt.Sprintf("%v", string(readBt)))
	}
}

func getQuery() []byte{
	res := make([]byte, 0)
	res = addInt32(res, 0)
	res = addInt32(res, 1)
	res = addInt32(res, 0)
	res = addInt32(res, 2004)

	res = addInt32(res, 0)
	res = addString(res, "local.student")
	res = addInt32(res, 0)
	res = addInt32(res, 1)
	res = addString(res, "\"code\":\"123\"")
	res = addString(res, "local.student")
	res = setInt32(res, 0, int32(len(res)))
	return res
}

func setInt32(bytes []byte, start int, val int32) []byte {
	bytes[start], bytes[start+1], bytes[start+2], bytes[start+3] =
		byte(val), byte(val >> 8), byte(val >> 16), byte(val >> 24)
	return bytes
}

func addInt32(bytes []byte, i int) []byte {
	return append(bytes, byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24))
}

func addString(bytes []byte, val string) []byte {
	bytes = append(bytes, []byte(val)...)
	return append(bytes, byte(0))
}
