package room

import (
	"game/algorithm"
	"lib/utils"
)

//抢杠胡处理
func (t *Desk) qiangKongHe(seat uint32, card byte) {
	if t.opt[seat]&algorithm.QIANG_GANG > 0 {
		msg := res_operate2(seat, t.seat, algorithm.QIANG_GANG,
			algorithm.QIANG_GANG, uint32(card))
		t.broadcast(msg)
		kongs := t.getKongCards(t.seat)
		var cs []uint32
		for i, v2 := range kongs { //杠
			_, c, mask := algorithm.DecodeKong(v2) //解码
			if c == card && int64(mask) == algorithm.BU_KONG {
				cs = append(kongs[:i], kongs[i+1:]...)
				t.kongCards[t.seat] = cs
				break
			}
		}
	}

}

// 海底捞,摸到属于每个玩家的最后一张牌
func (t *Desk) haidilaoHe() int64 {
	if len(t.cards) < 4 {
		return algorithm.HU_HAI_LAO
	}
	return 0
}

//天胡
func (t *Desk) tianHe(seat uint32) int64 {
	if uint32(len(t.cards)) != algorithm.TOTAL-(algorithm.HAND*4+1) {
		return 0
	}

	if seat != t.dealer {
		return 0
	}

	var l_o int = len(t.outCards)
	var l_p int = len(t.pongCards)
	var l_k int = len(t.kongCards)
	var l_c int = len(t.chowCards)
	if l_o == 0 && l_p == 0 && l_k == 0 && l_c == 0 {
		if t.discard == 0 {
			return algorithm.TIAN_HU
		}
	}
	return 0
}

//地胡
func (t *Desk) diHe(seat uint32) int64 {
	if uint32(len(t.cards)) < algorithm.TOTAL-(algorithm.HAND*4+4) {
		return 0
	}
	if t.dealer == seat {
		return 0
	}

	if len(t.outCards[seat]) > 0 {
		return 0
	}

	if len(t.pongCards[seat]) > 0 {
		return 0
	}
	if len(t.chowCards[seat]) > 0 {
		return 0
	}
	if len(t.kongCards[seat]) > 0 {
		return 0
	}

	for _, v := range t.pongCards {
		for _, c := range v {
			s, _ := algorithm.DecodePeng(c)
			if s == seat {
				return 0
			}
		}
	}

	for _, v := range t.kongCards {
		for _, c := range v {
			s, _, _ := algorithm.DecodeKong(c)
			if s == seat {
				return 0
			}
		}
	}

	tseat := seat + 1
	if tseat > 4 {
		tseat = 1
	}
	for s, v := range t.chowCards {
		if tseat == s {
			if len(v) > 0 {
				return 0
			}
			break
		}
	}
	return algorithm.DI_HU
}

// seat 胡牌位置玩家的座位  0表示黄庄
func (t *Desk) he(seat uint32, card byte) {
	if !t.state {
		return
	}
	t.qiangKongHe(seat, t.qiangKongCard) //抢杠胡处理
	huangZhuang := seat == 0             //是否黄庄

	// 是否包三家
	var bao bool
	if !huangZhuang && len(t.cards) <= 4 {
		if t.opt[seat]&algorithm.PAOHU > 0 {
			bao = true
		}
	}

	//算番
	total := t.gameOver(huangZhuang)

	// 牌墙只剩4张牌，放炮玩家包三家
	if bao {
		var gap int32
		for k, _ := range total {
			if k != t.seat && k != seat {
				gap += total[k]
				total[k] = 0
			}
		}
		total[t.seat] += gap
	}

	score := make(map[uint32]int32)

	for i, v := range total {
		p := t.getPlayer(i)
		uid := p.GetUserid()
		t.data.Score[uid] += v //总分
		score[i] = v           //当局分
	}
	if !huangZhuang { //胡牌
		//胡牌消息
		msg1 := res_he(seat, card)
		t.broadcast(msg1)
	}
	//结算消息广播
	msg2 := res_over(t.seat, t.lianCount, t.handCards, t.opt, total, total, score, t.maizi)
	t.broadcast(msg2)
	t.round++ //局数
	round, expire := t.getRound()
	msg := res_privateOver(t.id, round, expire, t.players, t.data.Score)
	t.broadcast(msg)                //私人局结束消息广播
	t.privaterecord(score)          //日志记录 TODO:goroutine
	t.lianDealer(huangZhuang, seat) //连庄
	t.overSet()                     //重置状态
	t.close(round, expire, false)   //结束牌局
}

//连庄计算
func (t *Desk) lianDealer(huangZhuang bool, seat uint32) {
	if !huangZhuang { //胡牌
		// 稳庄
		if t.dealer == seat {
			t.lianCount ++
		} else {
			t.dealer = seat //胡牌玩家接庄
			t.lianCount = 0
		}
	} else { //黄庄稳庄

	}
}

//结束牌局重置状态数据
func (t *Desk) overSet() {
	t.state = false //牌局状态
	t.dealer = 0    //庄家重置
	t.discard = 0   //重置打牌
	t.draw = 0      //重置摸牌
	t.dice = 0      //重置骰子
	t.seat = 0      //清除位置
	t.kong = false  //清除杠牌
	t.ready = make(map[uint32]bool, 4)
	t.maizi = make(map[uint32]uint32, 4)
}

//结束牌局,ok=true投票解散
func (t *Desk) close(round, expire uint32, ok bool) {
	var n uint32 = uint32(utils.Timestamp())
	if (round > 0 && expire > n) && !ok {
		return
	}
	if t.closeCh != nil {
		close(t.closeCh) //关闭计时器
		t.closeCh = nil  //消除计时器
	}
	for k, p := range t.players {
		// 投票解散不广播退出房间协议
		if !ok{
			msg := res_leave(k)
			t.broadcast(msg)
		}
		p.ClearRoom() //清除玩家房间数据

	}
	Del(t.data.Code) //从房间列表中清除
}
