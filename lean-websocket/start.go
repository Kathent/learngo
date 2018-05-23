package main

import (
	"net"
	"net/http"
	"fmt"
	"time"
	"github.com/gorilla/websocket"
)

func main() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(3 * time.Second)
		dialer := &websocket.Dialer{NetDial: func(network, addr string) (net.Conn, error) {
			return net.Dial(network, addr)
		}}
		conn, _, err := dialer.Dial("ws://localhost:12345/ws", nil)
		if err != nil {
			panic(err)
		}

		go readMsg(conn)
		go writeMsg(conn)
	}()

	mx := http.NewServeMux()
	mx.HandleFunc("/ws", wsHandler)
	http.Serve(listener, mx)
}


func wsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("enter wsHandler....")
	upgrader := &websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(fmt.Sprintf("upgrade faill...err:%v", err))
		return
	}


	go readMsg(conn)

	go writeMsg(conn)

	go func() {
		time.Sleep(time.Second * 10)
	}()
}

func writeMsg(conn *websocket.Conn) {
	func() {
		data := struct {
			A int
			B string
		}{
			A: 1,
			B: "2",
		}

		for {
			err := conn.WriteJSON(&data)
			if err != nil {
				fmt.Println(fmt.Sprintf("read msg err:%v", err))
				break
			}
		}
	}()
}

func readMsg(conn *websocket.Conn) {
	func() {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(fmt.Sprintf("read msg err:%v", err))
				break
			}

			fmt.Println(fmt.Sprintf("msg type:%d, content:%s", messageType, string(p)))
		}
	}()
}
