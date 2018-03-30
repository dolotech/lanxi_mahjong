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

	socket.Regist(&protocol.CLaunchVote{}, launchVote)

	socket.Regist(&protocol.CVote{}, vote)
}

// 发起房间解散投票
func launchVote(ctos *protocol.CLaunchVote, c interfacer.IConn) {
	//TODO:优化
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	//
	seat := player.GetSeat()
	stoc := &protocol.SLaunchVote{Seat: proto.Uint32(seat)}
	if rdata == nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotInPrivateRoom))
		c.Send(stoc)
		return
	}

	err := rdata.Vote(true, seat, 0)
	if err > 0 {
		stoc.Error = proto.Uint32(uint32(err))
		c.Send(stoc)
	}
}

// 玩家进行房间解散投票
func vote(ctos *protocol.CVote, c interfacer.IConn) {
	//TODO:优化
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	//
	seat := player.GetSeat()
	var vote uint32 = ctos.GetVote() //0同意,1不同意
	stoc := &protocol.SVote{
		Vote: proto.Uint32(vote),
		Seat: proto.Uint32(seat),
	}
	if rdata == nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotInPrivateRoom))
		c.Send(stoc)
		return
	}

	err := rdata.Vote(false, seat, vote)
	if err == 1 {
		stoc.Error = proto.Uint32(uint32(protocol.Error_RunningNotVote))
		c.Send(stoc)
	} else if err == 2 {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotVoteTime))
		c.Send(stoc)
	}
}
