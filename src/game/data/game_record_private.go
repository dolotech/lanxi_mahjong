package data

import (
	"time"
	"errors"
	"gopkg.in/mgo.v2/bson"
)

//牌局记录
type GameOverRecord struct {
	RoomId     uint32   `bson:"_id"`           //房间ID
	Invitecode string   `bson:"Invitecode"`    //房间邀请码
	TotalRound uint32   `bson:"TotalRound"`    //牌局局数
	Ante       uint32   `bson:"Ante"`          //底分
	Userids    []string `bson:"Userids"`       //用户账号id
	Ctime      uint32   `bson:"Ctime"`         //房间创建时间
	Cid        string   `bson:"Create_userid"` //房间创建人
	Expire     uint32   `bson:"Expire"`        //牌局设定的过期时间
	Create     uint32   `bson:"Create"`        //记录创建时间
	Update     uint32   `bson:"Update"`        //记录更新时间
	Rounds     []*GameOverRoundRecord `bson:"Rounds"`
}

type GameOverUserRecord struct {
	Seat    uint32 `bson:"Seat"`    //用户位置
	Coin    int32  `bson:"Coin"`    //输赢
	Userid  string `bson:"Userid"`  //用户账号id
	Huvalue uint32 `bson:"Huvalue"` //胡牌状态
	HuCard  uint32 `bson:"HuCard"`  //胡的牌
}

type GameOverRoundRecord struct {
	Users  []*GameOverUserRecord `bson:"Users"`
	Round  uint32 `bson:"Round"`  //牌局局数
	Zhuang uint32 `bson:"Zhuang"` //庄家的座位
	Hutype uint32 `bson:"Hutype"` //1:自摸，2:炮胡，3:黄庄
}

func (this *GameOverRecord) Add() error {
	this.Create = uint32(time.Now().Unix())
	return C(_GAMEOVER_RECORD).Insert(this)
}

func (this *GameOverRecord) Get(userid string) error {
	if userid == "" {
		return errors.New("userid can not empty")
	}

	if this.RoomId == 0 {
		return errors.New("roomid can not 0")
	}
	//return C(_GAMEOVER_RECORD).FindId(this.RoomId).Sort("-CTime").One(this)
	return C(_GAMEOVER_RECORD).FindId(this.RoomId).Sort("-_id").One(this)
}

func (this *GameOverRoundRecord) Push(roomId uint32) error {
	update := uint32(time.Now().Unix())
	C(_GAMEOVER_RECORD).Update(bson.M{"_id": roomId}, bson.M{"$set": bson.M{"Update": update}})
	return C(_GAMEOVER_RECORD).UpdateId(roomId, bson.M{"$push": bson.M{"Rounds": this}})
}

//获取记录
type GameOverRecords []*GameOverRecord

func (this *GameOverRecords) Get(userid string, page int, limit int) error {
	if userid == "" {
		return errors.New("userid can not empty")
	}

	if page < 1 {
		page = 1
	}
	if limit < LIMIT_MIN {
		limit = LIMIT_MIN
	} else if limit > LIMIT_MAX {
		limit = LIMIT_MAX
	}
	//return C(_GAMEOVER_RECORD).Find(bson.M{"Userids":userid}).Sort("-CTime").Skip((page - 1) * limit).Limit(limit).All(this)
	return C(_GAMEOVER_RECORD).Find(bson.M{"Userids": userid}).Sort("-_id").Skip((page - 1) * limit).Limit(limit).All(this)
}

//---------------------------------------------------------------------------------------

type GameRecordPrivate struct {
	Id            string `bson:"_id"` //房间ID
	RoomId        uint32 `bson:"RoomId"`
	Rtype         uint32 `bson:"Rtype"`     //房间类型
	Ante          uint32 `bson:"Ante"`      //底分
	Userids       []string `bson:"Userids"` //用户账号id
	Users         []*GameRecordPrivateUser `bson:"Users"`
	Ji            byte    `bson:"Ji"`     //
	HeroJi        uint32  `bson:"HeroJi"` // 0:无，1:英雄鸡，2：责任鸡，3：责任鸡碰方
	HuCard        byte  `bson:"HuCard"`
	Ctime         uint32 `bson:"Ctime"`
	Zhuang        uint32 `bson:"Zhuang"`
	Invitecode    string `bson:"Invitecode"`
	TotalRound    uint32 `bson:"TotalRound"`        //牌局局数
	Payment       uint32        `bson:"Payment"`    //付费方式1=AA or 0=房主支付
	Create_userid string     `bson:"Create_userid"` //房间创建人
	Expire        uint32      `bson:"Expire"`       //牌局设定的过期时间
	Rname         string `bson:"Rname"`             //房间名字
	Round         uint32 `bson:"Round"`             //牌局局数
}

type GameRecordPrivateUser struct {
	Seat           uint32 `bson:"Seat"`   //用户位置
	Coin           int32  `bson:"Coin"`   //输赢
	Userid         string `bson:"Userid"` //用户账号id
	StartHandCards []byte `bson:"StartHandCards"`
	EndHandCards   []byte `bson:"EndHandCards"`
	OutCards       []byte `bson:"OutCards"`
	Peng           []uint32 `bson:"Peng"`    //
	Kong           []uint32 `bson:"Kong"`    //
	HuValue        uint32 `bson:"HuValue"`   //胡牌状态
	TingValue      uint32 `bson:"TingValue"` //停牌状态
	HuCard         uint32 `bson:"HuCard"`
}

func (this *GameRecordPrivate) Save() error {
	this.Ctime = uint32(time.Now().Unix())
	return C(_GAMEOVER_PRIVATE).Insert(this)
}
