package room

import (
	"game/algorithm"
)

//结算,(明杠,放冲,庄家 - 收一家)
func (t *Desk) gameOver(huangZhuang bool) map[uint32]int32 {
	total := make(map[uint32]int32) // 总番数
	if huangZhuang { //黄庄
		return total
	}
	for k := uint32(1); k <= 4; k++ {
		if t.opt[k]&algorithm.HU > 0 {
			fan := algorithm.HuType(t.opt[k], t.data.Ante)
			fanType(t.lianCount, t.seat, k, fan, t.maizi, total)
			break
		}
	}
	return total
}

//lianCount连庄数,paoseat=放炮位置,huseat=胡牌位置
func fanType(lianCount, paoseat, huseat, score uint32, maizi map[uint32]uint32, huFan map[uint32]int32) map[uint32]int32 {
	if paoseat == huseat { //自摸,收三家
		for i := uint32(1); i <= 4; i++ {
			if i != huseat {
				huFan = over(huFan, huseat, i, (int32(maizi[huseat]+maizi[i])+1)*int32(score))
				//huFan = over(huFan, huseat, i, int32(score))
			}
		}
	} else { //炮胡,也收三家
		for i := uint32(1); i <= 4; i++ {
			if i != huseat {
				if i == paoseat {
					huFan = over(huFan, huseat, i, (int32(maizi[huseat]+maizi[i])+1)*int32(score))
					//huFan = over(huFan, huseat, i, int32(score))
				} else {
					huFan = over(huFan, huseat, i, (int32(maizi[huseat]+maizi[i])+1)*int32(float32(score)*0.5))
					//huFan = over(huFan, huseat, i, int32(float32(score)*0.5))
				}
			}
		}
	}
	return huFan
}

//收一家,s1=收的位置,s2=出的位置,val=收的番数
func over(m map[uint32]int32, s1, s2 uint32, val int32) map[uint32]int32 {
	m[s1] += val
	m[s2] -= val
	return m
}
