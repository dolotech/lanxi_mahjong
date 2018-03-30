package data

import (
	"errors"
	"gopkg.in/mgo.v2/bson"
	"lib/utils"
)

// 批量更改玩家的经济资源
func (this *User) UpdateResource(value map[string]int32) error {
	return C(_USER).UpdateId(this.Userid, bson.M{"$set": value})
}

// 获取昵称，性别和头像
func (this *User) GetPhotoSexName() error {
	return C(_USER).FindId(this.Userid).Select(bson.M{"Photo": 1, "Sex": 1, "Nickname": 1, "_id": -1}).One(&this)
}

func (this *User) GetPhotoFromDB() (string, error) {
	var user User
	err := C(_USER).FindId(this.Userid).Select(bson.M{"Photo": 1, "_id": -1}).One(&user)
	return user.Photo, err
}

func (this *User) UpdatePhoto() error {
	return C(_USER).UpdateId(this.Userid, bson.M{"$set": bson.M{"Photo": this.Photo}})
}

// 获取指手机用户的所以数据
func (this *User) Get() error {
	if this.Userid == "" {
		return errors.New("Userid can not empty")
	}
	return C(_USER).FindId(this.Userid).One(this)
}

func (this *User) Save() error {
	if this.Userid == "" {
		return errors.New("Userid  can not empty")
	}
	return C(_USER).Insert(this)
}

func (this *User) GetByWechat(wechat string) error {
	return C(_USER).Find(bson.M{"Wechat_uid": wechat}).One(this)
}

func (this *User) ExistsPhone(phone string) bool {
	count, _ := C(_USER).Find(bson.M{"Phone": phone}).Count()
	return count > 0
}

func (this *User) UpdateSex() error {
	return C(_USER).UpdateId(this.Userid, bson.M{"$set": bson.M{"Sex": this.Sex}})
}

func (this *User) UpdateParent() error {
	return C(_USER).UpdateId(this.Userid, bson.M{"$set": bson.M{"Build": this.Build}})
}

func (this *User) UpdateNickname() error {
	return C(_USER).UpdateId(this.Userid, bson.M{"$set": bson.M{"Nickname": this.Nickname}})
}
func (this *User) GetByPhone() string {
	var user User
	err := C(_USER).FindId(this.Userid).Select(bson.M{"Phone": 1, "_id": -1}).One(&user)
	if err != nil {
		return ""
	}
	return user.Phone
}

func (this *User) UpdatePWD(pwd string) error {
	if this.Userid == "" {
		return errors.New("Userid can not  empty")
	}
	tmp := &User{}
	err := C(_USER).FindId(this.Userid).Select(bson.M{"Auth": 1, "_id": -1}).One(tmp)
	if err != nil {
		return err
	}
	passwd := utils.Md5(pwd + tmp.Auth)
	return C(_USER).UpdateId(this.Userid, bson.M{"$set": bson.M{"Pwd": passwd}})
}

func (this *User) VerifyPwdByPhone(pwd string) bool {
	if this.Phone == "" {
		return false
	}
	tmp := &User{}
	err := C(_USER).Find(bson.M{"Phone": this.Phone}).Select(bson.M{"_id": 1, "Auth": 1, "Pwd": 1}).One(tmp)
	if err != nil {
		return false
	}

	if utils.Md5(pwd+tmp.Auth) == tmp.Pwd {
		this.Userid = tmp.Userid
		return true
	}

	return false
}

//  用户登陆密码验证
func (this *User) PWDIsOK(pwd string) bool {
	if this.Userid == "" {
		return false
	}
	tmp := &User{}
	err := C(_USER).FindId(this.Userid).Select(bson.M{"Auth": 1, "Pwd": 1, "_id": -1}).One(tmp)
	if err != nil {
		return false
	}

	if utils.Md5(pwd+tmp.Auth) == tmp.Pwd {
		return true
	}

	return false
}

type User struct {
	Userid        string `bson:"_id"`         // 用户id
	Nickname      string `bson:"Nickname"`    // 用户昵称
	Sex           uint32 `bson:"Sex"`         // 用户性别,男1 女2 非男非女3
	Email         string `bson:"Email"`       // 绑定的邮箱地址
	Phone         string `bson:"Phone"`       // 绑定的手机号码
	Auth          string `bson:"Auth"`        // 密码验证码
	Pwd           string `bson:"Pwd"`         // MD5密码
	Create_ip     uint32 `bson:"Create_ip"`   // 注册账户时的IP地址
	Create_time   uint32 `bson:"Create_time"` // 注册时间
	Terminal      string `bson:"Terminal"`    // 终端类型名字
	Status        uint32 `bson:"Status"`      // 正常1  锁定2  黑名单3
	Address       string `bson:"Address"`     //物理地址
	Photo         string `bson:"Photo"`       //头像
	Qq_uid        string `bson:"Qq_uid"`      //
	Wechat_uid    string `bson:"Wechat_uid"`
	Microblog_uid string `bson:"Microblog_uid"`
	Platform      uint32 `bson:"Platform"`

	Robot    bool   `bson:"Robot"`    //是否是机器人
	RoomCard uint32 `bson:"RoomCard"` //房卡
	Build    string `bson:"Build"`    //绑定id

	Longitude string `bson:"Longitude"` // 经度
	Latitude  string `bson:"Latitude"`  // 纬度
}
