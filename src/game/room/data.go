package room

import (
	"game/data"
	"lib/utils"
	"time"
)

func NewDeskData(rid, round, expire, ante, cost uint32, creator, invitecode string, maizi bool) *DeskData {
	return &DeskData{
		Rid:    rid,
		Ante:   ante,
		Cost:   cost,
		Cid:    creator,
		Expire: expire,
		Round:  round,
		Code:   invitecode,
		CTime:  uint32(utils.Timestamp()),
		Score:  make(map[string]int32),
		MaiZi:  maizi,
	}
}

type DeskData struct {
	Rid    uint32           //房间ID
	Cid    string           //房间创建人
	Expire uint32           //牌局设定的过期时间
	Code   string           //房间邀请码
	Round  uint32           // 总牌局数
	Ante   uint32           //私人房底分
	Cost   uint32           //创建消耗
	CTime  uint32           //创建时间
	Score  map[string]int32 //私人局用户战绩积分
	MaiZi  bool             // 是否买子
}

//// 开房记录
//type CreateRoomRecord struct {
//	Userid     string `bson:"Userid"`     // 开房玩家id
//	Date       uint32 `bson:"Date"`       // 时间戳
//	Roomid     uint32 `bson:"Roomid"`     // 房间id
//	Rtype      uint32 `bson:"Rtype"`      // 房间类型
//	Round      uint32 `bson:"Round"`      // 房间局数
//	Invitecode string `bson:"Invitecode"` // 房间邀请码
//	Maizi      bool   `bson:"Maizi"`      // 是否买子
//	Expire     uint32 `bson:"Expire"`     // 房间过期时间
//	Cost       uint32 `bson:"Cost"`       // 房间消耗
//	Ante       uint32 `bson:"Ante"`       // 底分
//}

func (this *DeskData) Convert() *data.CreateRoomRecord {
	return &data.CreateRoomRecord{
		Userid:     this.Cid,
		Roomid:     this.Rid,
		Date:       time.Unix(int64(this.CTime), 0),
		Round:      this.Round,
		Ante:       this.Ante,
		Cost:       this.Cost,
		Maizi:      this.MaiZi,
		Invitecode: this.Code,
		Expire:      time.Unix(int64(this.Expire), 0),
	}
}
