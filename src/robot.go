package main

import (
	"lib/utils"
	"bufio"
	//"data"
	"flag"
	"fmt"
	"robots"
	"strconv"
	"strings"
	"io"
	"net"
	"net/url"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"runtime/debug"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

func main() {
	//var config string
	//flag.StringVar(&config, "conf", "./conf.json", "config path")
	flag.Parse()
	//data.LoadConf(config)
	//glog.Infoln("Config: ", data.Conf)
	//serve_port := strconv.Itoa(data.Conf.Port)
	//robot_port := data.Conf.RobotPort
	defer glog.Flush()
	var serve_host string = "127.0.0.1"
	var serve_port string = "8005"
	var robot_addr string = "localhost:8085"
	robots.Start(serve_host, serve_port)
	ln, lnCh := server(robot_addr)
	//go client() //test
	//go testConf.runConfig("./conf.txt") //test
	robotsignalProc(ln, lnCh)
}

func server(addr string) (ln net.Listener, ch chan error) {
	ch = make(chan error)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}
	go func() {
		r := routes()
		ch <- http.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)}, r)
	}()
	return
}

func routes() (r *mux.Router) {
	r = mux.NewRouter()
	r.HandleFunc("/", wSHandler).Methods("GET")
	return
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  8,
	WriteBufferSize: 8,
}

var Key  string = "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"
var SIGN string = "qjby9vPheetlyYlsVjevzEltqh0b8b8FyESO+UqYPWc"

func wSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	var Token string = r.Header.Get("Token")
	//fmt.Printf("Token -> %s\n", Token)
	var TokenB []byte = []byte(Token)
	if !verifyToken(TokenB) {
		return
	}
	var CodeS string = string(TokenB[42:48])
	var NumS  string = string(TokenB[48:49])
	//fmt.Printf("CodeS -> %s, NumS -> %s\n", CodeS, NumS)
	socket, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	ipaddress := strings.Split(r.RemoteAddr, ":")[0]
	fmt.Printf("ipaddress -> %s\n", ipaddress)
	go func(){
		_, message, err := socket.ReadMessage() //TODO:timeout
		if err != nil {
			fmt.Printf("read err -> %v\n", err)
		}
		socket.Close()
		if len(message) == 50 {
			var msgS string = string(message[:43])
			var msgC string = string(message[43:49])
			var msgN string = string(message[49:50])
			if msgS == SIGN && msgC == CodeS && msgN == NumS {
				Num, _ := strconv.ParseInt(msgN, 10, 64)
				glog.Infoln("房间号：",msgC)
				robots.Msg2Robots(msgC, Num)
			} else {
				fmt.Printf("message err -> %s\n", string(message))
				fmt.Printf("CodeS -> %s, NumS -> %s\n", CodeS, NumS)
			}
		} else {
			fmt.Printf("message err -> %s\n", string(message))
		}
	}()
}

func client() {
	var addr string = "localhost:8085"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	var Now string = strconv.FormatInt(utils.Timestamp(), 10)
	var Code string = "123456"
	var Num string = "3"
	var Sign string = utils.Md5(Key+Now+Code+Num)
	var Token string = Sign+Now+Code+Num
	c, _, err := websocket.DefaultDialer.Dial(u.String(),
	http.Header{"Token":{Token}})
	if err != nil {
		fmt.Printf("dial err -> %v\n", err)
	}
	if c != nil {
		c.WriteMessage(websocket.TextMessage, []byte(SIGN+Code+Num))
		c.Close()
	}
}

func verifyToken(TokenB []byte) bool {
	// client
	// Key := "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"
	// Now := strconv.FormatInt(utils.Timestamp(), 10)
	// Sign := utils.Md5(Key+Now+Code+Num)
	// Token := Sign+Now+Code+Num
	// r.Header.Set("Token")
	// server
	// Key := "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"
	// Token := r.Header.Get("Token")
	// r.Header.Del("Token")
	// TokenB := []byte(Token)
	if len(TokenB) == 49 {
		//SignB := TokenB[:32]
		//TimeB := TokenB[32:42] TODO:过期验证
		if utils.Md5(Key+string(TokenB[32:])) == string(TokenB[:32]) {
			return true
		}
	}
	return false
}

func robotsignalProc(ln net.Listener, lnCh chan error) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(string(debug.Stack()))
		}
	}()
	ch := make(chan os.Signal, 1)
	//监听SIGINT和SIGKILL信号
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGHUP)
	glog.Infoln("signalProc ... ")
	for {
		msg := <-ch
		switch msg {
		default:
			//先关闭监听服务
			ln.Close()
			glog.Infoln(<-lnCh)
			//关闭服务
			robots.CloseRobots()
			//延迟退出，等待连接关闭，数据回存
			glog.Infof("get sig -> %v\n", msg)

			return
		case syscall.SIGHUP:
			glog.Infof("get sighup -> %v\n", msg)
		}
	}
}

// test
func testClient() {
	d := testGetData()
	if d == nil {
		fmt.Printf("get data err -> %v\n", 1)
		return
	}
	var addr string = d.Host + ":" + d.Port
	var Code string = d.Code
	var Num string = d.Num
	fmt.Printf("addr -> %s\n", addr)
	fmt.Printf("Code -> %s, Num -> %s\n", Code, Num)
	//
	//var addr string = "localhost:8085"
	//var Code string = "123456"
	//var Num string = "3"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/"}
	var Now string = strconv.FormatInt(utils.Timestamp(), 10)
	var Sign string = utils.Md5(Key+Now+Code+Num)
	var Token string = Sign+Now+Code+Num
	c, _, err := websocket.DefaultDialer.Dial(u.String(),
	http.Header{"Token":{Token}})
	if err != nil {
		fmt.Printf("dial err -> %v\n", err)
	}
	if c != nil {
		c.WriteMessage(websocket.TextMessage, []byte(SIGN+Code+Num))
		c.Close()
	}
}

// read conf.txt
//配置数据结构
type testData struct {
	Host    string `json:"host"`
	Port    string `json:"port"`
	State   string `json:"state"`
	Message string `json:"message"`
	Code    string `json:"code"`
	Num     string `json:"num"`
}

//回调消息
type testMsg struct {
	callback chan *testData
}

//配置结构
type Config struct {
	filepath string //配置文件路径
	conflist map[string]map[string]string
	modtime time.Time // 文件修改时间
	message chan *testMsg
}

//全局变量(外部调用)
var testConf = Config {
	message: make(chan *testMsg, 1024),
}

//获取数据
func testGetData() *testData {
	d := make(chan *testData)
	testConf.message <- &testMsg{callback: d}
	return <-d
}

//启动读取配置服务
func (this *Config) runConfig(filepath string) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(string(debug.Stack()))
		}
	}()

	this.filepath = filepath
	this.conflist = make(map[string]map[string]string)
	this.readConfig()

	for {
		select {
		case testMsg := <-this.message:
			// TODO:检测服务器状态然后分配
			d := &testData{}
			d.Host    = this.conflist["server1"]["host"]
			d.Port    = this.conflist["server1"]["port"]
			d.State   = this.conflist["server1"]["state"]
			d.Message = this.conflist["server1"]["message"]
			d.Code    = this.conflist["server1"]["code"]
			d.Num     = this.conflist["server1"]["num"]
			testMsg.callback <- d
		case <-time.After(10 * time.Second):
			var ok bool = this.readConfig()
			if ok { //文件有修改
				go testClient() //分配机器人
			}
		}
	}
}

func (this *Config) readConfig() bool {
	file, err := os.Open(this.filepath)
	if err != nil {
		glog.Errorln("open file err:", err)
	}
	fileinfo, err := file.Stat()
	if err != nil {
		glog.Errorln("fileinfo err:", err)
	}
	defer file.Close()
	if this.modtime != fileinfo.ModTime() {
		this.modtime = fileinfo.ModTime()
		var section string
		buf := bufio.NewReader(file)
		for {
			l, err := buf.ReadString('\n')
			line := strings.TrimSpace(l)
			len := len(line)
			if err != nil && err != io.EOF {
				glog.Errorln("read file err:", err)
			}
			if err == io.EOF {
				break
			}
			if len == 0 {
				continue
			}
			switch {
			case line[0] == '[' && line[len-1] == ']':
				section = strings.TrimSpace(line[1:len-1])
				this.conflist[section] = make(map[string]string)
			default:
				i := strings.IndexAny(line, "=")
				key := strings.TrimSpace(line[0:i])
				value := strings.TrimSpace(line[i+1 : len])
				this.conflist[section][key] = value
				//TODO:map concurrent read and write
			}
		}
		return true
	}
	return false
}
