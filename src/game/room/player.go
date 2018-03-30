package room

import (
	//"algorithm"
	"code.google.com/p/goprotobuf/proto"
	"game/data"
	"game/interfacer"
	"protocol"
	"sync"
)

func NewPlayer(data *data.User) *Player {
	return &Player{
		user: data,
	}
}

//用户的全部数据，服务器存取专用
type Player struct {
	seat uint32 // 玩家座位
	sync.RWMutex
	user *data.User
	conn interfacer.IConn

	inviteCode string
	roomType   uint32
	roomID     uint32 // 比赛场或私人局的房间ID
}

func (this *Player) Send(value interfacer.IProto) {
	this.RLock()
	defer this.RUnlock()
	this.conn.Send(value)
}
func (this *Player) GetConn() interfacer.IConn {
	this.RLock()
	defer this.RUnlock()
	return this.conn
}
func (this *Player) SetConn(value interfacer.IConn) {
	this.Lock()
	defer this.Unlock()
	this.conn = value
}

// ---------------------------房间属性---------------------------------
func (this *Player) ClearRoom() {
	this.inviteCode = ""
	this.roomType = 0
	this.roomID = 0
	this.seat = 0
}

// 分别为：房间类型ID，房间号，房间邀请码
func (this *Player) SetRoom(rid uint32, seat uint32, invitecode string) {
	this.inviteCode = invitecode
	this.roomID = rid
	this.seat = seat
}

func (this *Player) GetInviteCode() string {
	return this.inviteCode
}

func (this *Player) GetRoomID() uint32 {
	return this.roomID
}

func (this *Player) GetRoomType() uint32 {
	return this.roomType
}
func (this *Player) SetUserid(id string) { this.user.Userid = id }
func (this *Player) GetSeat() uint32     { return this.seat }
func (this *Player) GetUserid() string   { return this.user.Userid }

func (this *Player) SetLongitudeLatitude(longitude, latitude string) {
	this.user.Longitude = longitude
	this.user.Latitude = latitude
}

func (this *Player) SetEmail(email string) { this.user.Email = email }

func (this *Player) SetAuth(auth string) { this.user.Auth = auth }

func (this *Player) SetPwd(pwd string) { this.user.Pwd = pwd }

func (this *Player) SetIP(ip uint32) { this.user.Create_ip = ip }

func (this *Player) SetRoomCard(value uint32) { this.user.RoomCard = value }
func (this *Player) GetRoomCard() uint32      { return this.user.RoomCard }

func (this *Player) SetStatus(status uint32) { this.user.Status = status }
func (this *Player) GetPlatform() uint32     { return this.user.Platform }
func (this *Player) ConverProtoUser() *protocol.ProtoUser {
	return &protocol.ProtoUser{
		Userid:    &this.user.Userid,
		Position:  &this.seat,
		Nickname:  &this.user.Nickname,
		Sex:       &this.user.Sex,
		Photo:     &this.user.Photo,
		Address:   &this.user.Address,
		Terminal:  &this.user.Terminal,
		Email:     &this.user.Email,
		Roomcard:  &this.user.RoomCard,
		Platform:  &this.user.Platform,
		Longitude: &this.user.Longitude,
		Latitude:  &this.user.Latitude,
	}
}

func (this *Player) ConverDataUser() *protocol.UserData {
	return &protocol.UserData{
		Userid:     &this.user.Userid,
		Nickname:   &this.user.Nickname,
		Sex:        &this.user.Sex,
		Photo:      &this.user.Photo,
		Status:     &this.user.Status,
		Online:     proto.Bool(true),
		Phone:      &this.user.Phone,
		Address:    &this.user.Address,
		Terminal:   &this.user.Terminal,
		Email:      &this.user.Email,
		Ip:         &this.user.Create_ip,
		Roomid:     &this.roomID,
		Createtime: &this.user.Create_time,
		Platform:   &this.user.Platform,
		Roomcard:   &this.user.RoomCard,
		Build:      &this.user.Build,
	}
}
