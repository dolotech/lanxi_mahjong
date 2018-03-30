package room

import (
	"protocol"
	"time"
)

//玩家准备
func (t *Desk) Readying(seat uint32, ready bool) int32 {
	t.Lock() //房间加锁
	defer t.Unlock()
	if t.vote != 0 { //投票中不能准备
		return int32(protocol.Error_VotingCantLaunchVote)
	}
	t.ready[seat] = ready //设置状态
	msg := res_ready(seat, ready)
	t.broadcast(msg) //广播消息

	if !t.data.MaiZi {
		go func() {
			<-time.After(time.Millisecond * 200) //延迟200豪秒
			t.Diceing()                          //主动打骰
		}()
	}
	return 0
}

//玩家打骰子切牌,发牌
func (t *Desk) Diceing() bool {
	t.Lock() //房间加锁
	defer t.Unlock()
	if t.isDiceing() { //是否打骰

		//t.dealer_()   //打庄
		t.gameStart() //开始牌局
		return true
	} else {
		return false
	}
}

//是否可以打骰
func (t *Desk) isDiceing() bool {
	if t.state { //已经开始
		return false
	}

	if t.dealer != 0 { //已经打庄
		return false
	}

	if len(t.players) != 4 { //人数不够
		return false
	}

	if len(t.ready) < 4 {
		return false
	}

	if t.data.MaiZi {
		if len(t.maizi) < 4 {
			return false
		}
	}
	return true
}
