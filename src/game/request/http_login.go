package request

import (
	"code.google.com/p/goprotobuf/proto"
	"game/data"
	"game/interfacer"
	"lib/socket"
	"lib/utils"
	"protocol"
	"strconv"
)

func init() {
	socket.Regist(&protocol.CHttpLogin{}, httpLogin)
}

const KEY = "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"

//sign ==  "XG0e2Ye/KAUJRXaMNnJ5UH1haBvh2FXOoAggE6f2Utw"
//userid == 用户ID
//timestamp == 当前时间戳
//createtime == 注册时间戳
//md5(sign + userid + timestamp+ createtime)

func httpLogin(ctos *protocol.CHttpLogin, c interfacer.IConn) {
	stoc := &protocol.SHttpLogin{
		Error:  proto.Uint32(0),
		Userid: ctos.Userid,
	}

	if ctos.GetUserid() == "" {
		stoc.Error = proto.Uint32(uint32(protocol.Error_HTTP_LOGIN_USERID_NIL))
		c.Send(stoc)
		return
	}
	timestamp := strconv.FormatUint(uint64(ctos.GetTimestamp()), 10)
	createtime := strconv.FormatUint(uint64(ctos.GetCreatetime()), 10)

	if utils.Md5(KEY+ctos.GetUserid()+timestamp+createtime) != ctos.GetToken() {
		stoc.Error = proto.Uint32(uint32(protocol.Error_HTTP_LOGIN_TOKEN_FAIL))
		c.Send(stoc)
		return
	}

	member := &data.User{
		Userid: ctos.GetUserid(),
	}
	if member.Get() != nil {
		stoc.Error = proto.Uint32(uint32(protocol.Error_HTTP_LOGIN_USER_NOT_REGIST))
		c.Send(stoc)
		return
	}

	c.SetUserid(member.Userid)

	logining(member, c)
	c.Send(stoc)
}
