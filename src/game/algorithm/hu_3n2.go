package algorithm


 //检测是否满足(N * ABC + M *DDD)
func check3n(cs []byte, flag *int32) bool {
	if *flag == 0 {
		return true
	}
	flag_temp := *flag
	if removeABC(cs, flag) {
		if *flag == 0 || check3n(cs, flag) {
			return true
		}
	}
	flag = &flag_temp
	if removeDDD(cs, flag) {
		if *flag == 0 || check3n(cs, flag) {
			return true
		}
	}
	return false
}

//移除将牌后检测是否满足3n
func removeEE(cs []byte, flag int32, i uint8, j uint8) bool {
	flag = flag & (^(1 << i))
	flag = flag & (^(1 << j))
	return check3n(cs, &flag)
}

//尝试移除一个顺子
func removeABC(cs []byte, flag *int32) bool {
	le := uint8(len(cs))
	var first, second, third, wildcard_i, wildcard_j uint8
	firstFlag := true
	for i := uint8(0); i < le; i++ {
		if (((*flag) >> i) & 0x1) == 1 {
			if cs[i]>>4 < FENG { // 字牌不做顺子
				if firstFlag {
					first = i
					firstFlag = false
				} else if cs[i] == cs[first]+1 {
					second = i
				} else if cs[i] == cs[first]+2 {
					third = i
				}
				if second > 0 && third > 0 {
					*flag = (*flag) & (^(1 << first))
					*flag = (*flag) & (^(1 << second))
					*flag = (*flag) & (^(1 << third))
					return true
				}
			}

			if cs[i] == WILDCARD && wildcard_i == 0 {
				wildcard_i = i
			} else if cs[i] == WILDCARD && wildcard_j == 0 {
				wildcard_j = i
			}
		}
	}
	if wildcard_i != 0 {
		if second > 0 { //两张顺缺一张
			*flag = (*flag) & (^(1 << first))
			*flag = (*flag) & (^(1 << second))
			*flag = (*flag) & (^(1 << wildcard_i))
			return true
		} else if third > 0 { //卡窿缺一张
			*flag = (*flag) & (^(1 << first))
			*flag = (*flag) & (^(1 << third))
			*flag = (*flag) & (^(1 << wildcard_i))
			return true
		} else if wildcard_j != 0 && firstFlag { //缺两张
			*flag = (*flag) & (^(1 << first))
			*flag = (*flag) & (^(1 << wildcard_i))
			*flag = (*flag) & (^(1 << wildcard_j))
			return true
		}
	}
	return false
}

func removeDDD(cs []byte, flag *int32) bool {
	var count int8
	var a uint8 = 0xFF
	le := uint8(len(cs))
	for i := uint8(0); i < le; i++ {
		if (((*flag) >> i) & 0x1) == 1 {
			if a == 0xFF {
				a = i
				*flag = (*flag) & (^(1 << i))
				continue
			}
			if cs[a] == cs[i] || cs[i] == WILDCARD {
				count++
				*flag = (*flag) & (^(1 << i))
			}

			if count == 2 {
				return true
			}
		}

	}
	return false
}

func check3nPengPeng(cs []byte, flag *int32) bool {
	if *flag == 0 {
		return true
	}
	if removeDDD(cs, flag) {
		if *flag == 0 {
			return true
		}
		return check3nPengPeng(cs, flag)
	}
	return false
}

func hu3n2(cs []byte) bool {
	le := uint8(len(cs))
	var flag int32
	for i := uint8(0); i < le; i++ {
		flag = flag | (1 << i)
	}
	var i, j uint8
	for i = 0; i < le-1; i++ {
		j = i + 1
		if cs[i] == cs[j] { // 两张自身牌或两张万能牌
			if removeEE(cs, flag, i, j) {
				return true
			}
			if cs[i] == WILDCARD || j+1 < le && cs[j+1] != WILDCARD {
				i++
			}
		} else if cs[j] == WILDCARD { // 当c[i]不是万能牌c[j]是万能牌的时候
			for k := uint8(0); k < j; k++ {
				if removeEE(cs, flag, j, k) {
					return true
				}
			}
		}
	}
	return false
}

func detectAllTypeHu(cs []byte, ch, ps, ks []uint32, wildcard byte,luan int64) int64 {
	if hu3n2(cs) {
		//luan := existLuanFeng(getHandPongKong(cs, ch, ps, ks, wildcard))
		if luan > 0 {
			return HU | luan
		}
		return HU
	}
	thirteen := existThirteen(cs, wildcard)
	if thirteen > 0 {
		return HU | thirteen
	}

	//luan := existLuanFeng(getHandPongKong(cs, ch, ps, ks, wildcard))
	if luan > 0 {
		return HU | luan
	}

	pair7 := exist7pair(cs)
	if pair7 > 0 {
		return HU | pair7
	}
	return 0
}

// 3n +2 牌型胡牌 有序slices cs
func existHu3n2(cs []byte, ch, ps, ks []uint32, wildcard byte,luanfeng int64) int64 {
	hu := detectAllTypeHu(cs, ch, ps, ks, wildcard,luanfeng)
	if hu > 0 {
		return hu
	}
	// 财神是白板的情况作限制
	if wildcard != BAI {

		baicount := 0
		for _, v := range cs {
			if v == BAI {
				baicount ++
			}
		}

		if baicount > 0 {
			cards := make([]byte, len(cs))
			Sort(cs, 0, len(cs)-1)
			for k, v := range cs {
				if v == BAI {
					cs[k] = wildcard
					copy(cards, cs)
					Sort(cards, 0, len(cards)-1)
					hu := detectAllTypeHu(cards, ch, ps, ks, wildcard,luanfeng)
					if hu > 0 {
						return hu
					}
				}
			}
		}
	}

	return 0
}
