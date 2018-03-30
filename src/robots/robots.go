package robots

import (
	"lib/utils"
	"time"

	"github.com/golang/glog"
)

//通知消息体
type Message struct {
	code string
}

type Login struct {
	phone string
}

type Logout struct {
	phone string
}

//机器人服务
type Robots struct {
	host    string
	port    string
	phone   string //注册起始电话号码,10005000101
	online  map[string]bool //map[phone]状态,true=在线,
	offline map[string]bool //map[phone]状态,true=离线,false=登录中
	msgCh   chan interface{}
	closeCh chan bool
}

//通道
var RobotsCh chan interface{}

//消息通知
func Msg2Robots(code string, num int64) {
	for num > 0 {
		RobotsCh <- Message{code: code}
		num--
	}
}

//关闭
func CloseRobots() {
	RobotsCh <- true
}

//启动
func Start(host, port string) {
	r := &Robots{
		host:    host,
		port:    port,
		phone:   "10005008101",
		online:  make(map[string]bool),
		offline: make(map[string]bool),
		msgCh:   make(chan interface{}, 100),
		closeCh: make(chan bool, 1),
	}
	RobotsCh = r.msgCh //通道
	go r.Run() //启动
	//test
	//r.msgCh <- Message{}
	//r.msgCh <- Message{}
	//r.msgCh <- Message{}
	//test
	//go r.runTest()
}

//机器人测试
func (r *Robots) runTest() {
	glog.Infof("runTest started -> %d", 1)
	tick := time.Tick(time.Second)
	for {
		select {
		case <-tick:
			glog.Infof("r.online -> %d\n", len(r.online))
			glog.Infof("r.offline -> %d\n", len(r.offline))
			glog.Infof("r.phone -> %s\n", r.phone)
			//TODO:优化
			//运行指定数量机器人(每个创建一个牌局)
			//code = "create" 表示机器人创建房间
			if len(r.online) < 5000 {
				go Msg2Robots("create", 10)
			}
		}
	}
}

//处理
func (r *Robots) Run() {
	defer func() {
		glog.Infof("Robots closed -> %d", 1)
	}()
	glog.Infof("Robots started -> %d", 1)
	tick := time.Tick(time.Minute)
	for {
		select {
		case m := <-r.msgCh:
			switch m.(type) {
			case Message:
				msg := m.(Message)
				var code string = msg.code
				var phone string
				//for k, v := range r.offline {
				//	if v {
				//		phone = k
				//		r.offline[k] = false
				//		break
				//	}
				//}
				//if len(phone) == 0 {
				//	phone = r.phone
				//	r.phone = utils.StringAdd(r.phone)
				//}
				phone = r.phone
				r.phone = utils.StringAdd(r.phone)
				go RunRobot(r.host, r.port, phone, code, r.msgCh)
				glog.Infof("phone -> %s", phone)
			case Login:
				msg := m.(Login)
				glog.Infof("login -> %v", msg)
				delete(r.offline, msg.phone)
				r.online[msg.phone] = true
			case Logout:
				msg := m.(Logout)
				glog.Infof("logout -> %v", msg)
				delete(r.online, msg.phone)
				r.offline[msg.phone] = true
			default:
				glog.Infof("r.online -> %d\n", len(r.online))
				glog.Infof("r.offline -> %d\n", len(r.offline))
				glog.Infof("r.phone -> %s\n", r.phone)
				glog.Errorf("Robots err -> %v", m)
				close(r.closeCh)
			}
		case <-tick:
			//逻辑处理,TODO:
		case <-r.closeCh:
			glog.Infof("Robots closed -> %d", 1)
			//TODO:关闭机器人
			return
		}
	}
}
