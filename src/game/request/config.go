package request

import (
	"lib/socket"
	"game/interfacer"
	"protocol"
	log "github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
	"config"
)

func init() {
	socket.Regist(&protocol.CConfig{}, getconfig)
}

func getconfig(ctos *protocol.CConfig, c interfacer.IConn) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorln(err)
		}
	}()
	stoc := &protocol.SConfig{
		Sys: &protocol.SysConfig{
			Discardtimeout: proto.Uint32(uint32(config.Opts().Operate_timeout)),
			Version:        proto.String(config.Opts().Version),
			Shareaddr:      proto.String(config.Opts().Share_addr),
		},
	}

	c.Send(stoc)
}
