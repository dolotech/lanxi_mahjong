package room

import (
	"game/algorithm"
	"github.com/golang/glog"
	"protocol"
)

//出牌,加锁,托管自摸自动胡牌
func (t *Desk) Discard(seat uint32, card byte, ok bool) int32 {
	t.Lock()
	defer t.Unlock()
	if !t.state {
		glog.Infof("DiscardL err -> %d", seat)
		return int32(protocol.Error_NoStarted)
	}
	if seat < 1 || seat > 4 {
		glog.Infof("Discard seat err -> %d", seat)
		return int32(protocol.Error_NotInRoom)
	}
	if seat != t.seat || t.draw == 0 {
		return int32(protocol.Error_NotYourTurn)
	}
	t.discard_(card)
	return 0
}

//出牌,没加锁
func (t *Desk) discard_(card byte) {
	t.operateInit()  //清除操作记录
	t.discard = card //设置打牌状态
	t.draw = 0       //清除摸牌状态
	t.operate = 0    //清除操作状态
	//检测(胡,碰杠,吃)
	for s, _ := range t.players {
		if s == t.seat { //出牌人跳过
			t.out(t.seat, card) //移除牌
			continue
		}
		var cards []byte = t.getHandCards(s)
		//胡,杠碰,吃检测
		v_h := t.DiscardHu(card, cards, t.getChowCards(s), t.getPongCards(s), t.getKongCards(s), t.luckyCard,s) //胡
		if v_h > 0 && t.getSkip(s, v_h) { //是否过圈
			t.opt[s] = v_h
		}
		v_p := algorithm.DiscardPong(card, cards, t.luckyCard) //碰杠
		if v_p > 0 {
			t.opt[s] |= v_p
		}
		v_c := algorithm.DiscardChow(t.seat, s, card, cards, t.luckyCard) //吃

		if v_c > 0 {
			t.opt[s] |= v_c
		}
	}
	if t.kong { //杠操作出牌标识
		t.kong = false //杠后出牌清除
	}
	if t.optCount() == 0 { //无操作
		//出牌协议消息通知
		msg := res_discard(t.seat, t.discard)
		t.broadcast(msg) //消息广播
		 t.drawcard()  //摸牌
		return
	}
	t.turn()
}
