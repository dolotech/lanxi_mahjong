package data

import (
	"time"
	"lib/utils"
)

const (
	RESTYPE4      =4          //私人局
)

type DataResChange struct {
	Userid   string `bson:"Userid"`   //玩家ID
	Kind     uint32 `bson:"Kind"`     //道具、货币种类
	Time     uint32 `bson:"Time"`     //变动时间
	Channel  uint32 `bson:"Channel"`  //获取、扣除渠道
	Residual uint32 `bson:"Residual"` //剩余量
	Count    int32 `bson:"Count"`  // 变数量
}


type DataResChanges []*DataResChange
func (this *DataResChanges) Save(userid string) error {
	var list []interface{}
	for _,v:=range *this{
		v.Time =uint32(time.Now().Unix())
		v.Userid = userid
		list =append(list,utils.Struct2Map(v))
	}
	return C(_RESOURCE_RECORD).Insert(list...)
}
