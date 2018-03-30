package roomrequest

import (
	"lib/socket"
	"game/interfacer"
	"game/players"
	"protocol"
	"game/room"
	"code.google.com/p/goprotobuf/proto"
)

func init() {
	socket.Regist(&protocol.CReady{}, ready)
}

// 玩家准备
func ready(ctos *protocol.CReady, c interfacer.IConn) {
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	stoc := &protocol.SReady{}
	if rdata == nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotInRoom))
		c.Send(stoc)
		return
	}
	seat := player.GetSeat()
	ready := ctos.GetReady()
	err:= rdata.Readying(seat, ready)
	if err > 0 {
		stoc.Error = proto.Uint32(uint32(err))
		c.Send(stoc)
	}
}
