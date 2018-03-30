package algorithm

func replaceWildcard(cs []byte, wildcard byte, cp bool) []byte {
	if cp {
		list := make([]byte, len(cs))
		copy(list, cs)
		cs = list
	}
	// 有鬼牌
	if wildcard > 0 {
		for k, v := range cs {
			if v == wildcard {
				cs[k] = WILDCARD
			}
		}
	}
	return cs
}

// 打牌检测,明杠/碰牌
func DiscardPong(card byte, cs []byte, wildcard byte) (status int64) {
	//碰杠
	if existMingKong(card, cs, wildcard) {
		status |= MING_KONG
		status |= KONG
	}
	if existPeng(card, cs, wildcard) {
		status |= PENG
	}
	return status
}

// 打牌检测,吃
func DiscardChow(s1, s2 uint32, card byte, cs []byte, wildcard byte) (status int64) {
	if NextSeat(s1) == s2 {
		list := existChow(card, cs, wildcard)
		if len(list) >= 2 {
			status |= CHOW
		}
	}
	return
}

// 正常流程走牌令牌移到下一家
func NextSeat(seat uint32) uint32 {
	if seat == 4 {
		return 1
	}
	return seat + 1
}

// 是否存在n个牌
func Exist(c byte, cs []byte, n int) bool {
	for _, v := range cs {
		if n == 0 {
			return true
		}
		if c == v {
			n--
		}
	}
	return n == 0
}

// 移除一个牌
func Remove(c byte, cs []byte) []byte {
	for i, v := range cs {
		if c == v {
			cs = append(cs[:i], cs[i+1:]...)
			break
		}
	}
	return cs
}

// 移除n个牌
func RemoveN(c byte, cs []byte, n int) []byte {
	for n > 0 {
		for i, v := range cs {
			if c == v {
				cs = append(cs[:i], cs[i+1:]...)
				break
			}
		}
		n--
	}
	return cs
}
