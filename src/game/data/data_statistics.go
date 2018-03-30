package data

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// ONLINE_COUNTER
// ROOM_COUNTER
// PLAYING_ROOM_COUNTER
// DAILY_ROOM_COUNTER_EIGHT
// DAILY_ROOM_MAIZI_COUNTER_EIGHT
// DAILY_ROOM_NORMAL_COUNTER_EIGHT
// DAILY_ROOM_CONSUME_CARD_COUNTER_EIGHT
// DAILY_ROOM_COUNTER_SIXTEEN
// DAILY_ROOM_MAIZI_COUNTER_SIXTEEN
// DAILY_ROOM_NORMAL_COUNTER_SIXTEEN
// DAILY_ROOM_CONSUME_CARD_COUNTER_SIXTEEN

type Statistics struct {
	Name  string    `bson:"Name"`  // 统计在线的结构体名称
	Date  time.Time `bson:"Date"`  // 时间戳
	Total uint32    `bson:"Total"` // 当前在线的玩家的数据
}

// 获取统计数据所以数据
func (this *Statistics) Get() error {
	//if this.Name == "" {
	//	return errors.New("name can not empty")
	//}
	return C(_STATISTICS).FindId(this.Name).One(this)
}

func (this *Statistics) Save() error {
	//if this.Name == "" {
	//	return errors.New("name  can not empty")
	//}
	return C(_STATISTICS).Insert(this)
}

// 获取指手机用户的所以数据
func (this *Statistics) GetLatestOne() error {
	if this.Name == "" {
		return errors.New("name can not empty")
	}
	return C(_STATISTICS).Find(bson.M{"Name": "ONLINE_COUNTER"}).Sort("-Date").Limit(1).One(&this)
}

// 开房记录
type CreateRoomRecord struct {
	Userid     string    `bson:"Userid"`     // 开房玩家id
	Roomid     uint32    `bson:"Roomid"`     // 房间id
	Round      uint32    `bson:"Round"`      // 房间局数
	Invitecode string    `bson:"Invitecode"` // 房间邀请码
	Maizi      bool      `bson:"Maizi"`      // 是否买子
	Cost       uint32    `bson:"Cost"`       // 房间消耗
	Ante       uint32    `bson:"Ante"`       // 底分
	Expire     time.Time `bson:"Expire"`     // 房间过期时间
	Date       time.Time `bson:"Date"`       // 时间创建时间戳
}

func (this *CreateRoomRecord) Save() error {
	this.Date = time.Now()
	return C(_ROOM_CREATE_RECORD).Insert(this)
}
