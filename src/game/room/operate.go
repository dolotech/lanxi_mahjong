package room

import (
	"game/algorithm"
	"protocol"
)

// 玩家进行吃/碰/杠/胡操作
func (t *Desk) Operate(card uint32, value int64, seat uint32) int32 {
	t.Lock()
	defer t.Unlock()
	if !t.state {
		return int32(protocol.Error_NoStarted)
	}

	if t.opt[seat] == 0 {
		return int32(protocol.Error_NoOperate)
	}

	if value > 0 && t.opt[seat]&value == 0 {
		return int32(protocol.Error_NoOperate)
	}

	if t.opt[seat]&algorithm.HU > 0 && (value&algorithm.HU > 0 || value == 0) {
		if value&algorithm.HU > 0 {
			// 因为没有一炮多响，所以去除其他玩家的操作状态
			for k, _ := range t.opt {
				if k != seat {
					t.opt[k] = 0
				}
			}

			var card byte = t.operateCard() //所胡的牌
			if t.opt[seat]&algorithm.QIANG_GANG > 0 {
				card = t.qiangKongCard
			}

			t.he(seat, card)
		} else {
			if t.opt[seat]&algorithm.PAOHU > 0 {
				t.skipL(seat, t.opt[seat]) //跳过胡牌
			}

			count := 0
			for _, v := range t.opt {
				if v > 0 {
					count ++
				}
			}
			if count == 1 {
				t.opt[seat] = 0
			} else {
				t.opt[seat] = t.opt[seat] & 0x7E // 掩掉胡的操作
			}

			t.cancelOperate(seat, t.seat, value, card)

			if seat != t.seat { //如果暗杠时取消不应该摸牌,应该打牌
				t.turn() //取消时进入下一个操作
			}
		}
	} else if t.opt[seat]&algorithm.AN_KONG > 0 && (value&algorithm.AN_KONG > 0 || value == 0) {
		t.opt[seat] = 0
		if value&algorithm.AN_KONG > 0 {
			if !t.kong_(card, value, seat) {
				return int32(protocol.Error_NoOperate)
			}
		} else {
			t.cancelOperate(seat, t.seat, value, card)
		}
	} else if t.opt[seat]&algorithm.MING_KONG > 0 && (value&algorithm.MING_KONG > 0 || value == 0) {
		t.opt[seat] = 0
		if value&algorithm.MING_KONG > 0 {
			if !t.kong_(card, value, seat) {
				return int32(protocol.Error_NoOperate)
			}
		} else {
			t.cancelOperate(seat, t.seat, value, card)
			t.turn() // 取消时进入下一个操作
		}
	} else if t.opt[seat]&algorithm.BU_KONG > 0 && (value&algorithm.BU_KONG > 0 || value == 0) {
		t.opt[seat] = 0
		if value&algorithm.BU_KONG > 0 {
			if !t.kong_(card, value, seat) {
				return int32(protocol.Error_NoOperate)
			}

		} else {
			t.cancelOperate(seat, t.seat, value, card)
		}
	} else if t.opt[seat]&algorithm.PENG > 0 && (value&algorithm.PENG > 0 || value == 0) {
		t.opt[seat] = 0
		if value&algorithm.PENG > 0 {
			if !t.pong_(card, value, seat) {
				return int32(protocol.Error_NoOperate)
			}
		} else {
			t.cancelOperate(seat, t.seat, value, card)
			t.turn() // 取消时进入下一个操作
		}
	} else if t.opt[seat]&algorithm.CHOW > 0 && (value&algorithm.CHOW > 0 || value == 0) {
		t.opt[seat] = 0
		if value&algorithm.CHOW > 0 {
			if !t.chow_(card, value, seat) {
				return int32(protocol.Error_NoOperate)
			}
		} else {
			t.cancelOperate(seat, t.seat, value, card)
			t.turn() //取消时进入下一个操作
		}
	} else {
		return int32(protocol.Error_NoOperate)
	}
	return 0
}

func (t *Desk) turn() {
	count := t.optCount()
	if count == 0 { //无操作
		t.drawcard() //摸牌
		return
	}
	// 计算有几个人有操作

	if count > 1 {
		// 按逆时针优先胡牌
		for i := uint32(1); i <= 4; i++ {
			s := (t.seat+ i) % 4
			if s == 0 {
				s = 4
			}
			if t.opt[s]&algorithm.HU > 0 {
				// 只给胡的操作，掩掉吃、碰和杠
				value := t.opt[s] >> 7
				value <<= 7
				t.sendOperate(t.opt[s], s)
				return
			}
		}
		// 提示碰和杠
		for i := uint32(1); i <= 4; i++ {
			if t.opt[i]&algorithm.KONG > 0 || t.opt[i]&algorithm.PENG > 0 {
				t.sendOperate(t.opt[i], i)
				return
			}
		}
		// 提示吃
		for i := uint32(1); i <= 4; i++ {
			if t.opt[i]&algorithm.CHOW > 0 {
				t.sendOperate(t.opt[i], i)
				return
			}
		}
	} else {
		// 只有个人有操作
		for k, v := range t.opt {
			if v > 0 {
				t.sendOperate(v, k)
				return
			}
		}
	}
}

func (t *Desk) sendOperate(value int64, seat uint32) {
	card := t.operateCard() //所操作的牌
	p := t.getPlayer(seat)
	if t.operate == 0 { //第一次提示操作
		t.operate += 1 //操作状态变化
		msg1 := res_discard3(t.seat, value, card)
		p.Send(msg1)
		msg2 := res_discard(t.seat, card)
		t.broadcast_(seat, msg2) //消息广播
	} else { //二次提示操作
		t.operate += 1 //操作状态变化
		msg3 := res_pengkong(t.seat, value, card)
		p.Send(msg3)
	}
}
