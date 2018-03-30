package algorithm

import "github.com/golang/glog"

// 验证吃, c1,c2,c3有序
func VerifyChow(cards []byte, c1, c2, c3 byte, wildcard byte) bool {

	glog.Infof("%x %x %x",c1,c2,c3)
	if !Exist(c2, cards, 1) {
		return false
	}
	glog.Infof("%x %x %x",c1,c2,c3)
	if !Exist(c3, cards, 1) {
		return false
	}
	glog.Infof("%x %x %x",c1,c2,c3)
	var c []byte = []byte{c1, c2, c3}
	for k, v := range c {
		// 财神代替白板
		if v == BAI {
			c[k] = wildcard
		}
	}

	Sort(c, 0, len(c)-1)
	if (c[0]>>4) >= FENG || (c[1]>>4) >= FENG || (c[2]>>4) >= FENG {
		return false
	}
	return c[0]+1 == c[1] && c[1]+1 == c[2]
}

//检测,是否存在吃
func existChow(card byte, cs []byte, wildcard byte) ([]byte) {
	rel := []byte{}
	// 打出的是财神,不可以吃
	if card == wildcard {
		return rel
	}
	// 打出非白板的风牌
	if (card>>4) >= FENG && card != BAI {
		return rel
	}
	// 打出的是白板财神又是风牌的情况
	if card == BAI && (wildcard>>4) >= FENG {
		return rel
	}

	le := len(cs)
	if le == 1 {
		return rel
	}

	// 白板代替财神
	if card == BAI {
		card = wildcard
	}

	if (card & 0xF) > 1 {
		tcard:=findSame(cs, card-1, wildcard)
		if tcard > 0 {
			rel = append(rel, tcard)
			if (card & 0xF) > 2 {
				tcard:=findSame(cs, card-2, wildcard)
				if tcard > 0 {
					rel = append(rel, tcard)
				}
			}
		}
	}
	if (card & 0xF) <9 {
		tcard:=findSame(cs, card+1, wildcard)
		if tcard > 0 {
			rel = append(rel, tcard)
			if (card & 0xF) <8{
				tcard:=findSame(cs, card+2, wildcard)
				if tcard > 0 {
					rel = append(rel, tcard)
				}
			}
		}
	}
	return rel
}
func findSame(cs []byte, card byte, wildcard byte) byte {
	for _, v := range cs {
		// 白板代替财神
		if v == BAI && wildcard != BAI {
			v = wildcard
			if v == card {
				return BAI
			}
		} else {
			// 财神不参与吃
			if v == card && wildcard != v{
				return v
			}
		}
	}
	return 0
}
