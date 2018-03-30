package request

import (
	"lib/socket"
	"game/interfacer"
	"protocol"

	"code.google.com/p/goprotobuf/proto"
)

func init() {

	socket.Regist(&protocol.CPing{}, ping)
}

func ping(ctos *protocol.CPing, c interfacer.IConn) {
	stoc := &protocol.SPing{Error: proto.Uint32(0)}
	c.Send(stoc)
}
