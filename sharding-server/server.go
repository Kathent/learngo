package sharding_server

import (
	"net"
	"bufio"
	"github.com/alecthomas/log4go"
	"container/list"
	"bytes"
	"learngo/utils"
)

const (
	READ_BUF_SIZE = 16 * 1024
	WRITE_BUF_SIZE = 16 * 1024

	HEAD_PAYLOAD_LEN = 3
	HEAD_SEQ_LEN = 1

	READ_SIZE = 1024
)

type Server struct {
	maxConnChan chan byte
	closeChan chan byte
	handlerOp func(*Pipeline)
}

func NewServer(maxConNum int) *Server{
	return &Server{
		maxConnChan: make(chan byte, maxConNum),
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

func (pipeline *Pipeline) HandleConnectionActive() {
	front := pipeline.handlers.Front()
	if front == nil {
		return
	}

	if val, ok := front.Value.(*ConnectionHandlerContext); ok {
		val.ConnectionActive()
	}
}

func (pipeline *Pipeline) AddHandler(handler ConnectionHandler) {
	ctx := &ConnectionHandlerContext{}
	ctx.handler = handler

	back := pipeline.handlers.Back()
	if back == nil {
		pipeline.handlers.PushBack(ctx)
	}else if val, ok := back.Value.(*ConnectionHandlerContext); ok {
		ctx.pre = val
		val.next = ctx
		pipeline.handlers.PushBack(ctx)
	}
}
func (pipeline *Pipeline) HandleMsgRead(buffer *bytes.Buffer) {
	ctx := pipeline.GetCtx()
	if ctx != nil {
		ctx.ConnectionRead(buffer)
	}
}

func (pipeline *Pipeline) GetCtx() *ConnectionHandlerContext{
	front := pipeline.handlers.Front()
	if front == nil {
		return nil
	}

	return front.Value.(*ConnectionHandlerContext)
}

type ConnectionHandlerContext struct {
	handler ConnectionHandler
	pre *ConnectionHandlerContext
	next *ConnectionHandlerContext
	client ServerClient
}

func (context *ConnectionHandlerContext) ConnectionActive() {
	ctx := context.foundInBoundContext()
	if ctx != nil {
		ctx.handler.(ConnectionInBoundHandler).ConnectionActive(ctx)
	}
}

func (context *ConnectionHandlerContext) foundInBoundContext() *ConnectionHandlerContext {
	var tmp *ConnectionHandlerContext
	for tmp = context; tmp != nil ; {
		if _, ok := tmp.handler.(ConnectionInBoundHandler); ok {
			return context
		}else {
			tmp = tmp.next
		}
	}
	return nil
}

func (context *ConnectionHandlerContext) foundOutBoundContext() *ConnectionHandlerContext {
	var tmp *ConnectionHandlerContext
	for tmp = context; tmp != nil ; {
		if _, ok := tmp.handler.(ConnectionInBoundHandler); ok {
			return context
		}else {
			tmp = tmp.pre
		}
	}
	return nil
}

func (context *ConnectionHandlerContext) Write(packet interface{}) {
	ctx := context.foundOutBoundContext()
	if ctx != nil {
		ctx.handler.(ConnectionOutBoundHandler).ConnectionWrite(ctx, packet)
	}
}

func (context *ConnectionHandlerContext) ConnectionRead(buffer interface{}) {
	ctx := context.foundInBoundContext()
	if ctx != nil {
		ctx.handler.(ConnectionInBoundHandler).ConnectionRead(ctx, buffer)
	}
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
	ConnectionWrite(ctx *ConnectionHandlerContext, obj interface{})
}

type MysqlCodecs struct {
	ConnectionInBoundHandler
	ConnectionOutBoundHandler
}

func (*MysqlCodecs) ConnectionWrite(ctx *ConnectionHandlerContext, obj interface{}) {
	if val, ok := obj.(*MysqlPacket); ok {
		ctx.client.ss.writeCh <- val
	}
}

func (*MysqlCodecs) ConnectionRead(ctx *ConnectionHandlerContext, obj interface{}) {
	if val, ok := obj.(*bytes.Buffer); ok {
		bts := make([]byte, HEAD_PAYLOAD_LEN + HEAD_SEQ_LEN)

		n, err := val.Read(bts)
		if err != nil {
			ctx.client.ss.close()
			return
		}

		if n < HEAD_PAYLOAD_LEN + HEAD_SEQ_LEN {
			log4go.Info("MysqlCodecs ConnectionRead read length wrong..")
			return
		}

		header, err := unmarshal(bts)
		if err != nil {
			log4go.Info("unmarshal err:%v", err)
			return
		}

		bodyBts := make([]byte, header.length)
		n, err = val.Read(bodyBts)
		if err != nil {
			log4go.Info("read body:%v", err)
			ctx.client.ss.close()
		}

		if n < header.length {
			//读取到的长度不够
			return
		}

		packet := MysqlPacket{header, bodyBts}
		log4go.Info("readPacket packet:%+v", packet)

		if ctx.next != nil {
			ctx.next.ConnectionRead(packet)
		}
	}
}

type MysqlHandShakeHandler struct {
	ConnectionInBoundHandler
}

type AuthPluginData struct {
	authPluginDataPart1 []byte
	authPluginDataPart2 []byte
	authPluginData []byte
}

type MysqlHandShakePacket struct {
	protoVersion int
	serverVersion string
	capabilityFlagsLower int
	characterSet int
	statusFlag int
	capabilityFlagsUpper int
	connectionId int
	authPluginData AuthPluginData
}

func (packet *MysqlHandShakePacket) marshal() ([]byte, error) {
	res := make([]byte, 0)
	res = append(res, byte(packet.protoVersion))
	res = append(res, utils.GetStringNul(packet.serverVersion)...)
	res = append(res, byte(0))
	res = append(res, utils.GetByte4(packet.connectionId)...)
	res = append(res, utils.GetStringNul(string(packet.authPluginData.authPluginDataPart1))...)
	res = append(res, utils.GetByte2(packet.capabilityFlagsLower)...)
	res = append(res, byte(packet.characterSet))
	res = append(res, utils.GetByte2(packet.statusFlag)...)
	res = append(res, utils.GetByte2(packet.capabilityFlagsUpper)...)
	res = append(res, byte(0))
	res = append(res, []byte{0,0,0,0,0,0,0,0}...)
	res = append(res, utils.GetStringNul(string(packet.authPluginData.authPluginDataPart2))...)

	return res, nil
}

func NewMysqlHandShakePacket() *MysqlHandShakePacket{
	return &MysqlHandShakePacket{
		protoVersion: 0x0A,
		serverVersion: "5.5.59-Sharding-Proxy 2.1.0",
		capabilityFlagsLower: 0x1 | 0x2 | 0x4 | 0x8 | 0x40 | 0x100 | 0x200 | 0x400 | 0x1000 | 0x2000 | 0x8000,
		characterSet: 0x21,
		statusFlag: 0x2,
		capabilityFlagsUpper: 0,
		connectionId: 0,
		authPluginData: AuthPluginData{
			authPluginDataPart1: []byte{0,0,0,0,0,0,0,0},
			authPluginDataPart2: []byte{0,0,0,0,0,0,0,0,0,0,0,0},
			authPluginData: []byte{0, 0, 0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0},
		},
	}
}

func (MysqlHandShakeHandler) ConnectionActive(ctx *ConnectionHandlerContext) {
	ctx.Write(MysqlHandShakePacket{})
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
	pipe *Pipeline
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
	var readBytes = 1024
	for {
		bts, err := read(session.rb, readBytes)
		if err != nil {
			log4go.Info("readPacket err:%v", err)
			session.close()
			break
		}

		session.FireConnectionReadMsg(bytes.NewBuffer(bts))
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

}
func (session *ServerSession) FireConnectionReadMsg(buffer *bytes.Buffer) {
	session.pipe.HandleMsgRead(buffer)
}

func (client ServerClient) start() {
	log4go.Info("ServerClient start client:%+v", client)
	client.conn.SetKeepAlive(true)
	client.conn.SetNoDelay(true)
	client.conn.SetReadBuffer(READ_BUF_SIZE)
	client.conn.SetWriteBuffer(WRITE_BUF_SIZE)
	client.ss.pipe.HandleConnectionActive()
}

func NewServerClient(conn *net.TCPConn, pipe *Pipeline) *ServerClient{
	session := NewServerSession(conn, pipe)
	sc := ServerClient{
		conn: conn,
		ss: session,
	}

	return &sc
}

func NewServerSession(conn *net.TCPConn, pipe *Pipeline) *ServerSession {
	ss := ServerSession{
		conn: conn,
		rb: bufio.NewReaderSize(conn, READ_BUF_SIZE),
		wb: bufio.NewWriterSize(conn, READ_BUF_SIZE),
		readCh: make(chan *MysqlPacket, 1000),
		writeCh: make(chan *MysqlPacket, 1000),
		pipe: pipe,
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

			go server.handleConn(conn)
		}
	}
}

func (server *Server) Close() {

}

func (server *Server) handleConn(conn *net.TCPConn) {
	pipe := NewPipeline()
	server.handlerOp(pipe)
	NewServerClient(conn, pipe).start()
}

type HeadHandler struct{
	ConnectionInBoundHandler
}

func (headHandler *HeadHandler) ConnectionActive(ctx *ConnectionHandlerContext) {
	go ctx.client.ss.readPacket()
	go ctx.client.ss.writePacket()
	if ctx.next != nil {
		ctx.next.ConnectionActive()
	}
}

type TailHandler struct {
	ConnectionOutBoundHandler
}

func NewPipeline() *Pipeline {
	var pipe = Pipeline{
		handlers: list.New(),
	}
	pipe.AddHandler(HeadHandler{})
	pipe.AddHandler(TailHandler{})

	return &pipe
}

func (server *Server) HandleOp(f func(pipeline *Pipeline)) *Server{
	server.handlerOp = f
	return server
}

func Init() {
	NewServer(100).HandleOp(func(pipeline *Pipeline){
		pipeline.AddHandler(MysqlCodecs{})
		pipeline.AddHandler(MysqlHandShakeHandler{})
	}).Start()
}
