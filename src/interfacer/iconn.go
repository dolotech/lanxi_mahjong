/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 10:43
 * Filename      : iconn.go
 * Description   : 每个玩家socket连接的对象接口
 * *******************************************************/
package interfacer

type IConn interface {
	Close()
	Send(IProto)
	GetUserid() string
	SetUserid(string)
	GetLogin() bool
	SetLogin()
	GetIPAddr() uint32
	GetConnected() bool
}
