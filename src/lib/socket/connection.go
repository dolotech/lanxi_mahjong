package socket

import (
	"lib/event"
	"game/interfacer"
	"runtime/debug"
	"time"

	"github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
	"github.com/gorilla/websocket"
)

const (
	// 网络掉线事件
	OFFLINE = "offline"

	startBufSize   int64 = 512
	minBufSize     int64 = 128
	maxBufSize     int64 = 65536
)

func newConnection(socket *websocket.Conn,ip uint32) *Connection {
	c := &Connection{
		writeChan: make(chan interfacer.IProto, 128),
		ws:        socket,
		ReadChan:  make(chan *Packet, 128),
		connected: true,
		closeChan: make(chan bool, 1),
		ipAddr:ip,
	}
	return c
}

type Connection struct {
	writeChan chan interfacer.IProto
	userid    string // 玩家ID
	logined   bool   // true 标示已登录
	connected bool   // false标示连接断开
	ws        *websocket.Conn
	ReadChan  chan *Packet
	closeChan chan bool
	ipAddr    uint32 // 当前连得IP地址
	event.Dispatcher // 事件管理器
	count     uint32
}

func (c *Connection) GetConnected() bool {
	return c.connected
}
func (c *Connection) GetIPAddr() uint32 {
	return c.ipAddr
}
func (c *Connection) SetLogin() {
	c.logined = true
}
func (c *Connection) GetLogin() bool {
	return c.logined
}

func (c *Connection) SetUserid(userid string) {
	c.userid = userid
}
func (c *Connection) GetUserid() string {
	return c.userid
}
func (c *Connection) Close() {
	c.ws.Close()
}

func (c *Connection) LoginTimeout() {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(string(debug.Stack()))
		}
	}()
	//建立连接后一定时间没有登录断开连接
	select {
	case <-time.After(waitForLogin):
		if !c.logined {
			c.Close()
		}
	}
}

func (c *Connection) Reader(readChan chan *Packet) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(string(debug.Stack()))
		}
	}()

	for {
		select {
		// 如果管道关闭则退出for循环，因为管道关闭不会阻塞导致for进入死循环
		case packet, ok := <-readChan:
			if !ok {
				return
			}
			//glog.Infoln(packet.GetProto())
			c.count++
			c.count = c.count % 256
			if c.count != packet.count {
				glog.Errorln("count error -> ", c.count, packet.count)
				return
			}
			proxyHandle(packet.GetProto(), packet.GetContent(), c)
		}
	}
}

func (c *Connection) ReadPump() {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(string(debug.Stack()))
		}
	}()

	defer func() {
		//c.ws.Close()
		c.Close() //
		c.connected = false
		close(c.writeChan)
		close(c.ReadChan)
		close(c.closeChan)
		logout(c)
		if c.logined {
			c.Dispatch(OFFLINE, c)
		}
	}()
	c.ws.SetReadLimit(maxBufSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(
		func(string) error {
			c.ws.SetReadDeadline(time.Now().Add(pongWait))
			return nil
		})
	//声明一个临时缓冲区，用来存储被截断的数据
	tmpBuffer := make([]byte, startBufSize)
	var length uint32

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			return
		}
		copy(tmpBuffer[length:],message)
		length +=uint32( len(message))
		readCount := Unpack(tmpBuffer,length, c.ReadChan)

		reminder:= length - readCount
		if reminder <0{
			reminder = 0
		}
		if reminder > 0{
			if len(tmpBuffer) < int(maxBufSize){
				buf:= make([]byte,len(tmpBuffer)*2)
				copy(buf,tmpBuffer[readCount:reminder])
				tmpBuffer = buf
			}else{
				copy(tmpBuffer,tmpBuffer[readCount:reminder])
			}
		}else{
			if len(message) <len(tmpBuffer)/2 && len(tmpBuffer) > int(minBufSize) {
				tmpBuffer = make([]byte, len(tmpBuffer)/2)
			}
		}
		length = reminder
	}
}

func (c *Connection) Send(data interfacer.IProto) {
	defer func() {
		if err := recover(); err != nil {
			c.Close() //
			glog.Errorln(string(debug.Stack()))
		}
	}()
	if c.connected {
		c.writeChan <- data
	} else {
	}
}
func (c *Connection) write(mt int, packet interfacer.IProto) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	msg, _ := proto.Marshal((proto.Message)(packet))
	if len(msg) > 0 {
		b := Pack(packet.GetCode(), msg, 0)
		return c.ws.WriteMessage(mt, b)
	} else {
		return c.ws.WriteMessage(mt, msg)
	}
}

func (c *Connection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		//c.ws.Close()
		c.Close() //
		if err := recover(); err != nil {
			glog.Errorln(string(debug.Stack()))
		}
	}()

	for {
		select {
		// 如果管道关闭则退出for循环，因为管道关闭不会阻塞导致for进入死循环
		case proto, ok := <-c.writeChan:
			if !ok {
				c.write(websocket.CloseMessage, nil)
				return
			}
			if err := c.write(websocket.TextMessage, proto); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
