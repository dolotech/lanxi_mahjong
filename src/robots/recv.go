package robots

import (
	"game/algorithm"
	"protocol"

	"github.com/golang/glog"
	"time"
	"code.google.com/p/goprotobuf/proto"
)

func init() {
	regist(&protocol.SLogin{}, recvLogin)
	regist(&protocol.SEnterSocialRoom{}, recvComein)
	regist(&protocol.SDeal{}, recvDeal)
	regist(&protocol.SCreatePrivateRoom{}, recvCreate)
	regist(&protocol.SDiscard{}, recvDiscard)
	regist(&protocol.SDraw{}, recvDraw)
	regist(&protocol.SPrivateOver{}, recvGameover)
	regist(&protocol.SPrivateLeave{}, recvLeave)
	regist(&protocol.SUserData{}, recvdata)
	regist(&protocol.SZhuang{}, recvZhuang)
	regist(&protocol.SZhuangDeal{}, recvZhuangdeal)
	regist(&protocol.SOperate{}, recvOperate)
	regist(&protocol.SPengKong{}, recvPengKong)
	regist(&protocol.SRegist{}, recvRegist)
	regist(&protocol.SLaunchVote{}, launchVote)
	regist(&protocol.SHttpLogin{}, httpLogin)
}


// 机器人永远支持解散房间
func launchVote(stoc *protocol.SLaunchVote, c *Robot) {
	ctos:= &protocol.CVote{
		Vote:proto.Uint32(0),
	}
	c.Sender(ctos) // send data
}
// 接收到服务器登录返回
func recvRegist(stoc *protocol.SRegist, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch {
	case errcode == 13016:
		c.SendLogin() //已经注册,登录
	case errcode == 0:
		var uid string = stoc.GetUserid()
		c.data.Userid = uid
		c.SendLogin() //注册成功,登录
	default:
		glog.Infof("regist err -> %d", errcode)
	}
}
// 接收到服务器http登录返回
func httpLogin(stoc *protocol.SHttpLogin, c *Robot) {

	var errcode uint32 = stoc.GetError()
	glog.Infoln(stoc.GetCode(),errcode)
}
// 接收到服务器登录返回
func recvLogin(stoc *protocol.SLogin, c *Robot) {
	var errcode uint32 = stoc.GetError()

	switch {
	case errcode == 13009:
		glog.Infof("login passwd err -> %s", c.data.Phone)
	case errcode == 0:
		c.Logined() //登录成功
		c.data.Userid = stoc.GetUserid()
		c.SendUserData() // 获取玩家数据
	default:
		glog.Infof("login err -> %d", errcode)
	}
}

// 接收到玩家数据
func recvdata(stoc *protocol.SUserData, c *Robot) {
	var errcode uint32 = stoc.GetError()
	if errcode != 0 {
		glog.Infof("get data err -> %d", errcode)
	}
	data := stoc.GetData()
	// 设置数据
	c.data.Userid = data.GetUserid()     // 用户id
	c.data.Nickname = data.GetNickname() // 用户昵称
	c.data.Sex = data.GetSex()           // 用户性别,男1 女2 非男非女3
	c.data.Coin = data.GetCoin()         // 金币
	c.data.Exp = data.GetExp()           // 经验
	c.data.Diamond = data.GetDiamond()   // 钻石
	c.data.Ticket = data.GetTicket()     // 入场券
	c.data.Exchange = data.GetExchange() // 兑换券
	c.data.Vip = data.GetVip()
	//查找房间-进入房间
	c.SendEntry()
}

// 离开房间
func recvLeave(stoc *protocol.SPrivateLeave, c *Robot) {
	var seat uint32 = stoc.GetSeat()
	if seat == c.seat {
		c.Close() //下线
	}
	if seat >= 1 && seat <= 4 && seat != c.seat {
		c.SendLeave() //离开
	}
}

// 创建房间
func recvCreate(stoc *protocol.SCreatePrivateRoom, c *Robot) {
	var errcode uint32 = stoc.GetError()
	switch {
	case errcode == 0:
		var code string = stoc.GetRdata().GetInvitecode()
		if code != "" {
			glog.Infof("create room code -> %s", code)
			c.code = code       //设置邀请码
			c.SendEntry()       //进入房间
			Msg2Robots(code, 3) //创建房间成功,邀请3个人进入
		} else {
			glog.Errorf("create room code empty -> %s", code)
		}
	default:
		glog.Infof("create room err -> %d", errcode)
		c.Close() //进入出错,关闭
	}
}

// 进入房间
func recvComein(stoc *protocol.SEnterSocialRoom, c *Robot) {
	var errcode uint32 = stoc.GetError()

	glog.Infoln("机器人 进入房间", stoc.String())
	switch {
	case errcode == 0:
		c.seat = stoc.GetPosition()
		c.room = stoc.Room
		time.AfterFunc(time.Second, func() {
			c.SendReady() //准备

			if c.room.GetMaizi() {
				<-time.After(time.Second)
				c.SendMaizi()
			}
		})
	default:
		glog.Infof("comein err -> %d", errcode)
		c.Close() //进入出错,关闭
	}
}

// 打庄
func recvZhuang(stoc *protocol.SZhuang, c *Robot) {
	var zhuang uint32 = stoc.GetZhuang()
	switch {
	case zhuang == c.seat:
		c.SendBroken() //打骰了发牌
	default:
	}
}

// 非庄家发牌
func recvDeal(stoc *protocol.SDeal, c *Robot) {
	c.cards = stoc.GetCards()
}

// 庄家发牌,包含两个骰子数字
func recvZhuangdeal(stoc *protocol.SZhuangDeal, c *Robot) {
	value := stoc.GetValue()
	c.cards = stoc.GetCards()
	if value == 0 {
		c.SendDiscard2()
	} else {
		operate(value, 0, c) //胡或暗杠
	}
}

// 结束
func recvGameover(stoc *protocol.SPrivateOver, c *Robot) {
	var round uint32 = stoc.GetRound()
	if round == 0 {
		c.Close() //结束下线
	} else {
		time.AfterFunc(time.Second, func() {
			c.SendReady() //准备

			if c.room.GetMaizi() {
				<-time.After(time.Second)
				c.SendMaizi()
			}
		})
	}
}

// 此协议向有碰和明杠的玩家主动发送
func recvPengKong(stoc *protocol.SPengKong, c *Robot) {
	card := stoc.GetCard()
	value := stoc.GetValue()
	operate(value, card, c) //胡碰杠吃
}

// 操作结果返回,同步手牌
func recvOperate(stoc *protocol.SOperate, c *Robot) {
	seat := stoc.GetSeat()               // 操作位置
	card := stoc.GetCard()               // 被碰或杠牌的牌值
	value := stoc.GetValue()             // 碰或值杠,统一掩码标示
	discontinue := stoc.GetDiscontinue() // 抢杠
	if seat != c.seat { // 不是自己操作
		return
	}
	// 根据不同操作同步cards
	switch {
	case value&algorithm.CHOW > 0:
		c1, c2, _ := algorithm.DecodeChow(card) //解码
		c.cards = RemoveN(byte(c1), c.cards, 1)
		c.cards = RemoveN(byte(c2), c.cards, 1)
		c.SendDiscard() //出牌
	case value&algorithm.PENG > 0:
		c.cards = RemoveN(byte(card), c.cards, 2)
		c.pongCards = append(c.pongCards,byte(card))
		c.SendDiscard() //出牌
	case value&algorithm.MING_KONG > 0:
		c.cards = RemoveN(byte(card), c.cards, 3)
	case value&algorithm.AN_KONG > 0:
		c.cards = RemoveN(byte(card), c.cards, 4)
	case value&algorithm.BU_KONG > 0:
		c.cards = RemoveN(byte(card), c.cards, 1)
		c.pongCards = RemoveN(byte(card), c.pongCards, 1)
	case discontinue == algorithm.QIANG_GANG:
		c.SendHu()
	case discontinue > 0: //不同地址掩码值不一样
		c.SendHu()
	default:
		glog.Infof("SOperate err -> %d", value)
	}
}

// 玩家出牌广播
func recvDiscard(stoc *protocol.SDiscard, c *Robot) {
	var errcode uint32 = stoc.GetError()
	var seat uint32 = stoc.GetSeat()
	var card uint32 = stoc.GetCard()
	 value := stoc.GetValue()
	if errcode != 0 {
		glog.Infof("Discard err -> %d", errcode)
		return
	}
	if seat == c.seat { //自己出牌
		c.cards = RemoveN(byte(card), c.cards, 1) //移除
		return
	}
	if value == 0 { //没有操作
		return
	}
	operate(value, card, c) //胡碰杠吃
}

// 抓牌
func recvDraw(stoc *protocol.SDraw, c *Robot) {
	//card := stoc.GetCard()
	//value := stoc.GetValue()
	c.cards = stoc.GetCards()
	//PrintCards(c.cards)
	//if value == 0 {
		c.SendDiscard()
	//} else {
	//	operate(value, card, c) //胡或暗杠
	//}
}

//操作一定要正确,服务器有些操作没做验证
func operate(value int64, card uint32, c *Robot) {
	switch {
	case value&algorithm.HU > 0:
		c.SendHu()
	case value&algorithm.MING_KONG > 0:
		c.SendOperate(card, algorithm.MING_KONG)
	case value&algorithm.BU_KONG > 0:
		for _,v:=range c.cards{
			for _,p:=range  c.pongCards{
				if p == v {
					c.SendOperate(uint32(v), algorithm.BU_KONG)
					break
				}
			}
		}
	case value&algorithm.AN_KONG > 0:
		var kongcard byte = findKong(c.cards) //找一个暗扛
		c.SendOperate(uint32(kongcard), algorithm.AN_KONG)
	case value&algorithm.PENG > 0:
		c.SendOperate(card, algorithm.PENG)
	case value&algorithm.CHOW > 0:
		c.SendOperate(card, 0)	// 机器人不吃牌
		//c.cards = RemoveN(byte(chowcard[0]>>8&0xF), c.cards, 1)
		//c.cards = RemoveN(byte(chowcard[0]&0xF), c.cards, 1)
		//c.SendOperate(chowcard[0], algorithm.CHOW)
	}
}

//判断是否有4个相同的牌
func findKong(cards []byte) byte {
	m := make(map[byte]int, len(cards))
	for _, v := range cards {
		if i, ok := m[v]; ok {
			if i == 3 {
				return v
			}
			m[v] = i + 1
		} else {
			m[v] = 1
		}
	}
	return 0
}

// 移除n个牌
func RemoveN(c byte, cs []byte, n int) []byte {
	for n > 0 {
		for i, v := range cs {
			if c == v {
				cs = append(cs[:i], cs[i+1:]...)
				break
			}
		}
		n--
	}
	return cs
}
