package queuenet

import (
	"net"
	"utils"
	"encoding/binary"
	"fmt"
	"time"
	"bytes"
)

type TcpServer struct{
	addr string
	handler ConnHandler
}

type ConnHandler interface {
	HandleAccept(conn net.Conn)
	HandleWrite(conn net.Conn)
	HandleRead(conn net.Conn)
}

type ServerConnHandler struct {

}

func NewServer(addr string) *TcpServer{
	return &TcpServer{addr:addr, handler:new(ServerConnHandler)}
}

var msg = make(chan []byte)

func (server *TcpServer)StartServer() error{
	listen, err := net.Listen("tcp", server.addr)
	if err != nil {
		return err
	}

	accept, err := listen.Accept()
	if err != nil {
		return err
	}
	server.handler.HandleAccept(accept)
	return nil
}

func (handler *ServerConnHandler)HandleAccept(conn net.Conn) {
	go func() {
		for {
			//utils.FuncNoError(handler.HandleRead, conn)
			handler.HandleRead(conn)
		}
	}()
	go func() {
		for {
			//utils.FuncNoError(handler.HandleWrite, conn)
			handler.HandleWrite(conn)
		}
	}()
}

func (*ServerConnHandler)HandleWrite(conn net.Conn) {
	tick := time.Tick(1e9)
	<-tick
	t := time.Now()
	var buff bytes.Buffer
	n, err := buff.WriteString(t.String())

	utils.CheckError(err)

	fmt.Println(n)

	var lb bytes.Buffer

	length := make([]byte, 4)

	m := uint32(n)

	binary.BigEndian.PutUint32(length, m)

	lb.Write(length)

	lb.Write(buff.Bytes())

	conn.Write(lb.Bytes())
}

func (*ServerConnHandler)HandleRead(conn net.Conn) {
	by := make([]byte, 32)
	n, err := conn.Read(by)
	utils.CheckError(err)
	if n <= 0 {
		return
	}

	var buff bytes.Buffer
	_, err = buff.Read(by)
	utils.CheckError(err)

	var length int32
	binary.Read(&buff, binary.BigEndian, &length)

	by = make([]byte, length)
	n, err = conn.Read(by)
	utils.CheckError(err)
	fmt.Println("server receive msg :" , by)
	msg <- by
}
