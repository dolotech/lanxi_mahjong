/**********************************************************
 * Author        : Michael
 * Email         : dolotech@163.com
 * Last modified : 2016-01-23 10:42
 * Filename      : iroom.go
 * Description   : 房间的数据接口
 * *******************************************************/
package interfacer

type IDesk interface {
	Discard(uint32, byte, bool)int32
	Operate(uint32, int64, uint32) int32
	Readying(uint32, bool) int32
	Enter(IPlayer) int32
	Vote(bool, uint32, uint32) int32
	Broadcasts(IProto)
	Closed(bool)
	SetCheat( [][]byte,[]byte,byte)
	MaiZi(seat, count uint32) int32
	Offline( uint32, bool)

	ToString ()string
}
