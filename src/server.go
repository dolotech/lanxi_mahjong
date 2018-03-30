package main

import (
	"flag"
	"game/data"
	//_ "net/http/pprof"
	"cheat"
	"config"
	"fmt"
	"game/csv"
	_ "game/request"
	"game/room"
	_ "game/roomrequest"
	"github.com/golang/glog"
	"lib/socket"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
)

var (
	VERSION    = "0.0.0"
	BUILD_TIME = ""
)

func main() {
	fmt.Println("version: ", VERSION, "timestamp:", BUILD_TIME)
	var path string
	flag.StringVar(&path, "conf", "./config.toml", "config path")
	flag.Parse()
	config.ParseToml(path)

	socket.VERSION = VERSION
	socket.BUILD_TIME = BUILD_TIME
	glog.Infoln("Config: ", config.Opts())
	defer glog.Flush()
	glog.Infoln("逻辑服务器端口:", config.Opts().Server_port)

	data.InitIDGen()
	csv.InitShop()
	//csv.InitPattern()

	ln := socket.Server(config.Opts().Server_port)

	glog.Infoln("Server listening on", config.Opts().Server_port)
	//glog.Infoln("Server started at", ln.Addr())
	go cheat.Run(config.Opts().AdminPort)
	//go statistics.InitjobSechduler(5)
	gamesignalProc(ln)
}

func pprof() {
	if config.Opts().Oprof_port != "" {
		err := http.ListenAndServe(config.Opts().Oprof_port, nil)
		glog.Infoln("性能监控端口:", config.Opts().Oprof_port)
		if err != nil {
			glog.Fatal("ListenAndServe error: ", err)
		}
	}
}

func gamesignalProc(ln net.Listener) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(string(debug.Stack()))
		}
	}()
	ch := make(chan os.Signal, 1)
	//signal.Notify(ch, syscall.SIGUSR1, syscall.SIGUSR2)
	//signal.Notify(ch, syscall.SIGHUP)
	signal.Notify(ch, os.Interrupt, os.Kill, syscall.SIGHUP) //监听SIGINT和SIGKILL信号
	glog.Infoln("signalProc ... ")
	//for {
	msg := <-ch
	//switch msg {
	//default:
	//先关闭监听服务
	//	ch.Close()
	//关闭连接
	socket.Close()
	//关闭服务
	room.Close()
	//players.Close()
	//延迟退出，等待连接关闭，数据回存
	glog.Infof("get sig -> %v\n", msg)

	return
	//case syscall.SIGHUP:
	//
	//}
	//}
}
