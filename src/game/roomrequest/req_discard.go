package roomrequest

import (
	"protocol"
	"game/interfacer"
	"game/players"
	"game/room"
	"code.google.com/p/goprotobuf/proto"
	"lib/socket"
)

func init() {
	socket.Regist(&protocol.CDiscard{}, discard)
}

//打牌
func discard(ctos *protocol.CDiscard, c interfacer.IConn) {
	stoc := &protocol.SDiscard{}
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	if rdata == nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotInRoom))
		c.Send(stoc)
		return
	}

	var card uint32 = ctos.GetCard()
	seat := player.GetSeat()
	if card == 0 {
		stoc.Error = proto.Uint32(uint32(protocol.Error_CardValueZero))
		c.Send(stoc)
		return
	}
	err := rdata.Discard(seat, byte(card), false)
	if err > 0 {
		stoc.Error = proto.Uint32(uint32(err))
		c.Send(stoc)
	}
}
