package sharding_server

import (
	"net"
	"bufio"
)

const (
	READ_BUF_SIZE = 16 * 1024
	WRITE_BUF_SIZE = 16 * 1024

	HEAD_PAYLOAD_LEN = 3
	HEAD_SEQ_LEN = 1
)

type ServerClient struct {
	conn net.Conn
	ss *ServerSession
}

type MysqlPacket struct {
	length int
	seq int
	body []byte
}

type ServerSession struct {
	conn net.Conn
	rb *bufio.Reader
	wb *bufio.Writer
	readCh chan MysqlPacket
	writeCh chan MysqlPacket
}

func read(r *bufio.Reader, size int) ([]byte, error){
	idx := 0
	buf := make([]byte, size)
	for {
		n, err := r.Read(buf[idx:])
		if err != nil {
			return nil, err
		}

		if n < size {
			idx += n
		}else {
			break
		}
	}
}

func (session *ServerSession) readPacket() {
	for {
		bts, err := read(session.rb, HEAD_PAYLOAD_LEN + HEAD_SEQ_LEN)
		if err != nil {
			session.close()
			break
		}


	}
}

func (session *ServerSession) writePacket() {

}

func (session *ServerSession) close() {
	session.conn.Close()
	session.wb.Flush()
	close(session.readCh)
	close(session.writeCh)
}

func (client ServerClient) start() {
	go client.ss.readPacket()
	go client.ss.writePacket()
}

func NewServerClient(conn net.Conn) *ServerClient{
	session := NewServerSession(conn)
	sc := ServerClient{
		conn: conn,
		ss: session,
	}

	return &sc
}

func NewServerSession(conn net.Conn) *ServerSession {
	ss := ServerSession{
		conn: conn,
		rb: bufio.NewReaderSize(conn, READ_BUF_SIZE),
		wb: bufio.NewWriterSize(conn, READ_BUF_SIZE),
	}
	return &ss
}

func Init() {
	listener, err := net.Listen("tcp", ":9998")
	if err != nil {
		panic(err)
	}

	for {
		conn, e := listener.Accept()
		if e != nil {
			panic(e)
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	NewServerClient(conn)
}
