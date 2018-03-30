package config

import (
	"bytes"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
)

// Config 配置类型
type Config struct {
	Operate_timeout     int
	Trusteeship_timeout int

	Db_addr        string
	Server_port    string
	Server_ip      string
	Oprof_port     string
	Web_port       string
	Pay_port       string
	Pay_wx_pattern string

	Share_addr string
	Server_id  int
	Version    string

	AdminPort string

	// 可选买子数量
	SelectMaiziCt []uint32
	// 可选牌局数量
	SelectGameCountCt []uint32

	// 底分
	Ante uint32

	// 房卡价格
	Price  []uint32
}

// Opts Config 默认配置
var opts = Config{}

// Opts 获取配置
func Opts() Config {
	return opts
}

//创建房间房卡
type RoomCard struct {
	PlayCt   int32 // 创建房间局数
	NeedCard int32 // 所需房卡数
}

// Room 房间
type Room struct {
	Name       string   // 房间名称
	BasMoney   int64    // 底注
	MinMoney   int64    // 最小准入
	Tip        int      // 台费 （百分比 0 - 100）
	RaiseChips [4]int64 // 固定加注额度
}

// 房卡
type FangKa struct {
	ID     int
	Name   string
	Count  int
	Price  int
	Reward int
}

// ParseToml 解析配置文件
func ParseToml(file string) error {
	glog.Infoln("读取配置文件 ...")
	// 如果配置文件不存在
	if _, err := os.Stat(file); os.IsNotExist(err) {
		buf := new(bytes.Buffer)
		if err := toml.NewEncoder(buf).Encode(Opts()); err != nil {
			return err
		}
		glog.Infoln("没有找到配置文件，创建新文件 ...")
		// logs.Info(buf.String())
		return ioutil.WriteFile(file, buf.Bytes(), 0644)
	}
	var conf Config
	_, err := toml.DecodeFile(file, &conf)
	if err != nil {
		return err
	}
	opts = conf
	glog.Infoln("----------->，config.Opts().FangKas:%v", opts)

	return nil
}
