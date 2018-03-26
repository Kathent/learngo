package sharding_server

import (
	"net"
	"bufio"
	"github.com/alecthomas/log4go"
	"container/list"
)

const (
	READ_BUF_SIZE = 16 * 1024
	WRITE_BUF_SIZE = 16 * 1024

	HEAD_PAYLOAD_LEN = 3
	HEAD_SEQ_LEN = 1
)

type Server struct {
	Pipeline Pipeline
	maxConnChan chan byte
	closeChan chan byte
}

func NewServer(maxConNum int) *Server{
	return &Server{
		maxConnChan: make(chan byte, maxConNum),
		Pipeline: Pipeline{},
		closeChan: make(chan byte, 0),
	}
}

type ServerClient struct {
	conn *net.TCPConn
	ss *ServerSession
}

type MysqlPacketHeader struct {
	length int
	seq int
}

type MysqlPacket struct {
	*MysqlPacketHeader
	body []byte
}

type Pipeline struct {
	handlers *list.List
}

type ConnectionHandlerContext struct {
	handler ConnectionHandler
	pre *ConnectionHandlerContext
	next *ConnectionHandlerContext
}

type ConnectionHandler interface {
	HandlerAdded(*ConnectionHandlerContext)
	HandlerRemoved(*ConnectionHandlerContext)
	ErrCaught(*ConnectionHandler, error)
}

type ConnectionInBoundHandler interface {
	ConnectionHandler
	ConnectionActive(ctx *ConnectionHandlerContext)
	ConnectionInActive(ctx *ConnectionHandlerContext)
	ConnectionRead(ctx *ConnectionHandlerContext, obj interface{})
}

type ConnectionOutBoundHandler interface {
	ConnectionHandler
	ConnectionWrite(ctx *ConnectionHandlerContext)
}



func unmarshal(bts []byte) (*MysqlPacketHeader, error){
	_ = bts[3]
	header := &MysqlPacketHeader{}
	header.length = int(bts[0]) | int(bts[1]) << 8 | int(bts[2]) << 16
	header.seq = int(bts[3])
	return header, nil
}

type ServerSession struct {
	conn *net.TCPConn
	rb *bufio.Reader
	wb *bufio.Writer
	readCh chan *MysqlPacket
	writeCh chan *MysqlPacket
}

func read(r *bufio.Reader, size int) ([]byte, error){
	idx := 0
	buf := make([]byte, size)
	for {
		n, err := r.Read(buf[idx:])
		if err != nil {
			return nil, err
		}
		log4go.Info("read n:%d", n)

		if n < size {
			idx += n
		}else {
			break
		}
	}

	return buf, nil
}

func (session *ServerSession) readPacket() {
	for {
		bts, err := read(session.rb, HEAD_PAYLOAD_LEN + HEAD_SEQ_LEN)
		if err != nil {
			log4go.Info("readPacket err:%v", err)
			session.close()
			break
		}

		header, err := unmarshal(bts)
		if err != nil {
			log4go.Info("unmarshal err:%v", err)
			continue
		}

		bytes, err := read(session.rb, header.length)
		if err != nil {
			log4go.Info("read body:%v", err)
			session.close()
			break
		}

		packet := MysqlPacket{header, bytes}

		log4go.Info("readPacket packet:%+v", packet)
		session.readCh <- &packet
	}
}

func (session *ServerSession) writePacket() {
	// for val := range session.writeCh{
	//
	// }
}

func (session *ServerSession) close() {
	session.conn.Close()
	session.wb.Flush()
	close(session.readCh)
	close(session.writeCh)
}

func (session *ServerSession) dispatchCommand() {
	// for val := range session.readCh {
	//
	// }
}

func (client ServerClient) start() {
	log4go.Info("ServerClient start client:%+v", client)
	go client.ss.readPacket()
	go client.ss.writePacket()
	go client.ss.dispatchCommand()
}

func NewServerClient(conn *net.TCPConn) *ServerClient{
	session := NewServerSession(conn)
	sc := ServerClient{
		conn: conn,
		ss: session,
	}

	return &sc
}

func NewServerSession(conn *net.TCPConn) *ServerSession {
	ss := ServerSession{
		conn: conn,
		rb: bufio.NewReaderSize(conn, READ_BUF_SIZE),
		wb: bufio.NewWriterSize(conn, READ_BUF_SIZE),
		readCh: make(chan *MysqlPacket, 1000),
		writeCh: make(chan *MysqlPacket, 1000),
	}
	return &ss
}

func (server *Server) Start() {
	addr, err := net.ResolveTCPAddr("tcp", ":9998")
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <- server.closeChan:
			server.Close()
		case <- server.maxConnChan:
			conn, err := listener.AcceptTCP()
			if err != nil {
				server.closeChan <- 1
				break
			}

			handleConn(conn)
		}
	}
}

func (server *Server) Close() {

}

func handleConn(conn *net.TCPConn) {
	NewServerClient(conn).start()
}

func Init() {
	NewServer(100).Start()
}
