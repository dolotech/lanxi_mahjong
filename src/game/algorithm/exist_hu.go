package algorithm

func DetectKong(cs []byte, ps []uint32, wildcard byte) (status int64) {
	if len(existAnKong(cs, wildcard)) > 0 {
		status |= AN_KONG
		status |= KONG
	}
	if existBuKong(cs, ps, wildcard) {
		status |= BU_KONG
		status |= KONG
	}
	return
}

// 判断是否胡牌,0表示不胡牌,非0用32位表示不同的胡牌牌型
func ExistHu(cards []byte, ch, ps, ks []uint32, wildcard byte, card byte,color int64,luanfeng int64) int64 {
	le := len(cards)
	if card > 0 {
		le = le + 1
		cs := make([]byte, le)
		copy(cs, cards)
		cs[le-1] = card
		cards = cs
	} else {
		cs := make([]byte, le)
		copy(cs, cards)
		cards = cs
	}

	// 替换万能牌
	cards = replaceWildcard(cards, wildcard, false)
	Sort(cards, 0, len(cards)-1) //排序slices
	//单钓胡牌
	if le == 2 && (cards[0] == cards[1] || cards[0] == WILDCARD || cards[1] == WILDCARD) {
		value := HU
		//handPongKong := getHandPongKong(cards, ch, ps, ks, wildcard)
		qiuren := quanQiuRen(cards, ch, ps, ks)
		value = value | qiuren

		// 清风检测
		//if existLuanFeng(handPongKong) > 0 {
		if luanfeng> 0 {
			value &= (^HU_LUAN_FENG)
			value = HU_QING_FENG | value
			return value
		}

		if len(ch) == 0 {
			seven := ExistPengPeng(cards, wildcard)
			if seven > 0 {
				value = seven | value
			}
		}
		// 清一色检测
		//qing := existOneSuit(handPongKong, wildcard)
		if value > 0{
			value = value | color
		}
		return value
	}

	value := existHu3n2(cards , ch, ps, ks ,  wildcard,luanfeng)
	//是否3n+2牌型
	if value > 0 {
		// 一财一刻、两财一刻、三财一刻检测
		if value &HU_SINGLE ==0 && value&HU_SINGLE_ZI == 0 /*&& value &HU_SINGLE_ZI ==0*/{// 13不靠风内部已经检测，不再检测财神归位
			tv := ExistNCaiNKe(cards, ch, ps, ks, wildcard,luanfeng)
			value = value | tv
		}

		//handPongKong := getHandPongKong(cards, ch, ps, ks, wildcard)
		// 清风检测
		//if existLuanFeng(handPongKong) > 0 {
		if luanfeng> 0 {
			value &= (^HU_LUAN_FENG)
			value = HU_QING_FENG | value
			return value
		}
		// 清一色检测
		//color := existOneSuit(handPongKong, wildcard)
		if color > 0 {
			value = color | value
		}
		// 碰碰胡，有吃牌就不算碰碰胡
		if len(ch) == 0 {
			seven := ExistPengPeng(cards, wildcard)
			if seven > 0 {
				value = seven | value
			}
		}
		return value
	}

	return 0
}

//三个财神加倍
func ThreeWildcard(handcard []byte, wildcard byte) int64 {
	count := 0
	for _, v := range handcard {
		if wildcard == v {
			count ++
			if count == 3 {
				return HU_3_CAI_SHEN
			}
		}
	}
	return 0
}

// 全求人检测
func quanQiuRen(cards []byte, ch, ps, ks []uint32) int64 {
	leks := len(ks)
	for i := 0; i < leks; i++ {
		v, _, _ := DecodeKong(ks[i])
		if int64(v)&AN_KONG > 0 || int64(v)&BU_KONG > 0 {
			return 0
		}
	}
	return HU_QUAN_QIU_REN
}

// 合并手牌、杠牌、碰牌
func GetHandPongKong(cards []byte, ch, ps, ks []uint32, wildcard byte) []byte {
	le := len(cards)
	leps := len(ps)
	leks := len(ks)
	lechaw := len(ch)
	handpengkong := make([]byte, le,le+leps+leks+lechaw*3)
	copy(handpengkong, cards)

	for i := 0; i < leps; i++ {
		_, c := DecodePeng(ps[i])
		handpengkong = append(handpengkong, c)

	}
	for i := 0; i < leks; i++ {
		_, c, _ := DecodeKong(ks[i])
		handpengkong = append(handpengkong, c)
	}
	for i := 0; i < lechaw; i++ {
		c1, c2, c3 := DecodeChow(ch[i])
		// 把白板替换成财神本尊
		if c1 == BAI {
			c1 = wildcard
		} else if c2 == BAI {
			c2 = wildcard
		} else if c3 == BAI {
			c3 = wildcard
		}
		handpengkong = append(handpengkong, c1)
		handpengkong = append(handpengkong, c2)
		handpengkong = append(handpengkong, c3)
	}
	return handpengkong
}
