package queuenet

import (
	"net"
	"learngo/utils"
	"encoding/binary"
	"fmt"
	"bytes"
)

type TcpClient struct {
	addr string
	handler ConnHandler
}

type ClientHandler struct {

}



func NewClient(addr string) *TcpClient{
	return &TcpClient{addr:addr, handler:new(ClientHandler)}
}

func (client *TcpClient)StartClient(){
	dial, error := net.Dial("tcp", client.addr)
	utils.CheckError(error)

	client.handler.HandleAccept(dial)
}

func (handler *ClientHandler)HandleAccept(conn net.Conn) {
	go func() {
		for {
			utils.FuncNoError(handler.HandleRead, conn)
			//handler.HandleRead(conn)
		}
	}()
	go func() {
		for {
			utils.FuncNoError(handler.HandleWrite, conn)
			//handler.HandleWrite(conn)
		}
	}()
}

func (*ClientHandler)HandleWrite(conn net.Conn) {
	//by :=<- msg

	//var bf bytes.Buffer
	//bf.write
	//conn.Write(by)
}

func (*ClientHandler)HandleRead(conn net.Conn) {
	//by := make([]byte, 4)
	by := make([]byte, 4)
	n, err := conn.Read(by)
	utils.CheckError(err)
	if n <= 0 {
		return
	}

	fmt.Println(by)

	var buff bytes.Buffer
	_, err = buff.Write(by)
	utils.CheckError(err)

	var length int32
	binary.Read(&buff, binary.BigEndian, &length)

	by = make([]byte, length)
	n, err = conn.Read(by)
	utils.CheckError(err)
	fmt.Println("client receive msg :" , string(by[:]))
	//msg <- by
}

