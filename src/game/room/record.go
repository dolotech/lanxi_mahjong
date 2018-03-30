package room

import (
	"game/algorithm"
	"game/data"
)

// private record
func (t *Desk) privaterecord(coins map[uint32]int32) {
	var huType uint32 = 3 //胡牌类型
	for _, v := range t.opt {
		if v&algorithm.PAOHU > 0 { //放冲
			huType = 2
		} else if v&algorithm.ZIMO > 0 { //自摸
			huType = 1
		}
	}
	//---
	ante   := t.data.Ante
	code   := t.data.Code
	rounds := t.data.Round
	zhuang := t.dealer
	round  := t.round
	id     := t.id
	expire := t.data.Expire
	cid    := t.data.Cid
	ctime  := t.data.CTime
	//---
	roundRecord := &data.GameOverRoundRecord{
		Zhuang: zhuang,
		Hutype: huType,
		Round:  round,
	}
	//---
	var userids []string
	for k, v := range t.players {
		//var userid string = v.GetUserid()
		var coin int32 = 0
		//if n, ok := t.data.Score[userid]; ok {
		if n, ok := coins[k]; ok {
			coin = n //只记录当前局输赢
		}
		var huValue uint32
		var hucard uint32
		details := &data.GameOverUserRecord{
			Userid:  v.GetUserid(),
			Seat:    k,
			Coin:    coin, //只记录当前局输赢
			Huvalue: huValue,
			HuCard:  hucard,
		}
		roundRecord.Users = append(roundRecord.Users, details)
		userids = append(userids, v.GetUserid())
	}
	if round == 1 {
		record := &data.GameOverRecord{
			RoomId:     id,
			TotalRound: rounds,
			Invitecode: code,
			Ante:       ante,
			Userids:    userids,
			Cid:        cid,
			Expire:     expire,
			Ctime:      ctime,
		}
		record.Rounds = append(record.Rounds, roundRecord)
		record.Add()
	} else {
		roundRecord.Push(id)
	}
}
