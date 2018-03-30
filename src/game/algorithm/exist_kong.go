package algorithm

//是否存在暗杠
func existAnKong(cards []byte, wildcard byte) (kong []byte) {
	le := len(cards)
	for j := 0; j < le-3; j++ {
		count := 0
		for i := j + 1; i < le; i++ {
			if cards[j] == cards[i] && cards[i] != wildcard {
				count = count + 1
				if count == 3 {
					kong = append(kong, cards[i])
					break
				}
			}
		}
	}
	return
}

//是否存在碰
func existPeng(card byte, cards []byte, wildcard byte) bool {
	// 财神不可以碰
	if card == wildcard {
		return false
	}
	le := len(cards)
	count := 0
	for i := 0; i < le; i++ {
		if card == cards[i] {
			count = count + 1
			if count == 2 {
				return true
			}
		}
	}
	return false
}

//是否存在补杠
func existBuKong(handcards []byte, pongs []uint32, wildcard byte) bool {
	le := len(pongs)
	for i := 0; i < le; i++ {
		_, c := DecodePeng(pongs[i])
		if wildcard != c {
			for _, v := range handcards {
				if v == c {
					return true
				}
			}
		}

	}
	return false
}

//是否存在明杠
func existMingKong(card byte, cards []byte, wildcard byte) bool {
	// 财神不可以杠
	if card == wildcard {
		return false
	}
	le := len(cards)
	count := 0
	for i := 0; i < le; i++ {
		if card == cards[i] {
			count = count + 1
			if count == 3 {
				return true
			}
		}
	}
	return false
}
