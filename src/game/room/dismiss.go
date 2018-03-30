package room

import (
	"time"
	"protocol"
)

//投票解散,ok=true发起,=false投票
func (t *Desk) Vote(ok bool, seat, vote uint32) int32 {
	t.Lock() //房间加锁
	defer t.Unlock()
	if t.vote != 0 && ok { //发起
		return int32(protocol.Error_RunningNotVote)
	}
	if t.vote == 0 && !ok { //投票
		return int32(protocol.Error_VotingCantLaunchVote)
	}
	t.votes[seat] = vote //投票
	if ok {              //发起投票
		t.vote = seat //发起投票者
		t.voteT = time.AfterFunc(2*time.Minute,
			func() { t.dismiss(true) }) //超时设置
		msg := res_voteStart(seat)
		t.broadcast(msg)
		//return 0
	}
	msg := res_vote(seat, vote)
	t.broadcast(msg)
	t.dismiss(false)
	return 0
}

//投票解散,agree > unagree
func (t *Desk) dismiss(ok bool) {
	var agree int = 0
	var unagree int = 0
	var i uint32
	for i = 1; i <= 4; i++ {
		if _, ok2 := t.votes[i]; !ok2 && ok {
			t.votes[i] = 0 //超时不投票默认算同意
		}
		if _, ok3 := t.players[i]; !ok3 && !ok {
			t.votes[i] = 0 //空位置提前默认同意
		}
	}
	for i = 1; i <= 4; i++ {
		if v, ok4 := t.votes[i]; ok4 && v == 0 {
			agree++
		} else {
			unagree++
		}
	}
	if agree > unagree { //一半以上通过即可
		msg := res_voteResult(0) //0解散,1不解散
		t.broadcast(msg)
		if t.voteT != nil {
			t.voteT.Stop()
		}
		round, expire := t.getRound()
		t.close(round, expire, true) //解散
	} else if ok || len(t.votes) == 4 { //结束投票
		msg := res_voteResult(1)
		t.broadcast(msg)
		t.vote = 0
		t.votes = make(map[uint32]uint32)
		if t.voteT != nil {
			t.voteT.Stop()
		}
	}
}
