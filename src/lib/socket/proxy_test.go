/**********************************************************
 * Author : Michael
 * Email : dolotech@163.com
 * Last modified : 2016-06-11 16:18
 * Filename : proxy_test.go
 * Description :
 * *******************************************************/
package socket

import (
	"protocol"
	"reflect"
	"testing"

	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
)

type handler struct {
	f interface{}
	t reflect.Type
}

func Binding(c uint32, s interface{}, f interface{}) {
	//func Binding(s interface{}, f interface{}) {

	if reflect.TypeOf(f).Kind() == reflect.Func {
		m[c] = &handler{f: f, t: reflect.TypeOf(s)}
		//m[s.(interfacer.IProto).GetCode()] = &handler{f: f, t: reflect.TypeOf(s)}
	} else {
		glog.Errorln("must be function")
	}
}

var m map[uint32]*handler = make(map[uint32]*handler)

func parse(c uint32, b []byte) {
	if h, ok := m[c]; ok {
		v := reflect.New(h.t)
		if err := proto.Unmarshal(b, v.Interface().(proto.Message)); err == nil {
			reflect.ValueOf(h.f).Call([]reflect.Value{v})
		} else {
			glog.Errorln("proto unmarshal fail")
		}
	} else {
		glog.Errorln("proto not regist")
	}
}

func init() {
	p := protocol.SResource{}
	Binding(p.GetCode(), p, handle)
	//Binding(&p, handle)
}
func handle(data *protocol.SResource) {
	glog.Errorln(data)
}
func Test_1(t *testing.T) {
	b, _ := proto.Marshal(&protocol.SResource{Id: proto.Uint32(456), Count: proto.Int32(123)})
	parse(5000, b)
}
