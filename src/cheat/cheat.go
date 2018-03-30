package cheat

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"game/algorithm"
	"config"
	"game/room"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"lib/utils"
	"game/resource"
)

//用于测试指定房间牌型
type RoomCheat struct {
	RobotNum int      `json:"robot"`    // 机器人个数
	Wildcard int      `json:"wildcard"` // 财神牌
	RoomId   int      `json:"roomid"`   // 房间号(邀请码)
	Seat     [][]byte `json:"seat"`     // 每个座位的手牌 第一个是座家的，其它三家随机
	Card     []byte   `json:"card"`     // 剩余的牌
}

func Run(port string) {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	//e.Use(middleware.CSRF())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("1M"))

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusTemporaryRedirect, "/index.html")
	})

	e.Static("/", "./assets")

	e.POST("/create", create)
	e.POST("/autocreate", autoCreate)
	e.GET("/reloadcsv", reloadcsv)
	e.GET("/reloadconfig", reloadconfig)

	// 后台系统发房卡接口
	e.POST("/roomcard", updateRoomCard)

	e.GET("/reloadconfig", reloadconfig)

	e.Start(port)
}

type H map[string]interface{}

const KEY = "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"

type RoomInfoReq struct {
	RoomId int      `json:"roomid"` // 房间号(邀请码)
}

func updateRoomCard(c echo.Context) error {
	c.Response().CloseNotify()
	body, _ := ioutil.ReadAll(c.Request().Body)

	shipper := &Shipper{}
	err := json.Unmarshal(body, shipper)
	if err != nil {
		return c.JSON(http.StatusOK, H{"code": 1007})
	}

	sign := utils.Md5(KEY + strconv.Itoa(int(shipper.Timestamp)) + shipper.Userid +
		shipper.Transid + strconv.Itoa(int(shipper.Kind)) +
		strconv.Itoa(int(shipper.Count)) + strconv.Itoa(int(shipper.CurCount)))

	if sign != shipper.Key {
		return c.JSON(http.StatusOK, H{"code": 1003})
	}
	resource.NotifyChangeRes(shipper.Userid, shipper.Kind, int32(shipper.CurCount))

	return c.JSON(http.StatusOK, H{"code": 0})
}

func reloadcsv(c echo.Context) error {
	//csv.Init()
	glog.Infoln("reloadcsv")
	return c.JSON(http.StatusOK, H{"msg": "数据提交成功"})
}
func reloadconfig(c echo.Context) error {
	config.ParseToml("./game_config.toml")
	glog.Infoln("reloadconfig")
	return c.JSON(http.StatusOK, H{"msg": "数据提交成功"})
}

//返回{"code":0}表示成功code>0 表示 失败 ，如果失败，你要定时循环补单，直到逻辑服返回{"code":0}
//
//code错误码：
//
//1002：订单号相同我会忽略下单并返回
//1003：签名校验失败
//1007：数据解包失败
//1004：玩家不存在
//1005：商品不存在
//1006：订单过期(发货时间戳超过30分钟)

//1008：玩家不在线
// Key = Md5(SIGN+Timestamp+Userid + Transid + Kind + Count +CurCount)
type Shipper struct {
	Transid  string `json:"transid"`  // 订单号
	Key      string `json:"key"`      //  签名
	Userid   string `json:"userid"`   // 玩家id
	Kind     uint32 `json:"kind"`     // 商品类型4:房卡
	Count    uint32 `json:"count"`    // 商品数量
	CurCount uint32 `json:"curcount"` //玩家剩余商品数量

	Timestamp uint32 `json:"timestamp"` // 发货时间戳
}

func create(c echo.Context) error {
	c.Response().CloseNotify()
	body, _ := ioutil.ReadAll(c.Request().Body)

	ar := &RoomCheat{}
	err := json.Unmarshal(body, ar)
	if err != nil {
		glog.Errorln("json.Unmarshal failed ", err)
		return c.JSON(http.StatusOK, H{"msg": err.Error()})
	}
	r := room.Get(strconv.Itoa(ar.RoomId))

	glog.Errorln("_______邀请机器人进入房间--------", r, ar.RoomId)
	if r == nil {
		return c.JSON(http.StatusOK, H{"msg": "房间不存在"})
	}

	if ar.Wildcard == 0 {
		return c.JSON(http.StatusOK, H{"msg": "没有填写财神牌"})
	}

	if len(ar.Seat) != 4 {
		return c.JSON(http.StatusOK, H{"msg": "玩家数量不足"})
	}

	if len(ar.Seat[0]) != 14 {
		return c.JSON(http.StatusOK, H{"msg": "庄家手牌数量不足"})
	} else {
		for _, v := range ar.Seat[0] {
			if !algorithm.Legal(byte(v)) {
				return c.JSON(http.StatusOK, H{"msg": "庄家手牌，存在非法牌值"})
			}
		}
	}

	if len(ar.Seat[1]) != 13 {
		return c.JSON(http.StatusOK, H{"msg": "闲家1手牌数量不足"})
	} else {
		for _, v := range ar.Seat[1] {
			if !algorithm.Legal(byte(v)) {
				return c.JSON(http.StatusOK, H{"msg": "闲家1手牌，存在非法牌值"})
			}
		}
	}

	if len(ar.Seat[2]) != 13 {
		return c.JSON(http.StatusOK, H{"msg": "闲家2手牌数量不足"})
	} else {
		for _, v := range ar.Seat[2] {
			if !algorithm.Legal(byte(v)) {
				return c.JSON(http.StatusOK, H{"msg": "闲家2手牌，存在非法牌值"})
			}
		}
	}

	if len(ar.Seat[3]) != 13 {
		return c.JSON(http.StatusOK, H{"msg": "闲家3手牌数量不足"})
	} else {
		for _, v := range ar.Seat[3] {
			if !algorithm.Legal(byte(v)) {
				return c.JSON(http.StatusOK, H{"msg": "闲家3手牌，存在非法牌值"})
			}
		}
	}

	if len(ar.Card) != int(algorithm.TOTAL-algorithm.HAND*4-1) {
		return c.JSON(http.StatusOK, H{"msg": "剩余的牌的牌不足"})
	} else {
		for _, v := range ar.Card {
			if !algorithm.Legal(byte(v)) {
				return c.JSON(http.StatusOK, H{"msg": "剩余的牌的牌不足，存在非法牌值"})
			}
		}
	}

	if !algorithm.Legal(byte(ar.Wildcard)) {
		return c.JSON(http.StatusOK, H{"msg": "财神牌，非法牌值"})
	}

	r.SetCheat(ar.Seat, ar.Card, byte(ar.Wildcard))

	//if r.Len() >= 1 {
	// 邀请机器人进入房间
	glog.Infoln(" 邀请机器人进入房间", ar)
	go Client(strconv.Itoa(ar.RoomId), ar.RobotNum)
	//}

	glog.Infoln(" ", ar)
	return c.JSON(http.StatusOK, H{"msg": "数据提交成功"})
}

//召唤机器人
func Client(Code string, num int) {
	var addr string = "localhost:8085"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	var Key string = "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"
	var SIGN string = "qjby9vPheetlyYlsVjevzEltqh0b8b8FyESO+UqYPWc"
	var Now string = strconv.FormatInt(utils.Timestamp(), 10)
	//var Code string = "123456"
	var Num string = strconv.Itoa(num)
	var Sign string = utils.Md5(Key + Now + Code + Num)
	var Token string = Sign + Now + Code + Num
	c, _, err := websocket.DefaultDialer.Dial(u.String(),
		http.Header{"Token": {Token}})
	glog.Infoln(SIGN, Code, " : ", Num)
	if err != nil {
		glog.Errorf("dial err -> %v\n", err)
	}
	if c != nil {
		c.WriteMessage(websocket.TextMessage, []byte(SIGN+Code+Num))
		c.Close()
	}
}

//自动发牌，只需要填写部分的牌，如果全都没填写就走自动发牌流程
func autoCreate(c echo.Context) error {
	c.Response().CloseNotify()
	body, _ := ioutil.ReadAll(c.Request().Body)

	ar := &RoomCheat{}
	err := json.Unmarshal(body, ar)
	if err != nil {
		glog.Errorln("json.Unmarshal failed ", err)
		return c.JSON(http.StatusOK, H{"msg": err.Error()})
	}
	r := room.Get(strconv.Itoa(ar.RoomId))

	glog.Errorln("_______邀请机器人进入房间--------", r, ar.RoomId)
	if r == nil {
		return c.JSON(http.StatusOK, H{"msg": "房间不存在"})
	}

	if ar.RobotNum == 0 || ar.RobotNum > 3 {
		c.JSON(http.StatusOK, H{"msg": "邀请机器人的数量必须大于1 小于 3"})
	}
	go Client(strconv.Itoa(ar.RoomId), ar.RobotNum)
	return c.JSON(http.StatusOK, H{"msg": "成功邀请机器人"})
}
