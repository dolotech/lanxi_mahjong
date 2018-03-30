package room

import (
	"code.google.com/p/goprotobuf/proto"
	"config"
	"fmt"
	"game/algorithm"
	"game/interfacer"
	"github.com/golang/glog"
	"lib/utils"
	"protocol"
)

//新建一张牌桌
func NewDesk(data *DeskData) interfacer.IDesk {
	return &Desk{
		id:      data.Rid,
		data:    data,
		players: make(map[uint32]interfacer.IPlayer, 4),
		votes:   make(map[uint32]uint32, 4),
		ready:   make(map[uint32]bool, 4),
		maizi:   make(map[uint32]uint32, 4),
		opt:     make(map[uint32]int64, 4),
		offline: make(map[uint32]bool, 4),
	}
}

func (t *Desk) SetCheat(handcards [][]byte, leftcards []byte, wildcard byte) {
	t.Lock() //房间加锁
	defer t.Unlock()
	t.cheatHandCards = handcards
	t.cheatLeftCards = leftcards
	t.cheatwildcard = wildcard
}

//房间消息广播,聊天
func (t *Desk) Broadcasts(msg interfacer.IProto) {
	t.Lock() //房间加锁
	defer t.Unlock()
	t.broadcast(msg)
}

//关闭房间,停服or过期清除等等.TODO:玩牌中是否关闭?
func (t *Desk) Closed(ok bool) {
	round, expire := t.getRound()
	t.close(round, expire, ok) //ok=true强制解散,=false清理
}

//初始化操作值
func (t *Desk) operateInit() {
	t.qiangKongCard = 0
	for k, _ := range t.opt {
		t.opt[k] = 0
	}
}

//1.胡的牌(自摸,放炮),2.抢杠胡时操作的牌(补杠人的摸牌),TODO:优化
func (t *Desk) operateCard() byte {
	if t.discard != 0 {
		return t.discard
	} else {
		return t.draw
	}
}

//胡牌过圈
func (t *Desk) skipL(seat uint32, huValue int64) {
	t.skip[seat] = algorithm.HuType(huValue, config.Opts().Ante) //跳过胡牌过圈设置
}

//清除过圈(摸牌,抢杠)
func (t *Desk) unskipL(seat uint32) {
	if _, ok := t.skip[seat]; ok {
		delete(t.skip, seat)
	}
}

//胡牌过圈,1.杠后出牌,2.庄家出牌,3.正常过圈
func (t *Desk) skip_(s1, s2 uint32) {
	if len(t.skip) == 0 { //empty
		return
	}
	if s2 == t.dealer || t.kong { //1,2
		t.skip = make(map[uint32]uint32)
		return
	}
	var s []uint32 = make([]uint32, 0)
	for k, _ := range t.skip { //3
		if skiped(s1, s2, k) {
			s = append(s, k)
		}
	}
	for _, v := range s { //过圈位置
		delete(t.skip, v)
	}
}

//s1上一个位置,s2当前位置,s3胡牌位置
func skiped(s1, s2, s3 uint32) bool {
	if s1 == s2 { //s3不在s1,s2之间
		return false
	}
	var s4 = algorithm.NextSeat(s1)
	if s4 == s3 { //s3在s1,s2之间
		return true
	}
	return skiped(s4, s2, s3)
}

func (t *Desk) Offline(seat uint32, status bool) {
	if seat >= 1 || seat <= 4 {
		offline := t.offline[seat]

		if offline != status {
			t.offline[seat] = status
			stoc := &protocol.SOffline{Seat: proto.Uint32(seat), Status: proto.Bool(status)}
			t.broadcast(stoc)
		}
	}
}

//碰操作,已经验证通过
func (t *Desk) pong_(card uint32, value int64, seat uint32) bool {
	cards, boolean := t.ponging(seat, byte(card))
	if !boolean {
		return false
	}
	t.operateInit() //清除操作记录,操作成功后消除,防止重复提示
	v := algorithm.DetectKong(t.getHandCards(seat), t.getPongCards(seat), t.luckyCard)
	if v > 0 { //摸牌全部记录为胡(只自己操作)
		t.opt[seat] = v
	}
	//碰操作协议消息通知
	msg := res_operate(seat, t.seat, value, card)
	t.broadcast_(seat, msg)
	//glog.Error("碰后杠牌 %b",v)
	msg = res_operate2(seat, t.seat, value, v, card)
	t.getPlayer(seat).Send(msg)
	t.skip_(t.seat, seat) //过圈
	//状态设置
	t.seat = seat                //位置切换
	t.draw = cards[len(cards)-1] //设置摸牌,超时出牌时打出
	t.discard = 0                //重置出牌
	return true
	//等待出牌
}

//杠操作,已经验证通过
func (t *Desk) kong_(card1 uint32, value int64, seat uint32) bool {
	var card byte = byte(card1)
	switch value {
	case algorithm.BU_KONG:

		boolean := t.buKong(seat, card)

		if !boolean {
			return false
		}
	case algorithm.MING_KONG:

		boolean := t.mingKong(seat, card)
		if !boolean {
			return false
		}
	case algorithm.AN_KONG:

		boolean := t.anKong(seat, card)
		if !boolean {
			return false
		}
	}

	// 杠除8张
	if len(t.cards) > 8 {
		t.cards = t.cards[8:]
	} else {
		t.cards = t.cards[:0]
	}

	//杠操作协议消息通知
	msg := res_operate(seat, t.seat, value, card1)
	t.broadcast(msg)
	t.skip_(t.seat, seat) //过圈
	//状态设置
	t.kong = true //杠操作出牌标识
	t.seat = seat //位置切换
	if t.qiangKong(card, value) {
		t.operate++
		t.turn() //抢杠操作
	} else {
		t.drawcard() //摸牌
	}
	return true
}

//吃操作,已经验证通过
func (t *Desk) chow_(card uint32, value int64, seat uint32) bool {
	c1, c2 := algorithm.DecodeChow2(card)
	if !t.chowing(seat, c1, c2) {
		glog.Errorf("chow card error -> %x %x", c1, c2)
		return false
	}

	t.operateInit() //清除操作记录
	//吃操作协议消息通知
	card2 := algorithm.EncodeChow(c1, c2, t.discard)

	v := algorithm.DetectKong(t.getHandCards(seat), t.getPongCards(seat), t.luckyCard)
	if v > 0 { //摸牌全部记录为胡(只自己操作)
		t.opt[seat] = v
	}
	//glog.Error("吃后杠牌 %b", v)

	msg := res_operate(seat, t.seat, value, card2)
	t.broadcast_(t.seat, msg)
	msg = res_operate2(seat, t.seat, value, v, card2)

	t.getPlayer(t.seat).Send(msg)
	t.skip_(t.seat, seat) //过圈
	//状态设置
	t.seat = seat //位置切换
	var cards []byte = t.handCards[seat]
	t.draw = cards[len(cards)-1] //设置摸牌,超时出牌时打出
	t.discard = 0                //重置出牌

	return true
	//等待出牌
}

//补杠被抢杠,抢杠处理,抢杠胡牌
func (t *Desk) qiangKong(card byte, mask int64) bool {
	t.operateInit() //清除操作记录
	if mask == algorithm.AN_KONG {
		return false
	}
	//检测(抢杠胡)
	for s, _ := range t.players {
		if s == t.seat { //出牌人跳过
			continue
		}
		//抢杠不用过圈
		//var ok bool = t.getSkip(s) //是否过圈
		//if !ok {
		//var cards []byte = t.getHandCards(s)
		//胡,杠碰,吃检测
		v_h := t.DiscardHu(card, t.getHandCards(s), t.getChowCards(s), t.getPongCards(s), t.getKongCards(s), t.luckyCard, s) //胡
		if v_h > 0 {
			t.unskipL(s) //抢杠玩家一定过圈
			t.opt[s] = v_h | algorithm.QIANG_GANG
			t.qiangKongCard = card
		}
	}
	return len(t.opt) > 0
}

//取消操作时消息通知
func (t *Desk) cancelOperate(seat, beseat uint32, value int64, card uint32) {
	msg := res_operate(seat, beseat, value, card)
	p := t.getPlayer(seat)
	p.Send(msg)
}

//获取玩家
func (t *Desk) getPlayer(seat uint32) interfacer.IPlayer {
	if v, ok := t.players[seat]; ok && v != nil {
		return v
	}
	panic(fmt.Sprintf("getPlayer error:%d", seat))
}

//获取手牌
func (t *Desk) getHandCards(seat uint32) []byte {
	if v, ok := t.handCards[seat]; ok && v != nil {
		return v
	}
	panic(fmt.Sprintf("getHandCards error:%d", seat))
}

//获取海底牌
func (t *Desk) getOutCards(seat uint32) []byte {
	if v, ok := t.outCards[seat]; ok && v != nil {
		return v
	}
	return []byte{}
}

//获取碰牌
func (t *Desk) getPongCards(seat uint32) []uint32 {
	if v, ok := t.pongCards[seat]; ok && v != nil {
		return v
	}
	return []uint32{}
}

//获取杠牌
func (t *Desk) getKongCards(seat uint32) []uint32 {
	if v, ok := t.kongCards[seat]; ok && v != nil {
		return v
	}
	return []uint32{}
}

//获取吃牌
func (t *Desk) getChowCards(seat uint32) []uint32 {
	if v, ok := t.chowCards[seat]; ok && v != nil {
		return v
	}
	return []uint32{}
}

//获取玩家准备状态
func (t *Desk) getReady(seat uint32) bool {
	if v, ok := t.ready[seat]; ok {
		return v
	}
	return false
}

//获取玩家过圈状态
func (t *Desk) getSkip(seat uint32, huValue int64) bool {
	if v, ok := t.skip[seat]; ok {
		return algorithm.HuType(huValue, config.Opts().Ante) > v // 胡牌番值大于前一个胡牌番值则可以胡牌
	}
	return true
}

//获取剩余局数,结束时间
func (t *Desk) getRound() (uint32, uint32) {
	var expire uint32 = 0
	var now int64 = utils.Timestamp()
	if int64(t.data.Expire) > now {
		expire = t.data.Expire
	}
	var round uint32 = t.data.Round - t.round
	if round < 0 {
		round = 0
	}
	return round, expire
}

// 是否有操作值,有操作的玩家数量
func (t *Desk) optCount() int8 {
	var count int8
	for _, v := range t.opt {
		if v > 0 {
			count++
		}
	}
	return count
}

//房间消息广播
func (t *Desk) broadcast(msg interfacer.IProto) {
	for _, p := range t.players {
		p.Send(msg)
	}
}

//房间消息广播(除seat外)
func (t *Desk) broadcast_(seat uint32, msg interfacer.IProto) {
	for i, p := range t.players {
		if i != seat {
			p.Send(msg)
		}
	}
}

//--------操作
//摸牌
func (t *Desk) in(seat uint32, card byte) []byte {
	var cards []byte = t.getHandCards(seat)
	cards = append(cards, card)
	t.handCards[seat] = cards
	return cards
}

//玩家出牌
func (t *Desk) out(seat uint32, card byte) {
	var cards []byte = t.getHandCards(seat)
	cards = algorithm.Remove(card, cards)
	t.handCards[seat] = cards
}

//吃牌操作
func (t *Desk) chowing(seat uint32, c1, c2 byte) bool {
	var cards []byte = t.getHandCards(seat)

	if !algorithm.VerifyChow(cards, t.discard, c1, c2, t.luckyCard) {
		glog.Errorf("%x %x %x", t.discard, c1, c2)
		return false
	}
	cards = algorithm.Remove(c1, cards)
	cards = algorithm.Remove(c2, cards)
	var cs []uint32 = t.getChowCards(seat)
	cs = append(cs, algorithm.EncodeChow(c1, c2, t.discard))
	t.handCards[seat] = cards
	t.chowCards[seat] = cs
	return true
}

//碰牌操作
func (t *Desk) ponging(seat uint32, card byte) ([]byte, bool) {
	var cards []byte = t.getHandCards(seat)
	var isExist bool = algorithm.Exist(card, cards, 2)
	if !isExist {
		glog.Errorf("ponging card error -> %d", card)
		return cards, false
	}
	cards = algorithm.RemoveN(card, cards, 2)
	var cs []uint32 = t.getPongCards(seat)
	cs = append(cs, algorithm.EncodePeng(seat, card))
	t.handCards[seat] = cards
	t.pongCards[seat] = cs
	return cards, true
}

//暗扛操作
func (t *Desk) anKong(seat uint32, card byte) bool {
	var cards []byte = t.getHandCards(seat)
	var isExist bool = algorithm.Exist(card, cards, 4)
	if !isExist {
		glog.Errorf("anKong card error -> %d", card)
		return false
	}
	cards = algorithm.RemoveN(card, cards, 4)
	var cs []uint32 = t.getKongCards(seat)
	cs = append(cs, algorithm.EncodeKong(0, card, uint32(algorithm.AN_KONG)))
	t.handCards[seat] = cards
	t.kongCards[seat] = cs
	return true
}

//明杠操作
func (t *Desk) mingKong(seat uint32, card byte) bool {
	var cards []byte = t.getHandCards(seat)
	var isExist bool = algorithm.Exist(card, cards, 3)
	if !isExist {
		glog.Errorf("mingKong card error -> %d", card)
		return false
	}
	cards = algorithm.RemoveN(card, cards, 3)
	var cs []uint32 = t.getKongCards(seat)
	cs = append(cs, algorithm.EncodeKong(t.seat, card, uint32(algorithm.MING_KONG)))
	t.handCards[seat] = cards
	t.kongCards[seat] = cs
	return true
}

//补杠操作
func (t *Desk) buKong(seat uint32, card byte) bool {
	var cards []byte = t.getHandCards(seat)
	var isExist bool = algorithm.Exist(card, cards, 1)
	if !isExist {
		glog.Errorf("buKong card error -> %d", card)
		return false
	}
	var pongs []uint32 = t.getPongCards(seat)
	isExist = false
	for i, v := range pongs {
		_, c := algorithm.DecodePeng(v)
		if c == card {
			pongs = append(pongs[:i], pongs[i+1:]...)
			isExist = true
			break
		}
	}
	if !isExist {
		return false
	}
	cards = algorithm.Remove(card, cards)
	var cs []uint32 = t.getKongCards(seat)
	cs = append(cs, algorithm.EncodeKong(0, card, uint32(algorithm.BU_KONG)))
	t.handCards[seat] = cards
	t.kongCards[seat] = cs
	t.pongCards[seat] = pongs
	return true
}

// 摸牌检测,胡牌／暗杠／补杠
func (t *Desk) DrawDetect(card byte, cs []byte, ch, ps, ks []uint32, wildcard byte, seat uint32) int64 {
	//自摸胡检测
	handPongKong := algorithm.GetHandPongKong(cs, ch, ps, ks, wildcard)
	luangfeng := algorithm.ExistLuanFeng(handPongKong)
	color := algorithm.ExistOneSuit(handPongKong, wildcard)
	status := algorithm.ExistHu(cs, ch, ps, ks, wildcard, 0, color, luangfeng)

	//status := algorithm.ExistHu(cs, ch, ps, ks, wildcard, 0)
	if status > 0 {

		if t.tianHe(seat) > 0 { //天胡
			status = status | algorithm.TIAN_HU
		} else if t.diHe(seat) > 0 { //地胡
			status = status | algorithm.DI_HU
		} else if t.haidilaoHe() > 0 { // 海底捞
			status = status | algorithm.HU_HAI_LAO
		}
		// 13不靠没有爆头
		baotou := algorithm.ExistBaoTou(cs, ch, ps, ks, wildcard, card, true, status)
		if baotou > 0 {
			status = baotou

			// 爆头一定有财神
			if status&algorithm.HU_PENG_PENG > 0 {
				status = status & (^algorithm.HU_PENG_PENG)
				status = status | algorithm.HU_PENG_PENG_CAI
			}

			if status&algorithm.HU_ONE_SUIT > 0 {
				status = status & (^algorithm.HU_ONE_SUIT)
				status = status | algorithm.HU_ONE_SUIT_CAI
			}
		}

		threeW := algorithm.ThreeWildcard(cs, wildcard)
		if (threeW > 0) && (status&(^algorithm.HU)) == 0 {
			return 0
		}
		if threeW > 0 {
			status = algorithm.HU_3_CAI_SHEN | status
		}
		status = algorithm.ExistCaiShen(cs, status, wildcard, card, true)
		status |= algorithm.ZIMO
	}
	return status
}

// 打牌检测,胡牌, 接炮胡检测
func (t *Desk) DiscardHu(card byte, cs []byte, ch, ps, ks []uint32, wildcard byte, seat uint32) int64 {
	// 财神不能接炮胡
	if card == wildcard {
		return 0
	}
	handPongKong := algorithm.GetHandPongKong(cs, ch, ps, ks, wildcard)
	luangfeng := algorithm.ExistLuanFeng(handPongKong)
	color := algorithm.ExistOneSuit(handPongKong, wildcard)
	status := algorithm.ExistHu(cs, ch, ps, ks, wildcard, card, color, luangfeng)
	//status := algorithm.ExistHu(cs, ch, ps, ks, wildcard, card)
	if status > 0 {
		if t.diHe(seat) > 0 { //地胡
			status = status | algorithm.DI_HU
		}

		// 爆头不能炮胡,13不靠没有爆头
		if algorithm.ExistBaoTou(cs, ch, ps, ks, wildcard, card, false, status) > 0 {
			return 0
		}

		status = algorithm.ExistCaiShen(cs, status, wildcard, card, false)
		threeW := algorithm.ThreeWildcard(cs, wildcard)
		if (threeW > 0) && (status&(^algorithm.HU)) == 0 {
			return 0
		}
		if threeW > 0 {
			status = algorithm.HU_3_CAI_SHEN | status
		}

		status |= algorithm.PAOHU
	}
	return status
}
