package roomrequest

import (
	"protocol"
	"game/interfacer"
	"game/players"
	"game/room"
	"github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
	"lib/socket"
)

func init() {
	socket.Regist(&protocol.COperate{}, operate)
}

//操作
func operate(ctos *protocol.COperate, c interfacer.IConn) {
	stoc := &protocol.SOperate{}
	var card uint32 = ctos.GetCard()
	var value int64 = ctos.GetValue()
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	if rdata == nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotInRoom))
		c.Send(stoc)
		return
	}
	seat := player.GetSeat()
	code := rdata.Operate(card, value, seat)
	if code > 0 {
		stoc.Error = proto.Uint32(uint32(code))
		glog.Errorln(stoc.String())
		c.Send(stoc)
	}
}
