package roomrequest

import (
	"game/room"
	"code.google.com/p/goprotobuf/proto"
	"protocol"
	"game/interfacer"
	"game/players"
	"lib/socket"
	"github.com/golang/glog"
)

func init() {
	socket.Regist(&protocol.CEnterSocialRoom{}, entryroom)
}

func entryroom(ctos *protocol.CEnterSocialRoom, c interfacer.IConn) {
	stoc := &protocol.SEnterSocialRoom{}
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	if rdata != nil { //已经存在或重复进入
		code := rdata.Enter(player)
		if code == 0 {
			return
		}
	}
	player.ClearRoom()
	rdata = room.Get(ctos.GetInvitecode())
	if rdata == nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_RoomNotExist))
		c.Send(stoc)
		return
	}
	glog.Infoln("进入房间：",ctos.String())
	player.SetLongitudeLatitude(ctos.GetLongitude(),ctos.GetLatitude())
	code := rdata.Enter(player)

	if code > 0 {
		stoc.Error = proto.Uint32(uint32(code))
		c.Send(stoc)
	}
}
