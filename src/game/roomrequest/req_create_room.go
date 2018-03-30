package roomrequest

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"game/data"
	"game/interfacer"
	"game/players"
	"game/room"
	"lib/socket"
	"lib/utils"
	"protocol"
)

func init() {
	socket.Regist(&protocol.CCreatePrivateRoom{}, create)
}

// 私人局,创建房间
func create(ctos *protocol.CCreatePrivateRoom, c interfacer.IConn) {
	stoc := &protocol.SCreatePrivateRoom{}
	player := players.Get(c.GetUserid())
	rdata := room.Get(player.GetInviteCode())
	if rdata != nil {
		// 如果该玩家已经在私人局直接进入
		code := rdata.Enter(player)
		if code == 0 {
			return
		}
	}

	player.ClearRoom()
	creator := c.GetUserid()
	round := ctos.GetRound()

	exist := false
	for _, v := range config.Opts().SelectGameCountCt {
		if v == round {
			exist = true
			break
		}
	}

	if !exist {
		stoc.Error = proto.Uint32(uint32(protocol.Error_GameRoundIllegal))
		c.Send(stoc)
		return
	}

	expire := uint32(utils.Timestamp()) + round*600

	code := room.GenInvitecode(10)
	roomid, _ := data.GenRoomID()

	var cost uint32
	for k, v := range config.Opts().SelectGameCountCt {
		if v == round {
			if k < len(config.Opts().Price) {
				cost = config.Opts().Price[k]
				break
			}
		}
	}

	// 预判断房卡数量，是否足够创建房间
	if player.GetRoomCard() < cost {
		stoc.Error = proto.Uint32(uint32(protocol.Error_NotEnough_ROOM_CARD))
		c.Send(stoc)
		return
	}

	r := room.NewDeskData(uint32(roomid), round, expire, config.Opts().Ante, cost, creator, code, ctos.GetMaizi())
	r.Convert().Save()
	roomdata := &protocol.RoomData{
		Roomid:     proto.Uint32(uint32(roomid)),
		Rtype:      ctos.Rtype,
		Expire:     proto.Uint32(expire),
		Round:      proto.Uint32(round),
		Rname:      ctos.Rname,
		Invitecode: proto.String(code),
		Count:      proto.Uint32(1),
		Userid:     proto.String(creator),
		Maizi:      ctos.Maizi,
	}
	player.SetLongitudeLatitude(ctos.GetLongitude(), ctos.GetLatitude())
	rdata = room.NewDesk(r)
	room.Add(code, rdata)

	stoc.Rdata = roomdata
	c.Send(stoc)
}
