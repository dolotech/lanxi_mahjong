package algorithm


// 碰杠吃数据
func EncodePeng(seat uint32, card byte) uint32 {
	seat = seat << 8
	seat |= uint32(card)
	return seat
}

func DecodePeng(value uint32) (seat uint32, card byte) {
	seat = value >> 8
	card = byte(value & 0xFF)
	return
}

func EncodeKong(seat uint32, card byte, value uint32) uint32 {
	value = value << 16
	value |= (seat << 8)
	value |= uint32(card)
	return value
}

func DecodeKong(value uint32) (seat uint32, card byte, v uint32) {
	v = value >> 16
	seat = (value >> 8) & 0xFF
	card = byte(value & 0xFF)
	return
}

func EncodeChow(c1, c2, c3 byte) (value uint32) {
	value =  uint32(c1) << 16
	value |= uint32(c2) << 8
	value |= uint32(c3)
	return
}

func DecodeChow(value uint32) (c1, c2, c3 byte) {
	c1 = byte(value >> 16)
	c2 = byte(value >> 8 & 0xFF)
	c3 = byte(value & 0xFF)
	return
}

func DecodeChow2(value uint32) (c1, c2 byte) {
	c1 = byte(value >> 8)
	c2 = byte(value & 0xFF)
	return
}
