package request

import (
	"lib/socket"
	"game/data"

	"game/interfacer"
	"protocol"

	"github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
	"game/players"
	"game/room"
)

func init() {

	socket.Regist(&protocol.CUserData{}, getUserDataHdr)

	socket.Regist(&protocol.CArchieve{}, getArchieve)

}

func getUserDataHdr(ctos *protocol.CUserData, c interfacer.IConn) {
	stoc := &protocol.SUserData{}
	stoc.Data = &protocol.UserData{}
	stoc.Data.Roomid = proto.Uint32(0)
	stoc.Data.Invitecode = proto.String("")
	stoc.Data.Roomtype = proto.Uint32(0)
	// 获取玩家自己的详细资料
	if ctos.GetUserid() == c.GetUserid() {
		userdata := players.Get(c.GetUserid())
		stoc.Data = userdata.ConverDataUser()

		user:= &data.User{Userid:c.GetUserid()}
		photo,_:= user.GetPhotoFromDB()
		stoc.Data.Photo = proto.String(photo)
		// 再次认证是否在房间并判断房间是否过期
		if userdata.GetInviteCode() != "" && room.Get(userdata.GetInviteCode()) != nil {
			//glog.Infoln("已经在私人局", userdata.GetInviteCode())
			stoc.Data.Roomtype = proto.Uint32(userdata.GetRoomType())
			stoc.Data.Invitecode = proto.String(userdata.GetInviteCode())
			stoc.Data.Roomid = proto.Uint32(userdata.GetRoomID())
		} else {
			if userdata.GetInviteCode() != "" {
				glog.Infoln("不在私人局，或者房间过期", userdata.GetInviteCode())
				userdata.ClearRoom()
			}
		}
	} else {
		if ctos.GetUserid() != "" {
			member := &data.User{Userid: ctos.GetUserid()}
			err := member.Get()
			if err != nil {
				// TODO : err = err: No such field: %s in obj matchid
				// 老字段不能删除
				stoc.Error = proto.Uint32(uint32(protocol.Error_UserDataNotExist))
			} else {

				stoc.Data.Roomcard = &member.RoomCard
				stoc.Data.Userid = &member.Userid
				stoc.Data.Sex = &member.Sex
				stoc.Data.Nickname = &member.Nickname
				stoc.Data.Email = &member.Email
				stoc.Data.Phone = &member.Phone
				stoc.Data.Photo = &member.Photo
				stoc.Data.Ip = &member.Create_ip
				stoc.Data.Createtime = &member.Create_time
				stoc.Data.Terminal = &member.Terminal
				stoc.Data.Platform = &member.Platform
			}
		} else {
			stoc.Error = proto.Uint32(uint32(protocol.Error_UsernameEmpty))
		}
	}
	c.Send(stoc)
}

// 获取玩家的战绩
func getArchieve(ctos *protocol.CArchieve, c interfacer.IConn) {
	stoc := &protocol.SArchieve{}
	if ctos.GetUserid() == "" {
		stoc.Error = proto.Uint32(uint32(protocol.Error_UserDataNotExist))
	}
	c.Send(stoc)
}
