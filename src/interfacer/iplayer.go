/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 11:01
 * Filename      : iuser.go
 * Description   : 玩家自己的详细数据接口
 * *******************************************************/
package interfacer

import "protocol"

type IPlayer interface {
	GetUserid() string
	GetSeat() uint32
	GetNickname() string
	GetPhone() string

	SetUserid(string)
	SetNickname(string)
	SetSex(uint32)
	SetPwd(string)
	SetReady(bool)
	GetReady() bool

	SetLongitudeLatitude(longitude,latitude float32)
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

	GetBuild() string
	SetBuild(string)
	UserSave()
}
