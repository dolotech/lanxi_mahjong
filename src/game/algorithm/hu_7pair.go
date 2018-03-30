package algorithm

// 判断是否小七对(7个对子)
func exist7pair(cs []byte) int64 {
	if len(cs) != 14 {
		return 0
	}
	le := uint8(len(cs))
	var flag int16 = 0x3FFF

	// 先找出非万能的对子
	for i := uint8(0); i < le-1; i++ {
		for j := i + 1; j < le && ((flag>>i)&0x1) == 1; j++ {
			if ((flag>>j)&0x1) == 1 && cs[i] == cs[j] && cs[i] != WILDCARD {
				flag = flag & (^(1 << i))
				flag = flag & (^(1 << j))
				break
			}
		}
	}

	// 单个的杂牌跟万能匹对
	for i := uint8(0); i < le-1; i++ {
		for j := uint8(0); j < le && ((flag>>i)&0x1) == 1; j++ {
			if ((flag>>j)&0x1) == 1 && cs[j] != WILDCARD && cs[i] == WILDCARD {
				flag = flag & (^(1 << i))
				flag = flag & (^(1 << j))
				break
			}
		}
	}

	wilecardCount := 0
	dirtyCount := 0
	for i := uint8(0); i < le; i++ {
		if ((flag >> i) & 0x1) == 1 {
			if cs[i] == WILDCARD {
				wilecardCount ++
			} else {
				dirtyCount ++
			}
		}

	}

	if wilecardCount == dirtyCount || dirtyCount == 0 {

		return HU_SEVEN_PAIR
	}
	return 0
}
