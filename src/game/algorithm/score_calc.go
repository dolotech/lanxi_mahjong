package algorithm

//算番(牌型) 多个牌型时相加
func HuType(value int64, ante uint32) uint32 {
	if value&HU == 0 {
		return 0
	}
	var fan uint32

	for k, v := range Fan {
		if value&k > 0 && k != HU && k != HU_SINGLE {
			fan += (ante * v)
			if fan > 12 {
				fan = 6 * ante
				break
			}
		}
	}

	// 没有大牌则胡平胡
	if fan == 0 {
		fan += (ante * Fan[HU])
	}
	return fan
}
