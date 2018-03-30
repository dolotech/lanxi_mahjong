package algorithm
//
//func guiweiDetect(cs []byte, ch, ps, ks []uint32, wildcard byte)  {
//
//}

// 一财一刻、两财一刻、三财一刻检测
// 财神归位检测
// 升序的slice、一定为胡牌型
// todo 两财一刻和财神归位同时存在的情况
func guiwei(cs []byte, ch, ps, ks []uint32, wildcard byte,luanfeng int64) (value int64) {
	caishenCount := 0
	for _, v := range cs {
		if v == WILDCARD {
			caishenCount ++
		}
	}

	if caishenCount == 0 {
		return 0
	}

	list := make([]byte, len(cs))
	copy(list, cs)
	cs = list

	// 复原财神
	for k, v := range cs {
		if v == WILDCARD {
			cs[k] = wildcard
		}
	}
	Sort(cs, 0, len(cs)-1)

	count := caishenCount
	for i := 0; i < caishenCount; i++ {
		for k, v := range cs {
			if v == wildcard {

				if existHu3n2(cs,ch, ps, ks ,  wildcard,luanfeng) >0{
					if count == 1 {
						value = HU_GUI_WEI1
						return
					} else if count == 2 {
						value = HU_GUI_WEI2
						return
					} else if count == 3 {
						value = HU_GUI_WEI3
						return
					}
				}

				cs[k] = WILDCARD
				Sort(cs, 0, len(cs)-1)
				count --
				break
			}
		}
	}

	return
}

//  清一色／碰碰和／小七对  牌型有无财神对 检测
func ExistCaiShen(cards []byte, hu int64, wildcard byte, card byte, draw bool) int64 {
	le := len(cards)
	if !draw {
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

	count := 0

	if hu&HU_BAO_TOU1 > 0 || hu&HU_BAO_TOU2 > 0 || hu&HU_BAO_TOU3 > 0 {
		wildcardFlag := false
		cardFlag := false

		for _, v := range cards {
			if !cardFlag || v == card {
				cardFlag = true
				continue
			}
			if wildcard == v {
				if !wildcardFlag {
					wildcardFlag = true
					continue
				}
				count ++
			}
		}
	} else {
		for _, v := range cards {
			if wildcard == v {
				count ++
			}
		}
	}

	if count == 0 {
		return hu
	}

	use := 0
	if hu&HU_CAI_1 > 0 {
		use++
	}
	if hu&HU_GUI_WEI1 > 0 {
		use++
	}

	if hu&HU_CAI_2 > 0 {
		use += 2
	}
	if hu&HU_GUI_WEI2 > 0 {
		use += 2
	}
	if hu&HU_CAI_3 > 0 {
		use += 3
	}
	if hu&HU_GUI_WEI3 > 0 {
		use += 3
	}

	if use == count {
		return hu
	}

	if hu&HU_PENG_PENG > 0 {
		hu = hu & (^HU_PENG_PENG)
		hu = hu | HU_PENG_PENG_CAI
	}
	if hu&HU_SEVEN_PAIR > 0 {
		hu = hu & (^HU_SEVEN_PAIR)
		hu = hu | HU_SEVEN_PAIR_CAI
	}
	if hu&HU_ONE_SUIT > 0 {
		hu = hu & (^HU_ONE_SUIT)
		hu = hu | HU_ONE_SUIT_CAI
	}

	return hu

}
func ExistNCaiNKe(cs []byte, ch, ps, ks []uint32, wildcard byte,luanfeng int64) int64 {
	// 计算财神的数量
	caishenCount := 0
	for _, v := range cs {
		if v == WILDCARD {
			caishenCount ++
		}
	}

	if caishenCount == 0 {
		return 0
	}

	if len(cs) >= 5 {
		// 3财1刻
		if caishenCount == 3 {
			list := make([]byte, 0, len(cs))
			for _, v := range cs {
				if v != WILDCARD {
					list = append(list, v)
				}
			}
			Sort(list, 0, len(list)-1)
			if existHu3n2(list, ch, ps, ks , wildcard,luanfeng)>0  {
				return HU_CAI_3
			}
		}

		baiCount := 0
		for _, v := range cs {
			if v == BAI {
				baiCount ++
			}
		}

		// 2财1刻
		if baiCount >= 1 && caishenCount >= 2 {
			list := make([]byte, 0, len(cs))
			bai := 0
			cai := 0
			for _, v := range cs {
				if bai < 1 && v == BAI {
					bai ++
					continue
				}

				if cai < 2 && v == WILDCARD {
					cai ++
					continue
				}

				list = append(list, v)
			}

			Sort(list, 0, len(list)-1)
			if existHu3n2(list, ch, ps, ks , wildcard,luanfeng) >0{
				value := guiwei(list, ch, ps, ks, wildcard,luanfeng)
				return HU_CAI_2 | value
			}
		}
		// 1财1刻
		if baiCount >= 2 && caishenCount >= 1 {

			list := make([]byte, 0, len(cs))
			bai := 0
			cai := 0
			for _, v := range cs {
				if bai < 2 && v == BAI {
					bai ++
					continue
				}

				if cai < 1 && v == WILDCARD {
					cai ++
					continue
				}

				list = append(list, v)
			}

			Sort(list, 0, len(list)-1)
			if existHu3n2(list,ch, ps, ks ,  wildcard,luanfeng)>0  {
				value := guiwei(list, ch, ps, ks, wildcard,luanfeng)
				return HU_CAI_1 | value
			}
		}
	}
	return guiwei(cs, ch, ps, ks, wildcard,luanfeng)
}
