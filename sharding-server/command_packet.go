package sharding_server

import (
	"bytes"
	"fmt"
)

const(
	commandSleep   	=  0x00
	commandQuit   	=  0x01
	commandInitDb   	=  0x02
	commandQuery   	=  0x03
	commandFieldList   	=  0x04
	commandCreateDb   	=  0x05
	commandDropDb   	=  0x06
	commandRefresh   	=  0x07
	commandShutdown   	=  0x08
	commandStatistics   	=  0x09
	commandProcessInfo   	=  0x0a
	commandConnect   	=  0x0b
	commandProcessKill   	=  0x0c
	commandDebug    =  0x0d
	commandPing    =  0x0e
	commandTime    =  0x0f
	commandDelayedInsert    =  0x10
	commandChangeUser    =  0x11
	commandBinlogDump    =  0x12
	commandTableDump    =  0x13
	commandConnectOut    =  0x14
	commandRegisterSlave    =  0x15
	commandStmtPrepare    =  0x16
	commandStmtExecute    =  0x17
	commandStmtSendLongData    =  0x18
	commandStmtClose    =  0x19
	commandStmtReset    =  0x1a
	commandSetOption    =  0x1b
	commandStmtFetch    =  0x1c
	commandDaemon    =  0x1d
	commandBinlogDumpGtid    =  0x1e
	commandResetConnection    =  0x1f


	serverStatusInTrans            = 0x0001
   	serverStatusAutocommit         = 0x0002
   	serverMoreResultsExists        = 0x0008
   	serverStatusNoGoodIndexUsed    = 0x0010
   	serverStatusNoIndexUsed        = 0x0020
   	serverStatusCursorExists       = 0x0040
   	serverStatusLastRowSent        = 0x0080
   	serverStatusDbDropped          = 0x0100
   	serverStatusNoBackslashEscapes = 0x0200
   	serverStatusMetadataChanged    = 0x0400
   	serverQueryWasSlow             = 0x0800
   	serverPsOutParams              = 0x1000
   	serverStatusInTransReadonly    = 0x2000
   	serverSessionStateChanged      = 0x4000
)

type CommandPacket interface {
	Execute() ([]*MysqlPacket, error)
	Read(bytes.Buffer) error
}

type UnsupportedCommandPacket struct {
	commandType int
	*MysqlPacketHeader
}

func (packet *UnsupportedCommandPacket) Execute() ([]*MysqlPacket, error) {
	mysqlPacket, err := NewMysqlErrorPacket(packet.seq + 1,
		0xcc,
		"x",
		"xxx",
		fmt.Sprintf("Unsupported command packet '%s'.", packet.commandType),
			0xff).toMysqlPacket()

	if err != nil {
		return nil, err
	}
	return []*MysqlPacket{mysqlPacket}, nil
}

func (*UnsupportedCommandPacket) Read(bytes.Buffer) error {
	return nil
}

type QueryCommandPacket struct {

}

type QuitCommandPacket struct {

}

type InitDbCommandPacket struct {
	schemaName string
	*MysqlPacketHeader
}

func (packet *InitDbCommandPacket) Read(buf bytes.Buffer) error {
	packet.schemaName = string(buf.Bytes())
	return nil
}

func (packet *InitDbCommandPacket) Execute() ([]*MysqlPacket, error) {
	okPacket := NewMysqlOkPacket(packet.seq + 1, 0, 0, serverStatusAutocommit, 0, "")
	mysqlPacket, err := okPacket.toMysqlPacket()
	if err != nil {
		return nil, err
	}
	return []*MysqlPacket{mysqlPacket}, nil
}

type FieldListCommandPacket struct {

}

func Execute() ([]*MysqlPacket, error) {
	return nil, nil
}

func PacketFactory(command byte) CommandPacket {
	switch command {
	case commandQuery:
		return QueryCommandPacket{}
	case commandQuit:
		return QuitCommandPacket{}
	case commandInitDb:
		return InitDbCommandPacket{}
	case commandFieldList:
		return FieldListCommandPacket{}
	default:
		return &UnsupportedCommandPacket{commandType: int(command)}
	}
}