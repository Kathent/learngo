package dial_ss

import "net"

func Dial_ss(){
	conn, err := net.Dial("tcp", "45.76.66.211:8388")
	if err != nil {
		panic(err)
	}


	conn.Write([]byte{1,1,1,1,1})
}
