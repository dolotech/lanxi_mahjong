package interfacer


type IProto interface {
	GetCode() uint32
	Reset()
	String() string
	ProtoMessage()
}
