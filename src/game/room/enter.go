package room

import (
	"game/interfacer"
	"protocol"
)

//进入
func (t *Desk) Enter(p interfacer.IPlayer) int32 {
	t.Lock() //房间加锁
	defer t.Unlock()

	for k, v := range t.players {
		if p.GetUserid() == v.GetUserid() {
			p.SetRoom(t.id, k, t.data.Code)
			round, expire := t.getRound()
			// 判断玩家是否已经在房间
			msg1 := t.res_reEnter(t.id, k, round, expire,
				t.dealer, t.data, t.players, t.maizi)
			p.Send(msg1)
			return 0
		}
	}
	p.ClearRoom()

	// 玩家不在房间，查找空座位
	for i := uint32(1); i <= 4; i++ {
		if _, ok := t.players[i]; !ok {
			t.players[i] = p
			p.SetRoom(t.id, i, t.data.Code)
			round, expire := t.getRound()
			msg1 := t.res_enter(t.id, i, round, expire,
				t.dealer, t.data, t.players, t.ready, t.maizi)
			p.Send(msg1)
			uid := p.GetUserid()
			score := t.data.Score[uid]
			msg2 := res_othercomein(p, score)
			t.broadcast_(i, msg2)
			return 0
		}
	}

	// 房间已满员
	return int32(protocol.Error_RoomFull)
}
