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
	client *ServerClient
}

func (pipeline *Pipeline) HandleConnectionActive() {
	pipeline.GetCtx().foundInBoundContext().ConnectionActive()
}

func (pipeline *Pipeline) AddHandler(handler ConnectionHandler) {
	ctx := &ConnectionHandlerContext{}
	ctx.handler = handler
	ctx.client = pipeline.client

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
	log4go.Info("HandleMsgRead enter....")
	ctx := pipeline.GetCtx()
	if ctx != nil {
		ctx.ConnectionRead(buffer)
	}else {
		log4go.Info("HandleMsgRead ctx is nil....")
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
	client *ServerClient
}

func (context *ConnectionHandlerContext) ConnectionActive() {
	ctx := context.foundInBoundContext()
	if ctx != nil {
		ctx.handler.(ConnectionInBoundHandler).ConnectionActive(ctx)
	}else {
		log4go.Info("ConnectionActive no InBoundHandlerCtx found....")
	}
}

func (context *ConnectionHandlerContext) foundInBoundContext() *ConnectionHandlerContext {
	var tmp *ConnectionHandlerContext
	for tmp = context; tmp != nil ; {
		if _, ok := tmp.handler.(ConnectionInBoundHandler); ok {
			return tmp
		}else {
			tmp = tmp.next
		}
	}
	return nil
}

func (context *ConnectionHandlerContext) foundOutBoundContext() *ConnectionHandlerContext {
	var tmp *ConnectionHandlerContext
	for tmp = context; tmp != nil ; {
		if _, ok := tmp.handler.(ConnectionOutBoundHandler); ok {
			return tmp
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
	log4go.Info("ConnectionRead read enter...")
	ctx := context.foundInBoundContext()
	if ctx != nil {
		ctx.handler.(ConnectionInBoundHandler).ConnectionRead(ctx, buffer)
	}else {
		log4go.Info("ConnectionRead ctx is nil...")
	}
}

type ConnectionHandler interface {
	HandlerAdded(*ConnectionHandlerContext)
	HandlerRemoved(*ConnectionHandlerContext)
	ErrCaught(*ConnectionHandler, error)
}

type ConnectionInBoundHandler interface {
	ConnectionActive(ctx *ConnectionHandlerContext)
	ConnectionInActive(ctx *ConnectionHandlerContext)
	ConnectionRead(ctx *ConnectionHandlerContext, obj interface{})
}

type ConnectionOutBoundHandler interface {
	ConnectionWrite(ctx *ConnectionHandlerContext, obj interface{})
}

type BaseHandler struct {
	ConnectionHandler
}

func (BaseHandler) HandlerAdded(*ConnectionHandlerContext) {

}

func (BaseHandler) HandlerRemoved(*ConnectionHandlerContext) {

}

func (BaseHandler) ErrCaught(*ConnectionHandler, error) {

}

type MysqlCodecs struct {
	BaseInBoundHandler
	BaseOutBoundHandler
	BaseHandler
}

type BaseInBoundHandler struct {
	ConnectionInBoundHandler
}

func (*BaseInBoundHandler) ConnectionActive(ctx *ConnectionHandlerContext){
	ctx.next.ConnectionActive()
}

func (*BaseInBoundHandler) ConnectionInActive(ctx *ConnectionHandlerContext){
	ctx.next.ConnectionActive()
}

func (*BaseInBoundHandler) ConnectionRead(ctx *ConnectionHandlerContext, obj interface{}){
	ctx.next.ConnectionRead(obj)
}

type BaseOutBoundHandler struct {
	ConnectionOutBoundHandler
}

func (*BaseOutBoundHandler) ConnectionWrite(ctx *ConnectionHandlerContext, obj interface{}){

}

func (*MysqlCodecs) ConnectionWrite(ctx *ConnectionHandlerContext, obj interface{}) {
	if val, ok := obj.(*MysqlPacket); ok {
		ctx.client.ss.writeCh <- val
	}
}

func (*MysqlCodecs) ConnectionRead(ctx *ConnectionHandlerContext, obj interface{}) {
	log4go.Info("ConnectionRead read val:%+v", obj)
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

		log4go.Info("ConnectionRead read header %+v", header)
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
	BaseInBoundHandler
	BaseHandler
}

type MysqlPacketHandler struct {
	BaseInBoundHandler
	BaseHandler
}

func (*MysqlPacketHandler) ConnectionRead(ctx *ConnectionHandlerContext, obj interface{}) {
	if val, ok := obj.(*MysqlPacket); ok {
		buf := bytes.NewBuffer(val.body)
		commandBts := make([]byte, 1)
		n, err := buf.Read(commandBts)
		if err != nil {
			log4go.Info("ConnectionRead read command err:%v", err)
			return
		}

		if n < 1 {
			log4go.Info("ConnectionRead read command length wrong.")
			return
		}

		factory := PacketFactory(commandBts[0])
		packets, err := factory.Execute()
		if err != nil {
			log4go.Info("ConnectionRead execute err: %v", err)
			return
		}

		for _, v := range packets {
			ctx.Write(v)
		}
	}
}

type AuthPluginData struct {
	authPluginDataPart1 []byte
	authPluginDataPart2 []byte
	authPluginData []byte
}

func (*MysqlHandShakeHandler) ConnectionActive(ctx *ConnectionHandlerContext) {
	log4go.Info("MysqlHandShakeHandler ConnectionActive...")
	packet := NewMysqlPacket()
	marshal, err := NewMysqlHandShakePacket().marshal()
	if err != nil {
		log4go.Warn("ConnectionActive err:%v", err)
	}

	packet.body = marshal
	packet.length = len(marshal)
	ctx.Write(packet)
}

func NewMysqlPacket() *MysqlPacket {
	return &MysqlPacket{
		&MysqlPacketHeader{},
		nil,
	}
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
	log4go.Info("read size:%d", size)
	buf := make([]byte, size)
	n, err := r.Read(buf)
	if err != nil {
		return nil, err
	}
	log4go.Info("read n:%d", n)

	if n <= size {
		return buf[:n], nil
	}

	return buf, nil
}

func (session *ServerSession) readPacket() {
	var readBytes = 10240
	for {
		bts, err := read(session.rb, readBytes)
		if err != nil {
			log4go.Info("readPacket err:%v", err)
			session.close()
			break
		}else {
			log4go.Info("readPacket read packet len:%d", len(bts))
		}

		session.FireConnectionReadMsg(bytes.NewBuffer(bts))
	}
}

func (session *ServerSession) writePacket() {
	for val := range session.writeCh{
		session.wb.Write(utils.GetByte3(len(val.body)))
		session.wb.WriteByte(byte(val.seq))
		session.wb.Write(val.body)
		session.wb.Flush()
	}
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
		case server.maxConnChan <- 1:
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
	log4go.Info("Server handleConn....conn:%+v", conn)
	client := NewServerClient(conn)
	pipe := NewPipeline(client)
	client.ss.pipe = pipe
	server.handlerOp(pipe)
	client.start()
}

type HeadHandler struct{
	BaseInBoundHandler
}

func (*HeadHandler) HandlerAdded(*ConnectionHandlerContext){}
func (*HeadHandler)  HandlerRemoved(*ConnectionHandlerContext) {}
func (*HeadHandler) ErrCaught(*ConnectionHandler, error) {}

func (headHandler *HeadHandler) ConnectionActive(ctx *ConnectionHandlerContext) {
	log4go.Info("HeadHandler ConnectionActive enter...")
	go ctx.client.ss.readPacket()
	go ctx.client.ss.writePacket()
	ctx.next.ConnectionActive()
}

type TailHandler struct {
	BaseOutBoundHandler
}

func (*TailHandler) HandlerAdded(*ConnectionHandlerContext){}
func (*TailHandler)  HandlerRemoved(*ConnectionHandlerContext) {}
func (*TailHandler) ErrCaught(*ConnectionHandler, error) {}

func NewPipeline(client *ServerClient) *Pipeline {
	var pipe = Pipeline{
		handlers: list.New(),
		client: client,
	}
	pipe.AddHandler(&HeadHandler{})
	pipe.AddHandler(&TailHandler{})

	return &pipe
}

func (server *Server) HandleOp(f func(pipeline *Pipeline)) *Server{
	server.handlerOp = f
	return server
}

func Init() {
	NewServer(100).HandleOp(func(pipeline *Pipeline){
		pipeline.AddHandler(&MysqlCodecs{BaseInBoundHandler{}, BaseOutBoundHandler{}, BaseHandler{}})
		pipeline.AddHandler(&MysqlHandShakeHandler{BaseInBoundHandler{}, BaseHandler{}})
		pipeline.AddHandler(&MysqlPacketHandler{BaseInBoundHandler{}, BaseHandler{}})
	}).Start()
}
