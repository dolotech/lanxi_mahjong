package roomrequest

import (
	"protocol"
	"game/interfacer"
	"game/players"
	"game/room"
	"game/algorithm"
	"github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
	"lib/socket"
)

func init() {
	socket.Regist(&protocol.CHu{}, hu)
	socket.Regist(&protocol.CQiangKong{}, qiangkong)
}


//胡牌请求
func hu(ctos *protocol.CHu, c interfacer.IConn) {
	stoc := &protocol.SHu{}
	//TODO:优化
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	//
	if rdata != nil {
		seat := player.GetSeat()
		rdata.Operate(0, algorithm.HU, seat)
	} else {
		glog.Infof("hu room -> %s", c.GetUserid())
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotInRoom))
		c.Send(stoc)
	}
}

//抢杠胡牌请求
func qiangkong(ctos *protocol.CQiangKong, c interfacer.IConn) {
	stoc := &protocol.SQiangKong{}
	//TODO:优化
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	//
	if rdata != nil {
		seat := player.GetSeat()
		rdata.Operate(0, algorithm.HU, seat)
	} else {
		glog.Infof("qiangkong room -> %s", c.GetUserid())
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotInRoom))
		c.Send(stoc)
	}
}
