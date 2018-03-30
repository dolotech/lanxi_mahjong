package algorithm

// 乱风胡检测
// 所有牌都为风牌(东南西北中发白)，无需成胡牌型
// 碰杠牌要加入检测
func ExistLuanFeng(cards []byte) int64 {
	for _, v := range cards {
		if v != WILDCARD {
			if v>>4 < FENG {
				return 0
			}
		}
	}
	return HU_LUAN_FENG
}

// 清一色检测
func ExistOneSuit(cards []byte,wildcard byte) int64 {
	var c byte
	le := len(cards)
	for i := 0; i < le; i++ {
		card := cards[i]
		//if card == BAI{
		//	card = wildcard
		//}
		if card != wildcard {
			if c > 0 && c>>4 != card>>4 {
				return 0
			}
			c = card
		}
	}

	return HU_ONE_SUIT
}
