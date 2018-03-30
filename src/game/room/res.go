package room

import (
	"game/algorithm"
	"game/interfacer"
	"protocol"
	"code.google.com/p/goprotobuf/proto"
	"lib/utils"
)

//操作提示响应消息
func res_operate(seat, beseat uint32, value int64, card uint32) interfacer.IProto {
	stoc := &protocol.SOperate{}
	stoc.Seat = proto.Uint32(seat)
	stoc.Card = proto.Uint32(card)
	stoc.Value = proto.Int64(value)
	stoc.Beseat = proto.Uint32(beseat)
	return stoc
}

//操作提示响应消息
func res_operate2(seat, beseat uint32, value int64, qiang int64, card uint32) interfacer.IProto {
	stoc := &protocol.SOperate{}
	stoc.Seat = proto.Uint32(seat)
	stoc.Card = proto.Uint32(card)
	stoc.Value = proto.Int64(value)
	stoc.Beseat = proto.Uint32(beseat)
	stoc.Discontinue = proto.Int64(qiang)
	return stoc
}

//打牌响应消息
func res_discard(seat uint32, card byte) interfacer.IProto {
	stoc := &protocol.SDiscard{}
	stoc.Card = proto.Uint32(uint32(card))
	stoc.Seat = proto.Uint32(seat)
	return stoc
}

func res_discard3(seat uint32, v int64, card byte) interfacer.IProto {
	stoc := &protocol.SDiscard{}
	stoc.Card = proto.Uint32(uint32(card))
	stoc.Seat = proto.Uint32(seat)
	stoc.Value = proto.Int64(v)
	return stoc
}

//处理前面有玩家胡牌优先操作,如果该玩家跳过胡牌,此协议向有碰和明杠的玩家主动发送
func res_pengkong(seat uint32, v int64, card byte) interfacer.IProto {
	stoc := &protocol.SPengKong{}
	stoc.Card = proto.Uint32(uint32(card))
	stoc.Seat = proto.Uint32(seat)
	stoc.Value = proto.Int64(v)
	return stoc
}

//进入房间响应消息
func res_othercomein(p interfacer.IPlayer, score int32) interfacer.IProto {
	userinfo := p.ConverProtoUser()
	userinfo.Score = proto.Int32(score)
	userinfo.Offline = proto.Bool(false)
	stoc := &protocol.SOtherComein{Userinfo: userinfo}
	return stoc
}

//进入房间响应消息
func (t *Desk) res_enter(id, seat, round, expire, dealer uint32,
	data *DeskData, m map[uint32]interfacer.IPlayer,
	ready map[uint32]bool, maizi map[uint32]uint32) interfacer.IProto {
	var len uint32 = uint32(len(m))
	stoc := &protocol.SEnterSocialRoom{}
	roomdata := &protocol.RoomData{
		Roomid:     proto.Uint32(id),
		Rtype:      proto.Uint32(0),
		Rname:      proto.String(""),
		Expire:     proto.Uint32(uint32(expire)),
		Round:      proto.Uint32(round),
		Count:      proto.Uint32(len),
		Invitecode: proto.String(data.Code),
		Zhuang:     proto.Uint32(dealer),
		Userid:     proto.String(data.Cid),
		Maizi:      proto.Bool(data.MaiZi),
		Totalround: proto.Uint32(data.Round),
		Lian:       proto.Uint32(t.lianCount),
	}
	stoc.Room = roomdata
	var score int32
	for k, p := range m {
		uid := p.GetUserid()
		score1 := data.Score[uid]
		if k != seat {

			userinfo := p.ConverProtoUser()
			userinfo.Score = proto.Int32(score1)

			if zi, ok := maizi[k]; ok {
				userinfo.Maizi = proto.Int32(int32(zi))
			} else {
				userinfo.Maizi = proto.Int32(-1)
			}
			var r bool = t.getReady(k)
			userinfo.Ready = proto.Bool(r)
			ip := utils.InetTontoa(p.GetConn().GetIPAddr()).String()
			userinfo.Ip = &ip
			userinfo.Offline = proto.Bool(t.offline[k])
			stoc.Userinfo = append(stoc.Userinfo, userinfo)
		} else {
			score = score1
		}
	}

	stoc.Position = proto.Uint32(seat)
	stoc.Score = proto.Int32(score)
	stoc.Ready = proto.Bool(false)
	if zi, ok := maizi[seat]; ok {
		stoc.Mazi = proto.Int32(int32(zi))
	} else {
		stoc.Mazi = proto.Int32(-1)
	}
	stoc.Beginning = proto.Bool(false)
	stoc.Ma = proto.Uint32(0)
	stoc.LaunchSeat = proto.Uint32(t.vote)
	stoc.ConferenceId = proto.String("")
	if t.vote > 0 { //投票中
		for k, v := range t.votes {
			if v == 0 {
				stoc.VoteAgree = append(stoc.VoteAgree, k)
			} else {
				stoc.VoteDisagree = append(stoc.VoteDisagree, k)
			}
		}
	}
	stoc.LuckyCard = proto.Uint32(uint32(t.luckyCard))
	return stoc
}

func (t *Desk) orderValue(seat uint32) int64 {
	count := 0
	for _, v := range t.opt {
		if v&algorithm.HU > 0 {
			count ++
		}
	}

	// 是否有人胡牌
	if count > 0 {
		for i := uint32(1); i <= 4; i++ {
			s := (t.seat+ i) % 4
			if s == 0 {
				s = 4
			}

			if t.opt[s]&algorithm.HU > 0 {
				if s != seat {
					return 0
				}
				// 只给胡的操作，掩掉吃、碰和杠

				value := t.opt[s] >> 7
				value <<= 7
				return value
			}
		}
	}

	// 提示碰和杠
	for i := uint32(1); i <= 4; i++ {
		if t.opt[i]&algorithm.KONG > 0 || t.opt[i]&algorithm.PENG > 0 {
			if i != seat {
				return 0
			}
			return t.opt[i]
		}
	}
	return 0
}

//重新进入房间响应消息
func (t *Desk) res_reEnter(id, seat, round, expire, dealer uint32,
	data *DeskData, m map[uint32]interfacer.IPlayer, maizi map[uint32]uint32) interfacer.IProto {
	var lenght uint32 = uint32(len(m))
	stoc := &protocol.SEnterSocialRoom{}
	roomdata := &protocol.RoomData{
		Roomid:     proto.Uint32(id),
		Rtype:      proto.Uint32(0),
		Rname:      proto.String(""),
		Expire:     proto.Uint32(uint32(expire)),
		Round:      proto.Uint32(round),
		Count:      proto.Uint32(lenght),
		Invitecode: proto.String(data.Code),
		Zhuang:     proto.Uint32(dealer),
		Userid:     proto.String(data.Cid),
		Maizi:      proto.Bool(data.MaiZi),
		Totalround: proto.Uint32(data.Round),
		Lian:       proto.Uint32(t.lianCount),
	}
	stoc.Room = roomdata
	var score int32
	for k, p := range m {
		uid := p.GetUserid()
		score1 := data.Score[uid]
		if k != seat {
			var r bool = t.getReady(k)
			userinfo := p.ConverProtoUser()
			userinfo.Ready = proto.Bool(r)
			userinfo.Score = proto.Int32(score1)
			stoc.Userinfo = append(stoc.Userinfo, userinfo)
			if zi, ok := maizi[k]; ok {
				userinfo.Maizi = proto.Int32(int32(zi))
			} else {
				userinfo.Maizi = proto.Int32(-1)
			}

			ip := utils.InetTontoa(p.GetConn().GetIPAddr()).String()
			userinfo.Ip = &ip
			userinfo.Offline = proto.Bool(t.offline[k])
		} else {
			score = score1
		}
	}
	stoc.Position = proto.Uint32(seat)
	var r bool = t.getReady(seat)
	stoc.Ready = proto.Bool(r)

	if zi, ok := maizi[seat]; ok {
		stoc.Mazi = proto.Int32(int32(zi))
	} else {
		stoc.Mazi = proto.Int32(-1)
	}
	stoc.Beginning = proto.Bool(t.state)
	stoc.Score = proto.Int32(score)
	stoc.Ma = proto.Uint32(0)
	if zi, ok := maizi[seat]; ok {
		stoc.Mazi = proto.Int32(int32(zi))
	} else {
		stoc.Mazi = proto.Int32(-1)
	}
	stoc.LaunchSeat = proto.Uint32(t.vote)
	stoc.ConferenceId = proto.String("")

	if t.vote > 0 { //投票中
		for k, v := range t.votes {
			if v == 0 {
				stoc.VoteAgree = append(stoc.VoteAgree, k)
			} else {
				stoc.VoteDisagree = append(stoc.VoteDisagree, k)
			}
		}
	}
	//重连数据 TODO:优化
	if !t.state {
		stoc.Beginning = proto.Bool(false)
		return stoc
	}
	stoc.Beginning = proto.Bool(true)
	stoc.Turn = proto.Uint32(t.seat)
	stoc.Dice = proto.Uint32(t.dice)
	stoc.CardsCount = proto.Uint32(uint32(len(t.cards)))
	stoc.Handcards = t.getHandCards(seat)
	stoc.LuckyCard = proto.Uint32(uint32(t.luckyCard))
	value := t.opt[seat]

	count := 0
	for _, v := range t.opt {
		if v > 0 {
			count ++
		}
	}
	if count > 1 {
		value = t.orderValue(seat)
	}

	stoc.Value = proto.Int64(value) //操作值
	var kongnum uint32
	for i, _ := range m {
		pcard := &protocol.ProtoCard{}
		pcard.Seat = proto.Uint32(i)
		pongs := t.getPongCards(i) //碰牌数据
		for _, v := range pongs {
			j, card := algorithm.DecodePeng(v)
			var peng uint32 = j << 24
			peng |= (uint32(card) << 16)
			pcard.Peng = append(pcard.Peng, peng)
		}
		kongs := t.getKongCards(i) //杠牌数据
		for _, v := range kongs {
			j, card, classify := algorithm.DecodeKong(v)
			var kong uint32 = j << 24
			kong |= (uint32(card) << 16)
			kong |= (classify << 8)
			pcard.Kong = append(pcard.Kong, kong)
			kongnum = kongnum + 1
		}
		chows := t.getChowCards(i) //吃牌数据
		pcard.Chow = chows
		pcard.Outcards = t.getOutCards(i)

		//加上进行中的牌
		if t.seat == i && t.discard != 0 {
			pcard.Outcards = append(pcard.Outcards, t.discard)
		}
		stoc.Cards = append(stoc.Cards, pcard)
	}
	// 桌上有几个杠(未摸起的牌，牌尾被摸走几张牌)
	stoc.KongCount = proto.Uint32(kongnum)
	return stoc
}

//打庄响应消息
func res_dealer(dealer uint32) interfacer.IProto {
	stoc := &protocol.SZhuang{}
	stoc.Zhuang = proto.Uint32(dealer)
	//stoc.Lian = proto.Uint32(t.lian)
	return stoc
}

//庄家响应消息
func res_zhuangDeal(v int64, dice uint32, cards []byte, luckyCard uint32) interfacer.IProto {
	stoc := &protocol.SZhuangDeal{}
	stoc.Value = proto.Int64(v)
	stoc.Dice = proto.Uint32(dice)
	stoc.Cards = cards
	stoc.LuckyCard = proto.Uint32(luckyCard)
	return stoc
}

//闲家响应消息
func res_deal(v int64, dice uint32, cards []byte, luckyCard uint32) interfacer.IProto {
	stoc := &protocol.SDeal{}
	stoc.Value = proto.Int64(v)
	stoc.Dice = proto.Uint32(dice)
	stoc.Cards = cards
	stoc.LuckyCard = proto.Uint32(luckyCard)
	return stoc
}

//闲家响应消息
func res_otherDraw(seat, value, remainder uint32) interfacer.IProto {
	stoc := &protocol.SOtherDraw{}
	stoc.Seat = proto.Uint32(seat)
	stoc.Kong = proto.Uint32(value)
	stoc.Remainder = proto.Uint32(remainder)
	return stoc
}

//摸牌协议消息响应消息
func res_draw(kong uint32, v int64, card byte, remainder uint32, cards []byte) interfacer.IProto {
	stoc := &protocol.SDraw{}
	stoc.Card = proto.Uint32(uint32(card))
	stoc.Cards = cards
	stoc.Kong = proto.Uint32(kong)
	stoc.Value = proto.Int64(v)
	stoc.Remainder = proto.Uint32(remainder)
	return stoc
}

//玩家准备消息响应消息
func res_ready(seat uint32, ready bool) interfacer.IProto {
	stoc := &protocol.SReady{}
	stoc.Ready = proto.Bool(ready)
	stoc.Seat = proto.Uint32(seat)
	return stoc
}

//结束牌局响应消息
func res_he(seat uint32, card byte) interfacer.IProto {
	stoc := &protocol.SHu{}
	stoc.Seat = proto.Uint32(seat)
	stoc.Card = proto.Uint32(uint32(card))
	return stoc
}

//结束牌局响应消息,huType:0:黄庄，1:自摸，2:炮胡
func res_over(seat, lian uint32,
	handCards map[uint32][]byte, huValue map[uint32]int64,
	huFan, total, coin map[uint32]int32, maizi map[uint32]uint32) interfacer.IProto {
	stoc := &protocol.SGameover{
		Data: make([]*protocol.ProtoCount, 4),
	}
	var huType uint32 = 0  //胡牌类型
	var paoSeat uint32 = 0 //放冲玩家
	var i uint32
	for i = 1; i <= 4; i++ {
		val := huValue[i] //胡牌掩码
		if val&algorithm.PAOHU > 0 { //放冲
			huType = 2
			paoSeat = seat
		} else if val&algorithm.ZIMO > 0 { //自摸
			huType = 1
		}
		val >>= 8 //显示处理
		val <<= 8 //显示处理

		protoCount := &protocol.ProtoCount{
			Seat:  proto.Uint32(i),
			Maizi: proto.Int32(int32(maizi[i])),
			Hu:    proto.Int64(val),
			// 结算时要显示该玩家的手牌
			Cards:     handCards[i],
			Total:     proto.Int32(total[i]),
			Coin:      proto.Int32(coin[i]),
			HuTypeFan: proto.Int32(huFan[i]), //前端只显示牌型分,所以给的一样
			HuFan:     proto.Int32(huFan[i]),
		}
		stoc.Data[i-1] = protoCount
	}
	stoc.LianCount = proto.Uint32(lian)
	stoc.PaoSeat = proto.Uint32(paoSeat)
	stoc.HuType = proto.Uint32(huType)
	return stoc
}

//离开房间响应消息
func res_leave(seat uint32) interfacer.IProto {
	stoc := &protocol.SPrivateLeave{}
	stoc.Seat = proto.Uint32(seat)
	return stoc
}

//私人局结束响应消息
func res_privateOver(id, round, expire uint32,
	m map[uint32]interfacer.IPlayer, n map[string]int32) interfacer.IProto {
	stoc := &protocol.SPrivateOver{}
	// 如果是私人房，且房间过期 ，踢掉房间玩家
	stoc.Cid = proto.Uint32(0)
	stoc.Roomid = proto.Uint32(id)
	stoc.Round = proto.Uint32(round)
	stoc.Expire = proto.Uint32(expire)
	for _, p := range m {
		userid := p.GetUserid()
		s := &protocol.PrivateScore{
			Userid: proto.String(userid),
			Score:  proto.Int32(n[userid]),
		}
		stoc.List = append(stoc.List, s)
	}
	return stoc
}

//发起投票申请解散房间
func res_voteStart(seat uint32) interfacer.IProto {
	stoc := &protocol.SLaunchVote{Seat: proto.Uint32(seat)}
	return stoc
}

//投票解散房间事件结果
func res_voteResult(vote uint32) interfacer.IProto {
	stoc := &protocol.SVoteResult{Vote: proto.Uint32(vote)}
	return stoc
}

//投票
func res_vote(seat, vote uint32) interfacer.IProto {
	stoc := &protocol.SVote{
		Vote: proto.Uint32(vote),
		Seat: proto.Uint32(seat),
	}
	return stoc
}
