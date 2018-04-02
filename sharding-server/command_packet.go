package sharding_server

import (
	"bytes"
	"fmt"
	"strings"
	"strconv"
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



	mysqlTypeDECIMAL = 0x00
	mysqlTypeTINY = 0x01
	mysqlTypeSHORT = 0x02
	mysqlTypeLONG	= 0x03
	mysqlTypeFLOAT	= 0x04
	mysqlTypeDOUBLE	= 0x05
	mysqlTypeNULL	= 0x06
	mysqlTypeTIMESTAMP	= 0x07
	mysqlTypeLONGLONG	= 0x08
	mysqlTypeINT24	= 0x09
	mysqlTypeDATE	= 0x0a
	mysqlTypeTIME	= 0x0b
	mysqlTypeDATETIME	= 0x0c
	mysqlTypeYEAR	= 0x0d
	mysqlTypeNEWDATE	= 0x0e
	mysqlTypeVARCHAR	= 0x0f
	mysqlTypeBIT	= 0x10
	mysqlTypeTIMESTAMP2	= 0x11
	mysqlTypeDATETIME2	= 0x12
	mysqlTypeTIME2	= 0x13
	mysqlTypeNEWDECIMAL	= 0xf6
	mysqlTypeENUM	= 0xf7
	mysqlTypeSET	= 0xf8
	mysqlTypeTINY_BLOB	= 0xf9
	mysqlTypeMEDIUM_BLOB	= 0xfa
	mysqlTypeLONG_BLOB	= 0xfb
	mysqlTypeBLOB	= 0xfc
	mysqlTypeVAR_STRING	= 0xfd
	mysqlTypeSTRING	= 0xfe
	mysqlTypeGEOMETRY	= 0xff
)

func getColumnType(str string) int {
	switch str {
	case "bit":
		return mysqlTypeBIT
	case "tiny":
		return mysqlTypeTINY
	case "int":
		return mysqlTypeINT24
	case "bigint":
		return mysqlTypeLONG
	case "decimal":
		return mysqlTypeDOUBLE
	case "char":
		return mysqlTypeSTRING
	case "varchar":
		return mysqlTypeVARCHAR
	case "date":
		return mysqlTypeDATE
	case "time":
		return mysqlTypeTIME
	case "timestamp":
		return mysqlTypeTIMESTAMP
	case "blob":
		return mysqlTypeBLOB
	default:
		return -1
	}
}

type CommandPacket interface {
	Execute() ([]*MysqlPacket, error)
	Read(bytes.Buffer) error
}

type UnsupportedCommandPacket struct {
	commandType int
	*MysqlPacketHeader
}

func (packet *UnsupportedCommandPacket) Execute() ([]*MysqlPacket, error) {
	errorPacket := NewMysqlErrorPacket(packet.seq+1,
		0xcc,
		"x",
		"xxx",
		fmt.Sprintf("Unsupported command packet '%d'.", packet.commandType),
		0xff)
	mysqlPacket, err := errorPacket.toMysqlPacket()

	if err != nil {
		return nil, err
	}
	return []*MysqlPacket{mysqlPacket}, nil
}

func (*UnsupportedCommandPacket) Read(bytes.Buffer) error {
	return nil
}

type QueryCommandPacket struct {
	sql string
	*MysqlPacketHeader
}

func (packet *QueryCommandPacket) Read(buf bytes.Buffer) error{
	packet.sql = string(buf.Bytes())
	return nil
}

func (packet *QueryCommandPacket) Execute() ([]*MysqlPacket, error) {
	return nil, nil
}

type QuitCommandPacket struct {
	*MysqlPacketHeader
}

func (packet *QuitCommandPacket) Read(buf bytes.Buffer) error{
	return nil
}

func (packet *QuitCommandPacket) Execute() ([]*MysqlPacket, error) {
	okPacket := NewMysqlOkPacket(packet.seq + 1, 0, 0, 0x002, 0, "")
	mysqlPacket, _ := okPacket.toMysqlPacket()
	return []*MysqlPacket{mysqlPacket}, nil
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
	table string
	fieldWildcard string
	*MysqlPacketHeader
}

func (packet *FieldListCommandPacket) Read(buf bytes.Buffer) error {
	line, err := buf.ReadBytes(0)
	if err != nil {
		return err
	}
	packet.table = string(line)
	packet.fieldWildcard = string(buf.Bytes()[len(line):])
	return nil
}

func (packet *FieldListCommandPacket) Execute() ([]*MysqlPacket, error) {
	sql := fmt.Sprintf("SHOW COLUMNS FROM %s FROM %s", packet.table, "robot")
	seq := packet.seq
	rows, err := GetDb().Query(sql)
	if err != nil {
		errPacket := NewMysqlErrorPacket(seq + 1, 0, "", "", err.Error(), 0xff)
		mysqlPacket, err := errPacket.toMysqlPacket()
		if err != nil {
			return nil, err
		}
		return []*MysqlPacket{mysqlPacket}, nil
	}

	res := make([]*MysqlPacket, 0)
	for rows.Next(){
		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		field := columns[1]
		lenStr := columns[2]
		columnLenStr := lenStr[strings.Index(lenStr[strings.Index(lenStr, ")"):], "(")+1:]
		realLen, _ := strconv.Atoi(columnLenStr)
		columnType := getColumnType(lenStr[strings.Index(lenStr, "(") + 1:])

		seq++
		tmpPacket := NewColumnDefinition41Packet(seq, "sharding_db", packet.table,
			packet.table, field, field, realLen, columnType, 0)
		mysqlPacket, err := tmpPacket.toMysqlPacket()
		if err != nil {
			return nil, err
		}
		res = append(res, mysqlPacket)
	}
	return res, nil
}

func Execute() ([]*MysqlPacket, error) {
	return nil, nil
}

func PacketFactory(command byte) CommandPacket {
	switch command {
	case commandQuery:
		return &QueryCommandPacket{}
	case commandQuit:
		return &QuitCommandPacket{}
	case commandInitDb:
		return &InitDbCommandPacket{}
	case commandFieldList:
		return &FieldListCommandPacket{}
	default:
		return &UnsupportedCommandPacket{commandType: int(command)}
	}
}