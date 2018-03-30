package data

import (
	"errors"
	"time"
)

type DataUserActive struct {
	Userid string    `bson:"Userid"` // 用户账号id
	IP     uint32    `bson:"IP"`     // 登录IP
	Time   time.Time `bson:"Time"`   // 时间戳
	Action uint32    `bson:"Action"` // 1:上线，2：下线
}

// 记录登陆时间,没有该玩家数据则插入
func (this *DataUserActive) Login() error {
	this.Time = time.Now()
	this.Action = 1
	if this.Userid != "" {
	} else {
		return errors.New("user id is empty!")
	}

	return C(_LOGIN_LOG_OUT_RECORD).Insert(this)
}

// 记录退出时间，并累积在线时长
func (this *DataUserActive) Logout() error {
	this.Time = time.Now()
	this.Action = 2
	if this.Userid != "" {
	} else {
		return errors.New("user id is empty!")
	}
	return C(_LOGIN_LOG_OUT_RECORD).Insert(this)
}
