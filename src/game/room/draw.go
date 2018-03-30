package room

import (
	"game/algorithm"
	"time"
)

//模牌,kong==false普通摸牌,kong==true扛后摸牌
func (t *Desk) drawcard() {
	<-time.After(time.Millisecond * 100)
	if len(t.cards) == 0 {
		t.he(0, 0) //结束牌局
		return
	}

	var value uint32 = 0
	//var seat uint32 = t.seat //过圈
	if t.kong { //杠后摸牌
		value = 1
	} else { //普通摸牌, t.discrad != 0
		var outs []byte = t.getOutCards(t.seat)
		outs = append(outs, t.discard)
		t.outCards[t.seat] = outs
		//位置切换
		t.seat = algorithm.NextSeat(t.seat)
	}
	//t.skip_(seat, t.seat) //过圈
	t.unskipL(t.seat) //摸牌一定过圈
	t.operateInit()   //清除操作记录
	var card byte = t.cards[0]
	t.draw = card //设置摸牌状态
	t.discard = 0 //清除打牌状态
	t.operate = 0 //清除操作状态
	t.cards = t.cards[1:]
	cards := t.in(t.seat, card)

	v := t.DrawDetect(card, cards, t.getChowCards(t.seat), t.getPongCards(t.seat), t.getKongCards(t.seat), t.luckyCard, t.seat)
	v |= algorithm.DetectKong(cards, t.getPongCards(t.seat), t.luckyCard)

	if v&algorithm.HU > 0 {
		if t.kong {
			v = v | algorithm.HU_KONG_FLOWER
		}
	}
	if v > 0 { //摸牌全部记录为胡(只自己操作)
		t.opt[t.seat] = v
	}

	//其他玩家消息
	msg1 := res_otherDraw(t.seat, value, uint32(len(t.cards)))
	//摸牌协议消息
	msg2 := res_draw(value, v, card, uint32(len(t.cards)), cards)
	//摸牌协议消息通知
	for s, o := range t.players {
		if s == t.seat {
			//摸牌玩家消息
			o.Send(msg2)
		} else {
			//其他玩家消息
			o.Send(msg1)
		}
	}
}
