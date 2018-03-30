package interfacer

type IDesk interface {
	Discard(uint32, byte, bool) int32
	Operate(uint32, int64, uint32) int32
	Readying(uint32, bool) int32
	Enter(IPlayer) int32
	Vote(bool, uint32, uint32) int32
	Broadcasts(IProto)
	Closed(bool)
	SetCheat([][]byte, []byte, byte)
	MaiZi(seat, count uint32) int32
	Offline(uint32, bool)
	ToString() string
}
