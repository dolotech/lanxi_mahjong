package robots

import (
	"lib/utils"
	"lib/socket"
	"game/interfacer"
	"net/url"
	"reflect"
	"time"

	"github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
	"github.com/gorilla/websocket"
	"runtime/debug"
	"protocol"
)

// 机器人连接数据
type Robot struct {
	readCh  chan *socket.Packet
	conn    *websocket.Conn
	data    *user            //数据
	code    string           //邀请码
	index   uint32           //包序
	seat    uint32           //位置
	cards   []byte           //手牌
	pongCards []byte	// 碰的牌
	msgCh   chan interface{} //
	closeCh chan bool        //
	room *protocol.RoomData
}

// 基本数据
type user struct {
	Userid   string // 用户id
	Nickname string // 用户昵称
	Sex      uint32 // 用户性别,男1 女2 非男非女3
	Phone    string // 绑定的手机号码
	Coin     uint32 // 金币
	Exp      uint32 // 经验
	Diamond  uint32 // 钻石
	Ticket   uint32 // 入场券
	Exchange uint32 // 兑换券
	Vip      uint32 // vip
	RoomCard uint32 // 房卡
}

type handler struct {
	f interface{}
	t reflect.Type
}

func regist(s interface{}, f interface{}) {
	msg := s.(interfacer.IProto)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {

		v = v.Elem()
		s = v.Interface()
		v = reflect.ValueOf(s)
	}

	if reflect.TypeOf(f).Kind() == reflect.Func {
		m[msg.GetCode()] = &handler{f: f, t: reflect.TypeOf(s)}
	} else {
		glog.Errorln("must be function")
	}
}

var m map[uint32]*handler = make(map[uint32]*handler)

//启动一个机器人
func RunRobot(host, port, phone, code string, msgCh chan interface{}) {
	glog.Infof("run robot -> %s", phone,code)
	u := url.URL{Scheme: "ws", Host: host + ":" + port, Path: "/"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		glog.Errorf("robot run dial -> %v", err)
		return
	}
	this := &Robot{
		readCh:  make(chan *socket.Packet),
		data:    &user{Phone: phone},
		msgCh:   msgCh,
		code:    code,
		conn:    c,
		closeCh: make(chan bool, 1),
	}
	rl := GetRobotList()
	var r_l int32 = int32(len(rl))
	var num int32 = utils.RandInt32N(r_l)
	r_d := rl[num]
	this.data.Nickname = r_d.Nickname
	this.data.Sex = 2
	go this.RecvMsg() //接收消息
	go this.Reader()  //协议路由
	this.SendRegist() //注册
	this.ticker()     //计时器
}

// 关闭连接
func (this *Robot) Close() {
	this.conn.Close()
}

// 关闭连接
func (this *Robot) Closed() {
	close(this.closeCh)
	close(this.readCh)
	this.msgCh <- Logout{phone: this.data.Phone}
}

// 关闭连接
func (this *Robot) Logined() {
	this.msgCh <- Login{phone: this.data.Phone}
}

//计时器
func (this *Robot) ticker() {
	tick := time.Tick(time.Minute)
	glog.Infof("ticker -> %s", this.data.Phone)
	for {
		select {
		case <-tick:
		//逻辑处理
		//TODO:不在游戏中时下线或者重新开始,
		//this.Close()
		case <-this.closeCh:
			glog.Infof("ticker closed -> %s", this.data.Phone)
			return
		}
	}
}

// 接收消息
func (this *Robot) RecvMsg() {
	defer func() {
		glog.Infof("RecvMsg closed -> %s", this.data.Phone)
		this.Close()
		this.Closed()
	}()

	tmpBuffer := make([]byte, 10*1024)
	var length uint32
	for {
		_, message, err :=this.conn.ReadMessage()
		if err != nil {
			return
		}
		copy(tmpBuffer[length:],message)
		length +=uint32( len(message))
		readCount := socket.Unpack(tmpBuffer,length, this.readCh)
		reminder:= length - readCount
		if reminder <0{
			reminder = 0
		}
		if length - readCount > 0{
			copy(tmpBuffer,tmpBuffer[readCount:reminder])
		}
		length = reminder
	}
}

// 协议路由
func (this *Robot) Reader() {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln("reader err -> ", err)
		}
		glog.Infof("reader closed -> %s", this.data.Phone)
	}()
	for {
		select {
		// 如果管道关闭则退出for循环，因为管道关闭不会阻塞导致for进入死循环
		case packet, ok := <-this.readCh: //接受消息
			if !ok {
				glog.Infof("reader closed -> %s", this.data.Phone)
				return
			}
			this.proxyHandle(packet.GetProto(), packet.GetContent())
		}
	}
}

func (this *Robot) proxyHandle(c uint32, b []byte) {

	defer func() {
		if e := recover(); e != nil {
			glog.Errorln(c, string(debug.Stack()))
		}
	}()

	if h, ok := m[c]; ok {
		v := reflect.New(h.t)
		if err := proto.Unmarshal(b, v.Interface().(proto.Message)); err == nil {
			reflect.ValueOf(h.f).Call([]reflect.Value{v, reflect.ValueOf(this)})
		} else {
			glog.Errorln("protocol  unmarshal fail: ", c)
		}
	} else {
		//glog.Errorln("protocol not regist:", c)
	}
}

//发送请求
//TODO:concurrent write to websocket connection
func (this *Robot) Sender(packet interfacer.IProto) {
	m, _ := proto.Marshal((proto.Message)(packet))
	this.index++
	this.index = this.index % 256
	//glog.Infoln("index -> ", this.index)
	msg := socket.Pack(packet.GetCode(), m, this.index)
	this.conn.WriteMessage(websocket.TextMessage, msg)
}
