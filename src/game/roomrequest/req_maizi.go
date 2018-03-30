package roomrequest

import (
	"protocol"
	"game/interfacer"
	"game/players"
	"code.google.com/p/goprotobuf/proto"
	"game/room"
	"lib/socket"
	"config"
)

func init() {
	socket.Regist(&protocol.CMaiZi{}, maizi)
}

func maizi(ctos *protocol.CMaiZi, c interfacer.IConn) {
	stoc := &protocol.SMaiZi{}
	player := players.Get(c.GetUserid())
	if player == nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_UserDataNotExist))
		c.Send(stoc)
		return
	}
	rdata := room.Get(player.GetInviteCode())
	if rdata == nil {
		player.ClearRoom()
		stoc.Error = proto.Uint32(uint32(protocol.Error_RoomNotExist))
		c.Send(stoc)
		return
	}

	exist := false
	for _,v:=range config.Opts().SelectMaiziCt{
		if v == stoc.GetCount(){
			exist = true
			break
		}
	}

	if !exist {
		stoc.Error = proto.Uint32(uint32(protocol.Error_GameMaiziIllegal))
		c.Send(stoc)
		return
	}

	rel := rdata.MaiZi(player.GetSeat(), ctos.GetCount())
	if rel > 0 {
		stoc.Error = proto.Uint32(uint32(rel))
		c.Send(stoc)
		return
	}
}
