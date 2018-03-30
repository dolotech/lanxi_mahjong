/**********************************************************
 * Author : Michael
 * Email : dolotech@163.com
 * Last modified : 2016-06-11 16:27
 * Filename : interface.go
 * Description :  零散的接口
 * *******************************************************/
package interfacer


type IProto interface {
	GetCode() uint32
	Reset()
	String() string
	ProtoMessage()
}
