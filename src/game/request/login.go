package request

import (
	"game/data"
	"game/room"
	"lib/event"
	"lib/socket"
	"lib/utils"

	"game/interfacer"
	"game/players"
	"protocol"
	"runtime/debug"
	//"strconv"
	"time"

	"code.google.com/p/goprotobuf/proto"
	"github.com/golang/glog"
)

func init() {

	socket.Regist(&protocol.CLogin{}, login)
}

func login(ctos *protocol.CLogin, c interfacer.IConn) {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(string(debug.Stack()))
		}

	}()
	//glog.Infoln(ctos.Userid, ctos.GetPassword())
	stoc := &protocol.SLogin{}
	var member *data.User
	if ctos.GetPhone() != "" && utils.PhoneRegexp(ctos.GetPhone()) {
		member = &data.User{Phone: ctos.GetPhone()}
		if ctos.GetPassword() != "" && len(ctos.GetPassword()) == 32 && member.VerifyPwdByPhone(ctos.GetPassword()) {
			c.SetUserid(member.Userid)
		} else {
			//glog.Errorln(member.Userid)
			stoc.Error = proto.Uint32(uint32(protocol.Error_UsernameOrPwdError))
		}

	} else {
		stoc.Error = proto.Uint32(uint32(protocol.Error_UsernameOrPwdError))
	}
	if member != nil {
		err := member.Get()
		if err != nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_UsernameOrPwdError))
		}
	}

	if stoc.Error != nil {
		c.Send(stoc)
		time.AfterFunc(time.Millisecond*200, c.Close)
		return
	}
	//重复登录检测
	logining(member, c)
	stoc.Userid = &member.Userid
	c.Send(stoc)

}

func logining(member *data.User, c interfacer.IConn) {
	//	重复登录检测
	glog.Infoln("玩家ID : ", c.GetUserid())
	// 已经在房间打牌
	userdata := players.Get(member.Userid)
	if userdata != nil {
		conn := userdata.GetConn()
		conn.(event.IDispath).RemoveAll()
		conn.Close()
	} else {
		// 登陆成功把用户的数据从数据库取出存入服务内存
		userdata = room.NewPlayer(member)
		players.Set(userdata.GetUserid(), userdata)
	}
	userdata.SetConn(c)

	c.SetLogin()

	rdata := room.Get(userdata.GetInviteCode())
	if rdata != nil {
		rdata.Offline(userdata.GetSeat(), false)
	}
	//注释掉下面玩家下线处理，玩家数据会越积越多，造成内存泄漏
	c.(event.IDispath).ListenOnce(socket.OFFLINE, func(t string, args interface{}) {
		active := &data.DataUserActive{Userid: member.Userid, IP: c.GetIPAddr()}
		active.Logout()

		if rdata != nil {
			rdata.Offline(userdata.GetSeat(), true)
		}

	})
	// 记录登陆时间和IP地址
	go func() {
		active := &data.DataUserActive{Userid: member.Userid, IP: c.GetIPAddr()}
		active.Login()
		//tradeOff(userdata) //发货失败订单检测
	}()

}
