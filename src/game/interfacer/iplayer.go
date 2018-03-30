package interfacer

import "protocol"

type IPlayer interface {
	GetUserid() string
	GetSeat() uint32
	SetUserid(string)
	SetLongitudeLatitude(longitude,latitude string)
	ConverDataUser() *protocol.UserData
	ConverProtoUser() *protocol.ProtoUser
	GetInviteCode() string // 私人局邀请码
	GetRoomType() uint32   // 房间类型ID,对应房间表
	GetRoomID() uint32     // 比赛场或金币场房间id
	// 分别为：房间类型ID，房间号，房间邀请码
	SetRoom( uint32, uint32, string)
	ClearRoom()
	GetPlatform() uint32

	SetConn(IConn)
	GetConn() IConn
	Send(IProto)

	GetRoomCard() uint32
	SetRoomCard(uint32)
}
