package algorithm

// 爆头检测（敲响）：
// 已经有十二张牌组合好，剩下一张牌为财神，并无将牌，抓到任意一张则为敲响。普通牌2倍，摸到白板4倍，摸到财神6倍。
// (已经有十二张牌组合好，剩下一张牌为财神，并无将牌，别人打出财神牌，或者白板算放炮吗)
// （此种情况为爆头，只能自摸，别人打出牌一律不能胡）
// （该种情况可与其他胡牌情况倍数累加）
func ExistBaoTou(cards []byte, ch, ps, ks []uint32, wildcard byte, card byte, draw bool, hu int64) int64 {
	if hu&HU_SINGLE > 0 || hu&HU_SINGLE_ZI > 0 ||
		hu&HU_SEVEN_PAIR > 0 ||
		hu&HU_KONG_FLOWER > 0 || hu&QIANG_GANG > 0 ||
		hu&TIAN_HU > 0 || hu&DI_HU > 0 ||
		hu&HU_QING_FENG > 0 || hu&HU_LUAN_FENG > 0 {
		return 0
	}
	handPongKong := GetHandPongKong(cards, ch, ps, ks, wildcard)

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
	cards = replaceWildcard(cards, wildcard, false)
	Sort(cards, 0, len(cards)-1) // 排序slices
	//glog.Errorln(cards[len(cards)-1] != WILDCARD)
	if cards[len(cards)-1] != WILDCARD {
		return 0
	}

	cards[len(cards)-1] = 0xFE
	for k, v := range cards {
		if v == card {
			cards[k] = 0xFE
			break
		}
	}

	luangfeng := ExistLuanFeng(handPongKong)
	color := ExistOneSuit(handPongKong, wildcard)
	value := ExistHu(cards, ch, ps, ks, wildcard, 0, color, luangfeng)

	if value > 0 {
		if card == BAI {
			return value | HU_BAO_TOU2
		} else if card == wildcard {
			return value | HU_BAO_TOU3
		}
		return value | HU_BAO_TOU1
	}
	return 0
}
