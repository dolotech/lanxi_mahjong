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
