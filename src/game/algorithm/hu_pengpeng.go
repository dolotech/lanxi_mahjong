package algorithm

// 碰碰胡
// 有序slices cs
func ExistPengPeng(cs []byte,wildcard byte) int64 {
	//hasWildcard:= false
	//for _,v:=range cs{
	//	if v == WILDCARD{
	//		hasWildcard = true
	//		break
	//	}
	//}

	value := hu3n2PengPeng(cs)
	if value > 0 {
		//if hasWildcard{
			//return  HU_PENG_PENG_CAI
		//}
		return  HU_PENG_PENG
	}
	// 财神是白板的情况作限制
	if wildcard != BAI {
		cards := make([]byte, len(cs))
		copy(cards, cs)
		for k, v := range cards {
			if v == BAI {
				cards[k] = wildcard
				Sort(cards, 0, len(cards)-1)
				value := hu3n2PengPeng(cards)
				if value > 0 {
					//if hasWildcard{
					//	return  HU_PENG_PENG_CAI
					//}
					return  HU_PENG_PENG
				}
			}
		}
	}
	return 0
}



func hu3n2PengPeng(cs []byte) int64 {
	le := uint8(len(cs))
	var bgflag int32
	for i := uint8(0); i < le; i++ {
		bgflag = bgflag | (1 << i)
	}
	for n := uint8(0); n < le-1; n++ {
		 flag  := bgflag
		for g := n + 1; g < le; g++ {
			if cs[n] == cs[g] || cs[n] == WILDCARD || cs[g] == WILDCARD {
				flag = flag & (^(1 << n))
				flag = flag & (^(1 << g))
				break
			}
		}
		if flag == bgflag {
			continue
		}
		if check3nPengPeng(cs, &flag) {
			return HU
		}
	}
	return 0
}
