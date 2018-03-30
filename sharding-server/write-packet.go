package sharding_server

import (
	"learngo/utils"
	"github.com/alecthomas/log4go"
)

type MysqlWritePacket interface{
	marshal() ([]byte, error)
}

type BaseMysqlWritePacket struct {
	*MysqlPacketHeader
	MysqlWritePacket
}

func (*BaseMysqlWritePacket) marshal() ([]byte, error) {
	return nil, nil
}

func (base *BaseMysqlWritePacket) toMysqlPacket() (*MysqlPacket, error) {
	packet := NewMysqlPacket()
	packet.seq = base.seq + 1
	bts, err := base.marshal()
	if err != nil {
		log4go.Info("toMysqlPacket err:%v", err)
		return nil, err
	}

	packet.body = bts
	return packet, nil
}

type MysqlErrorPacket struct {
	*BaseMysqlWritePacket

	header int
	errorCode int
	sqlStateMarker string
	sqlState	string
	errorMessage	string
}

func NewMysqlErrorPacket(seq int, errCode int, stateMarker string,
	sqlState string, msg string, header int) *MysqlErrorPacket {
	packet := &MysqlErrorPacket{
		header:         header,
		errorCode:      errCode,
		sqlStateMarker: stateMarker,
		sqlState:       sqlState,
		errorMessage:   msg,
	}

	packet.seq = seq
	return packet
}

func (err *MysqlErrorPacket) marshal() ([]byte, error) {
	bts := make([]byte, 0)
	bts = append(bts, byte(err.header))
	bts = append(bts, utils.GetByte2(err.errorCode)...)
	bts = append(bts, []byte(err.sqlStateMarker)...)
	bts = append(bts, []byte(err.sqlState)...)
	bts = append(bts, utils.GetStringNul(err.errorMessage)...)
	return bts, nil
}

type MysqlOkPacket struct {
	*BaseMysqlWritePacket
	header int
	affectedRows int64
	lastInsertId int64
	statusFlags  int
	warnings	 int
	info		 string
}

func NewMysqlOkPacket(seq int, row int64, id int64, flags int, warnings int, info string) *MysqlOkPacket {
	packet := &MysqlOkPacket{
		header: 0x00,
		affectedRows: row,
		lastInsertId: id,
		statusFlags: flags,
		warnings: warnings,
		info: info,
	}

	packet.seq = seq
	return packet
}

func (packet *MysqlOkPacket) marshal() ([]byte, error) {
	bts := make([]byte, 0)
	bts = append(bts, byte(packet.header))
	bts = append(bts, utils.GetLongToInt(packet.affectedRows)...)
	bts = append(bts, utils.GetLongToInt(packet.lastInsertId)...)
	bts = append(bts, utils.GetByte2(packet.statusFlags)...)
	bts = append(bts, utils.GetByte2(packet.warnings)...)
	bts = append(bts, utils.GetStringNul(packet.info)...)
	return bts, nil
}

type MysqlHandShakePacket struct {
	MysqlWritePacket
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
