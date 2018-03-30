package algorithm

import (
	"time"
	"math/rand"
)

func SearchDirtyCard(cs []byte, wildcard byte) byte {
	le := uint8(len(cs))
	if le == 0 {
		return 0
	}

	//
	if le == 2{
		return cs[1]
	}

	var flag int32
	for i := uint8(0); i < le; i++ {
		flag = flag | (1 << i)
	}
	Sort(cs, 0, len(cs)-1)

	check3n(cs,&flag)

	// 去掉万能牌
	for i := uint8(0); i < le; i++ {
		if ((flag>>i)&0x1) == 1 && cs[i] == wildcard {
			flag = flag & (^(1 << i))
		}
	}

	// 去掉对子
	for i := uint8(0); i < le-1; i++ {
		for j := i + 1; j < le && (((flag >> i) & 0x1) == 1); j++ {
			if ((flag>>j)&0x1) == 1 && cs[j] == cs[i] {
				flag = flag & (^(1 << i))
				flag = flag & (^(1 << j))
				break
			}
		}
	}

	// 去掉搭子
	for i := uint8(0); i < le-1; i++ {
		for j := i + 1; j < le && (((flag >> i) & 0x1) == 1); j++ {
			if ((flag>>j)&0x1) == 1 && cs[i]>>4 < FENG && cs[j]+1 == cs[i] && cs[j]&0x0f != 1 && cs[i]&0x0f != 9 {
				flag = flag & (^(1 << i))
				flag = flag & (^(1 << j))
				break
			}
		}
	}

	// 去掉卡窿
	for i := uint8(0); i < le-1; i++ {
		for j := i + 1; j < le && (((flag >> i) & 0x1) == 1); j++ {
			if ((flag>>j)&0x1) == 1 && cs[i]>>4 < FENG && cs[j]+2 == cs[i] && cs[j]&0x0f != 1 && cs[i]&0x0f != 9 {
				flag = flag & (^(1 << i))
				flag = flag & (^(1 << j))
				break
			}
		}
	}

	// 风牌先除
	for i := uint8(0); i < le; i++ {
		if ((flag>>i)&0x1) == 1 && (cs[i]>>4 >= FENG) {
			return cs[i]
		}
	}

	// 1/9先除
	for i := uint8(0); i < le; i++ {
		if ((flag>>i)&0x1) == 1 && ( cs[i]&0x0f == 1 || cs[i]&0x0f == 9 ) {
			return cs[i]
		}
	}
	bet := make([]byte, 0, len(cs))

	for i := uint8(0); i < le; i++ {
		if ((flag >> i) & 0x1) == 1 {
			bet = append(bet, cs[i])
		}
	}

	if len(bet) ==1{
		return bet[0]
	}else if len(bet) > 1 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		return bet[int(r.Intn(len(bet)))]
	}


	return cs[le-1]
}
