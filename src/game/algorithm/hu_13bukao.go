package algorithm


// 风未集齐，三家x（5或者6张财神）
// 风集齐，三家2x
//  财神为风，三家4x
// 13不靠(数牌相隔2个或以上,字牌不重复)/13不靠风未齐(未集齐7个不同字牌)
// 有序slices cs && len(cs) == 14
func existThirteen(cs []byte, wildcard byte) int64 {
	if len(cs) != 14 {
		return 0
	}
	thirteen_single := true
	thirteen_single_zi := true
	var c byte
	var n int
	var huType int64
	for _, v := range cs {
		if v != WILDCARD{
			if v>>4 < FENG {//不是风牌
				//不是相同类型的牌继续检测下一张牌
				if c == 0 || c>>4 != v>>4 {
					c = v
					continue
				}
				//没有相隔两张牌
				if v-c < 0x03 {
					thirteen_single = false
					break
				}
				c = v
			} else { //风牌或字牌
				n++
				if v == c { //出现两张一样的风牌时中断检测
					thirteen_single_zi = false
					break
				}
				c = v
			}
		}
	}

	wildcardCount:= 0
	for _, v := range cs {
		if v == WILDCARD{
			wildcardCount++
		}
	}

	//13不靠风集齐 
	if thirteen_single && thirteen_single_zi && n +wildcardCount >= 7 {
		// 13不靠检测财神归位
		if wildcard>>4 >= FENG {
			huType |= HU_GUI_WEI1
		}
		huType |= HU_SINGLE_ZI
	//13不靠风未齐
	} else if thirteen_single && thirteen_single_zi {
		huType |= HU_SINGLE
	}
	return huType
}
