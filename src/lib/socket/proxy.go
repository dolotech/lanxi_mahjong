package socket

import (
	"reflect"
	"runtime/debug"

	"github.com/golang/glog"
	"code.google.com/p/goprotobuf/proto"
	"game/interfacer"
)

type handler struct {
	f interface{}
	t reflect.Type
}

func Regist(s interface{}, f interface{}) {
	msg:= s.(interfacer.IProto)
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {

		v = v.Elem()
		s = v.Interface()
		v = reflect.ValueOf(s)
	}



	if reflect.TypeOf(f).Kind() == reflect.Func {
		m[msg.GetCode()] = &handler{f: f, t: reflect.TypeOf(s)}
	} else {
		glog.Errorln("must be function")
	}
}

var m map[uint32]*handler = make(map[uint32]*handler)

func proxyHandle(c uint32, b []byte, conn *Connection) {
	defer func() {
		if e := recover(); e != nil {
			glog.Errorln(c, string(debug.Stack()))
		}
	}()

	if h, ok := m[c]; ok && (conn.GetLogin() || c == 1000 || c == 1022  || c==1026) {
		v := reflect.New(h.t)
		if err := proto.Unmarshal(b, v.Interface().(proto.Message)); err == nil {
			reflect.ValueOf(h.f).Call([]reflect.Value{v, reflect.ValueOf(conn)})
		} else {
			glog.Errorln("protocol  unmarshal fail: ", c)
		}
	} else {
		glog.Errorln("protocol not regist:", c)
	}
}
