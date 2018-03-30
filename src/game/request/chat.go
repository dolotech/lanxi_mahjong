package request

import (
	"lib/socket"
	"game/interfacer"
	"game/players"
	"protocol"
	"game/room"
	"code.google.com/p/goprotobuf/proto"
)

func init() {
	socket.Regist(&protocol.CBroadcastChatText{}, chattext)
	socket.Regist(&protocol.CBroadcastChat{}, chatsound)
}

// 文本聊天
func chattext(ctos *protocol.CBroadcastChatText, c interfacer.IConn) {
	stoc := &protocol.SBroadcastChatText{}
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	if rdata != nil {
		seat := player.GetSeat()
		stoc.Seat = &seat
		stoc.Content = ctos.Content
		rdata.Broadcasts(stoc)
	} else {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotInRoom))
	}
	if stoc.Error != nil {
		c.Send(stoc)
	}
}

// 语音聊天
func chatsound(ctos *protocol.CBroadcastChat, c interfacer.IConn) {
	stoc := &protocol.SBroadcastChat{}
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	if rdata != nil {
		seat := player.GetSeat()
		stoc.Seat = &seat
		stoc.Content = ctos.Content
		rdata.Broadcasts(stoc)
	} else {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotInRoom))
	}
	if stoc.Error != nil {
		c.Send(stoc)
	}
}
