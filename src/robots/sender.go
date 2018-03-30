package robots

import (
	"lib/utils"
	"crypto/md5"
	"encoding/hex"
	"protocol"

	"code.google.com/p/goprotobuf/proto"
	"game/algorithm"
	"github.com/golang/glog"
)

// 发送注册请求
func (c *Robot) SendRegist() {
	ctos := &protocol.CRegist{}
	ctos.Phone = proto.String(c.data.Phone)
	ctos.Nickname = proto.String(c.data.Nickname)
	h := md5.New()
	h.Write([]byte("piaohua")) // 需要加密的字符串为 123456
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos.Pwd = proto.String(pwd)
	c.Sender(ctos)
}

// 发送登录请求
func (c *Robot) SendLogin() {
	ctos := &protocol.CLogin{}
	ctos.Phone = proto.String(c.data.Phone)
	h := md5.New()
	h.Write([]byte("piaohua")) // 需要加密的字符串为 123456
	pwd := hex.EncodeToString(h.Sum(nil))
	ctos.Password = proto.String(pwd)
	c.Sender(ctos)


	//ctos:=&protocol.CHttpLogin{}
	//c.Sender(ctos)
}

// 获取玩家数据
func (c *Robot) SendUserData() {
	ctos := &protocol.CUserData{}
	ctos.Userid = proto.String(c.data.Userid)
	c.Sender(ctos)
}

// 玩家创建房间
func (c *Robot) SendCreate() {
	ctos := &protocol.CCreatePrivateRoom{}
	var a1 []uint32 = []uint32{4, 8, 16}
	var a2 []uint32 = []uint32{1, 5, 9}
	var i int32 = utils.RandInt32N(3) //随机
	ctos.Round = proto.Uint32(a1[i])
	ctos.Rtype = proto.Uint32(a2[i])
	ctos.Ante = proto.Uint32(1)
	ctos.Ma = proto.Uint32(1)
	ctos.Payment = proto.Uint32(1)
	ctos.Rname = proto.String("")
	ctos.Maizi = proto.Bool(false)
	c.Sender(ctos)
}

// 玩家进入房间
func (c *Robot) SendEntry() {
	if c.code == "create" { //表示创建房间
		c.SendCreate() //创建一个房间
	} else { //表示进入房间
		glog.Infoln("表示进入房间",c.code)
		ctos := &protocol.CEnterSocialRoom{}
		ctos.Invitecode = proto.String(c.code)
		c.Sender(ctos)
	}
}

// 玩家打骰子
func (c *Robot) SendBroken() {
	ctos := &protocol.CBroken{}
	c.Sender(ctos)
}



// 玩家买子
func (c *Robot) SendMaizi() {
	ctos := &protocol.CMaiZi{
		Count:proto.Uint32(0),
	}
	c.Sender(ctos)
}


// 玩家准备
func (c *Robot) SendReady() {
	ctos := &protocol.CReady{}
	ctos.Ready = proto.Bool(true)
	c.Sender(ctos)
}

// 离开
func (c *Robot) SendLeave() {
	ctos := &protocol.CPrivateLeave{}
	c.Sender(ctos)
}

// 庄家出牌
func (c *Robot) SendDiscard2() {
	utils.Sleep(6) //展示动画时间
	c.SendDiscard()
}

// 自己出牌
func (c *Robot) SendDiscard() {
	ctos := &protocol.CDiscard{}
	card := algorithm.SearchDirtyCard(c.cards,0)
	if card == 0{
		card = c.cards[0]
		glog.Errorf("机器人提交的牌值为0 %+x",c.cards)

	}
	ctos.Card = proto.Uint32(uint32(card))
	ctos.Value = proto.Int64(0)

	utils.SleepRand(3)
	c.Sender(ctos)
}

// 玩家碰杠操作
func (c *Robot) SendOperate(card uint32, value int64) {
	ctos := &protocol.COperate{}
	ctos.Card = proto.Uint32(card)
	ctos.Value = proto.Int64(value)
	utils.SleepRand(3)
	c.Sender(ctos) // send data
}

// 胡牌
func (c *Robot) SendHu() {
	ctos := &protocol.CHu{}
	c.Sender(ctos) // send data
}
