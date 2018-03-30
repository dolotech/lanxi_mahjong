package main

import (
	"encoding/json"
	"io/ioutil"
	"fmt"
	"os"
	"strconv"
	"strings"
)

/*
//json data
{
	"pao":0,       //放冲位置(自摸时填自己位置)
	"dealer":0,    //庄家位置
	"seat1":{                  //位置
		"hu":"HU,PAOHU",       //胡牌类型,多个逗号隔开
		"cards":"W8,W8,W8,W8", //手牌,(大对子暗杠数,胡牌牌型)
		"pong": {       //碰  
			"seat":"2", //被碰位置 
			"card":"B1" //碰的牌值
		},
		"chow": {              //吃
			"seat":"4",        //上家位置(被吃位置)
			"card":"T1,T2,T3"  //吃的牌
		},
		"mingKong": {     //明杠
			"seat":"2",   //被杠位置
			"card":"T7"   //明杠牌值
		},
		"anKong": {       //暗杠
			"seat":"0",   //暗杠位置可以为0
			"card":"B7"   //暗杠牌值
		},
		"buKong": {        //补杠
			"seat":"0,0",  //补杠位置可以为0
			"card":"F1,F2" //补杠牌值
		}
	},
	"seat2":{
	},
	"seat3":{
	},
	"seat4":{
	}
}
*/

var Conf Config

type Config struct {
	Pao    uint32 `json:"pao"`
	Dealer uint32 `json:"dealer"`
	Seat1  Seats  `json:"seat1"`
	Seat2  Seats  `json:"seat2"`
	Seat3  Seats  `json:"seat3"`
	Seat4  Seats  `json:"seat4"`
}

type Seats struct {
	Hu       string `json:"hu"`
	Cards    string `json:"cards"`
	Pong     Kong   `json:"pong"`     //多个碰
	Chow     Kong   `json:"chow"`     //多个吃
	MingKong Kong   `json:"mingKong"` //多个杠
	AnKong   Kong   `json:"anKong"`
	BuKong   Kong   `json:"buKong"`
}

type Kong struct {
	Seat  string `json:"seat"`
	Card  string `json:"card"`
}

//解析配置
func Parse() {
	//fmt.Println("Conf -> ", Conf)
	PAOSEAT = Conf.Pao   //放冲位置
	DEALER = Conf.Dealer //庄家位置
	parsing(1, Conf.Seat1)
	parsing(2, Conf.Seat2)
	parsing(3, Conf.Seat3)
	parsing(4, Conf.Seat4)
}

func parsing(seat uint32, c Seats) {
	parseHu(seat, c.Hu)
	parseCards(seat, c.Cards)
	parseKongs(MING_KONG, seat, c.MingKong)
	parseKongs(AN_KONG, seat, c.AnKong,)
	parseKongs(BU_KONG, seat, c.BuKong)
	parsePong(seat, c.Pong)
	parseChow(seat, c.Chow)
}

func parseHu(seat uint32, hu string) {
	if len(hu) == 0 {
		return
	}
	var sc []string = strings.Split(hu, ",")
	var mask uint32
	for _, v := range sc {
		val := HU2MASK[v]
		if val == 0 {
			panic(fmt.Sprintf("parse hu err -> %v", v))
		}
		mask |= val
	}
	if mask > 0 {
		HUS[seat] = mask
	}
}

func parseCards(seat uint32, cards string) {
	if len(cards) == 0 {
		return
	}
	var sc []string = strings.Split(cards, ",")
	var cs []byte = []byte{}
	for _, v := range sc {
		c := STR2CARD[v]
		if c == 0 {
			panic(fmt.Sprintf("parse cards err -> %v", v))
		}
		cs = append(cs, c)
	}
	if len(cs) > 0 {
		HANDS[seat] = cs
	}
}

func parseKongs(mask, seat uint32, ks Kong) {
	if len(ks.Seat) == 0 {
		return
	}
	var s []string = strings.Split(ks.Seat, ",")
	var c []string = strings.Split(ks.Card, ",")
	if len(s) != len(c) {
		panic(fmt.Sprintf("parse kong err -> %v, %v", s, c))
	}
	var cs []uint32 = []uint32{}
	if cs1, ok := KONGS[seat]; ok {
		cs = cs1
	}
	for k, v := range s {
		s_v, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("parse kong err -> %v", err))
		}
		c_v := STR2CARD[c[k]]
		if c_v == 0 {
			panic(fmt.Sprintf("parse kong err -> %v", c[k]))
		}
		cs = append(cs, EncodeKong(uint32(s_v), byte(c_v), mask))
	}
	if len(cs) > 0 {
		KONGS[seat] = cs
	}
}

func parsePong(seat uint32, ps Kong) {
	if len(ps.Seat) == 0 {
		return
	}
	var s []string = strings.Split(ps.Seat, ",")
	var c []string = strings.Split(ps.Card, ",")
	if len(s) != len(c) {
		panic(fmt.Sprintf("parse pong err -> %v, %v", s, c))
	}
	var cs []uint32 = []uint32{}
	if cs1, ok := PONGS[seat]; ok {
		cs = cs1
	}
	for k, v := range s {
		s_v, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("parse pong err -> %v", err))
		}
		c_v := STR2CARD[c[k]]
		if c_v == 0 {
			panic(fmt.Sprintf("parse pong err -> %v", c[k]))
		}
		cs = append(cs, EncodePeng(uint32(s_v), byte(c_v)))
	}
	if len(cs) > 0 {
		PONGS[seat] = cs
	}
}

func parseChow(seat uint32, ws Kong) {
	if len(ws.Seat) == 0 {
		return
	}
	var s []string = strings.Split(ws.Seat, ",")
	var c []string = strings.Split(ws.Card, ",")
	if len(s) != (len(c) / 3) {
		panic(fmt.Sprintf("parse chow err -> %v, %v", s, c))
	}
	var cs []uint32 = []uint32{}
	if cs1, ok := CHOWS[seat]; ok {
		cs = cs1
	}
	for k, v := range s {
		_, err := strconv.Atoi(v)
		if err != nil {
			panic(fmt.Sprintf("parse pong err -> %v", err))
		}
		c_v1 := STR2CARD[c[k]]
		c_v2 := STR2CARD[c[k+1]]
		c_v3 := STR2CARD[c[k+2]]
		if c_v1 == 0 || c_v2 == 0 || c_v3 == 0 {
			panic(fmt.Sprintf("parse pong err -> %v", c[k]))
		}
		cs = append(cs, EncodeChow(byte(c_v1),byte(c_v1),byte(c_v2)))
	}
	if len(cs) > 0 {
		CHOWS[seat] = cs
	}
}

//加载配置
func LoadConf() {
	var path string = "./conf.json"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	//fmt.Println(data)
	err = json.Unmarshal(data, &Conf)
	if err != nil {
		panic(err)
	}
}

// 万牌
const (
	W1  = iota + 0x01
	W2
	W3
	W4
	W5
	W6
	W7
	W8
	W9
)

// 条牌
const (
	T1  = iota + 0x11
	T2
	T3
	T4
	T5
	T6
	T7
	T8
	T9
)

// 筒牌
const (
	B1  = iota + 0x21
	B2
	B3
	B4
	B5
	B6
	B7
	B8
	B9
)

// 字牌
const (
	Z1  = iota + 0x41
	Z2
	Z3
	Z4
)

// 风牌
const (
	F1  = iota + 0x51
	F2
	F3
	F4
)

const (
	// 牌局基础的常量
	TOTAL uint32 = 136 //一副贵州麻将的总数
	BING  uint32 = 2   //同子类型
	TIAO  uint32 = 1   //条子类型
	WAN   uint32 = 0   //万字类型
	FENG  uint32 = 4   //风牌类型
	ZI    uint32 = 5   //字牌类型
	HAND uint32 = 13 //手牌数量
	SEAT uint32 = 4  //最多可参与一桌打牌的玩家数量,不算旁观
	// 碰杠胡掩码,用32位每位代表不同的状态
	DRAW      uint32 = 0      // 摸牌 
	DISCARD   uint32 = 1      // 打牌
	PENG      uint32 = 2 << 0 // 碰
	MING_KONG uint32 = 2 << 1 // 明杠
	AN_KONG   uint32 = 2 << 2 // 暗杠
	BU_KONG   uint32 = 2 << 3 // 补杠
	KONG      uint32 = 2 << 4 // 杠(代表广义的杠)
	CHOW      uint32 = 2 << 5 // 吃
	HU        uint32 = 2 << 6 // 胡(代表广义的胡)
	//胡牌类型
	HU_PING            uint32 = 2 << 8  // 平胡
	HU_SINGLE          uint32 = 2 << 9  // 十三烂
	HU_SINGLE_ZI       uint32 = 2 << 10 // 七星十三烂
	HU_SEVEN_PAIR_BIG  uint32 = 2 << 11 // 大七对
	HU_SEVEN_PAIR      uint32 = 2 << 12 // 小七对
	HU_SEVEN_PAIR_KONG uint32 = 2 << 13 // 豪华小七对
	HU_ONE_SUIT        uint32 = 2 << 14 // 清一色
	HU_ALL_ZI          uint32 = 2 << 15 // 字一色
	ZIMO           uint32 = 2 << 17 // 自摸
	PAOHU          uint32 = 2 << 18 // 炮胡,也叫放冲
	QIANG_GANG     uint32 = 2 << 19 // 抢杠,其他家胡你补杠那张牌
	HU_KONG_FLOWER uint32 = 2 << 20 // 杠上开花
	HU_MENGQING    uint32 = 2 << 21 // 门清
	HU_DANDIAO     uint32 = 2 << 22 // 单钓
	TIAN_HU        uint32 = 2 << 23 // 天胡
	DI_HU          uint32 = 2 << 24 // 地胡
)

var HUS map[uint32]uint32 = make(map[uint32]uint32)   //胡数据
var KONGS map[uint32][]uint32 = make(map[uint32][]uint32) //杠数据
var PONGS map[uint32][]uint32 = make(map[uint32][]uint32) //碰数据
var CHOWS map[uint32][]uint32 = make(map[uint32][]uint32) //吃数据
var HANDS map[uint32][]byte = make(map[uint32][]byte) //手牌
var DEALER uint32 //庄家位置
var PAOSEAT uint32 //放冲位置
//掩码值转换
var HU2MASK map[string]uint32 = make(map[string]uint32)
var MASK2HU map[uint32]string = make(map[uint32]string)
var HU2STR map[uint32]string = make(map[uint32]string)
//牌值转换
var STR2CARD map[string]byte = make(map[string]byte)
//番值
var FAN map[uint32]int32 = make(map[uint32]int32)
//初始化
func init() {
	//牌型
	FAN[HU_PING           ] = 1	 // 平胡
	FAN[HU_SINGLE         ] = 2	 // 十三烂
	FAN[HU_SINGLE_ZI      ] = 3	 // 七星十三烂
	FAN[HU_SEVEN_PAIR_BIG ] = 3	 // 大七对
	FAN[HU_SEVEN_PAIR     ] = 4	 // 小七对
	FAN[HU_SEVEN_PAIR_KONG] = 4	 // 豪华小七对
	FAN[HU_ONE_SUIT       ] = 5	 // 清一色
	FAN[HU_ALL_ZI         ] = 12 // 字一色
	//胡牌方式
	FAN[QIANG_GANG    ] = 2  // 抢杠,其他家胡你补杠那张牌
	FAN[HU_KONG_FLOWER] = 2  // 杠上开花,杠完牌抓到的第一张牌自摸了
	FAN[HU_MENGQING   ] = 2  // 门清
	FAN[HU_DANDIAO    ] = 2  // 单钓
	FAN[TIAN_HU       ] = 12 // 天胡
	FAN[DI_HU         ] = 12 // 地胡
	//杠牌
	FAN[MING_KONG] = 1 // 明杠
	FAN[AN_KONG  ] = 2 // 暗杠
	FAN[BU_KONG  ] = 1 // 补杠

	//初始化
	HU2MASK["HU"]                = HU                 //胡
	HU2MASK["HU_PING"]           = HU_PING            //平胡
	HU2MASK["HU_SINGLE"]         = HU_SINGLE          //十三烂
	HU2MASK["HU_SINGLE_ZI"]      = HU_SINGLE_ZI       //七星十三烂
	HU2MASK["HU_SEVEN_PAIR_BIG"] = HU_SEVEN_PAIR_BIG  //大七对
	HU2MASK["HU_SEVEN_PAIR"]     = HU_SEVEN_PAIR      //小七对
	HU2MASK["HU_SEVEN_PAIR_KONG"]= HU_SEVEN_PAIR_KONG //豪华小七对
	HU2MASK["HU_ONE_SUIT"]       = HU_ONE_SUIT        //清一色
	HU2MASK["HU_ALL_ZI"]         = HU_ALL_ZI          //字一色
	HU2MASK["ZIMO"]              = ZIMO               //自摸
	HU2MASK["PAOHU"]             = PAOHU              //炮胡,也叫放冲
	HU2MASK["QIANG_GANG"]        = QIANG_GANG         //抢杠,其他家胡你补杠那张牌
	HU2MASK["HU_KONG_FLOWER"]    = HU_KONG_FLOWER     //杠上开花,杠完牌抓到的第一张牌自摸了
	HU2MASK["HU_MENGQING"]       = HU_MENGQING        //门清
	HU2MASK["HU_DANDIAO"]        = HU_DANDIAO         //单钓
	HU2MASK["TIAN_HU"]           = TIAN_HU            //天胡
	HU2MASK["DI_HU"]             = DI_HU              //地胡
	//初始化
	HU2STR[HU]                = "胡"
	HU2STR[HU_PING]           = "平胡"
	HU2STR[HU_SINGLE]         = "十三烂"
	HU2STR[HU_SINGLE_ZI]      = "七星十三烂"
	HU2STR[HU_SEVEN_PAIR_BIG] = "大七对"
	HU2STR[HU_SEVEN_PAIR]     = "小七对"
	HU2STR[HU_SEVEN_PAIR_KONG]= "豪华小七对"
	HU2STR[HU_ONE_SUIT]       = "清一色"
	HU2STR[HU_ALL_ZI]         = "字一色"
	HU2STR[ZIMO]              = "自摸"
	HU2STR[PAOHU]             = "炮胡"
	HU2STR[QIANG_GANG]        = "抢杠"
	HU2STR[HU_KONG_FLOWER]    = "杠上开花"
	HU2STR[HU_MENGQING]       = "门清"
	HU2STR[HU_DANDIAO]        = "单钓"
	HU2STR[TIAN_HU]           = "天胡"
	HU2STR[DI_HU]             = "地胡"

	//
	STR2CARD["W1"] = W1 //万
	STR2CARD["W2"] = W2
	STR2CARD["W3"] = W3
	STR2CARD["W4"] = W4
	STR2CARD["W5"] = W5
	STR2CARD["W6"] = W6
	STR2CARD["W7"] = W7
	STR2CARD["W8"] = W8
	STR2CARD["W9"] = W9
	STR2CARD["T1"] = T1 //条
	STR2CARD["T2"] = T2
	STR2CARD["T3"] = T3
	STR2CARD["T4"] = T4
	STR2CARD["T5"] = T5
	STR2CARD["T6"] = T6
	STR2CARD["T7"] = T7
	STR2CARD["T8"] = T8
	STR2CARD["T9"] = T9
	STR2CARD["B1"] = B1 //筒
	STR2CARD["B2"] = B2
	STR2CARD["B3"] = B3
	STR2CARD["B4"] = B4
	STR2CARD["B5"] = B5
	STR2CARD["B6"] = B6
	STR2CARD["B7"] = B7
	STR2CARD["B8"] = B8
	STR2CARD["B9"] = B9
	STR2CARD["Z1"] = Z1 //字
	STR2CARD["Z2"] = Z2
	STR2CARD["Z3"] = Z3
	STR2CARD["Z4"] = Z4
	STR2CARD["F1"] = F1 //风
	STR2CARD["F2"] = F2
	STR2CARD["F3"] = F3
}

func Mask2HuInit() {
	for k, v := range HU2MASK {
		MASK2HU[v] = k
	}
}

func main() {
	LoadConf()    //加载配置
	Mask2HuInit() //初始化
	Parse()       //解析配置
	//结算
	huFan, mingKong, beMingKong, anKong, buKong, total := gameOver()
	//打印输出
	fmt.Printf("庄家位置:%d\n", DEALER)
	fmt.Printf("放炮位置:%d\n", PAOSEAT)
	fmt.Printf("胡牌掩码:%+v\n", HUS)
	fmt.Printf("杠牌掩码:%+v\n", KONGS)
	fmt.Printf("碰牌掩码:%+v\n", PONGS)
	fmt.Printf("吃牌掩码:%+v\n", CHOWS)
	fmt.Printf("手牌数据:%+v\n", HANDS)
	//打印结算
	fmt.Printf("牌型分:%+v\n", huFan)
	fmt.Printf("明杠分:%+v\n", mingKong)
	fmt.Printf("被杠分:%+v\n", beMingKong)
	fmt.Printf("暗杠分:%+v\n", anKong)
	fmt.Printf("补杠分:%+v\n", buKong)
	fmt.Printf("总分:%+v\n", total)
	//打印牌型
	for k, v := range HUS {
		var handCards []byte //手牌
		if v_hs, ok := HANDS[k]; ok {
			handCards = v_hs
		}
		if len(handCards) == 0 {
			continue
		}
		if v&PAOHU > 0 {
			v_h := DiscardHu(handCards) //胡
			v_h |= heType(v_h, k, handCards) //大七对,清一色,字一色
			fmt.Printf("放炮胡牌,位置:%d, 牌型:%+v\n", k, HuTypeStr(v_h))
		} else if v&ZIMO > 0 {
			v_h := DrawDetect(handCards) //胡
			v_h |= heType(v_h, k, handCards) //大七对,清一色,字一色
			fmt.Printf("自摸胡牌,位置:%d, 牌型:%+v\n", k, HuTypeStr(v_h))
		}
	}
}

/*
庄自摸	    庄炮胡	        闲自摸	            闲炮胡(闲放)	闲炮胡(庄放)
牌型*庄*3	牌型*放炮*庄	牌型+牌型+牌型*庄	牌型*放炮	    牌型*庄*放炮
6*牌型	    4*牌型	        4*牌型	            2*牌型	        4*牌型
庄家自摸时,其他玩家需每人多支付2倍,如:庄家平胡自摸所得番数为:1*2*3
庄家炮胡时,放炮者需支付给庄家2倍,其他玩家不需要给,如:庄家平胡炮胡所得番数为:1*2*2
其他玩家自摸时,庄家需要额外支付2倍,如:闲家平胡自摸所得番数为:1+1+1*2
其他玩家炮胡时,如果是闲家放炮,则庄家不需要给番数,所得番数为:1*2
其他玩家炮胡时,如果是庄家放炮,则所得番数为:1*2*2
*/
//结算,(明杠,放冲,庄家 - 收一家)
func gameOver() (huFan, mingKong, beMingKong, anKong, buKong, total map[uint32]int32) {
	huFan      = make(map[uint32]int32)// 胡牌牌型番数
	mingKong   = make(map[uint32]int32)// 闷豆的番数
	beMingKong = make(map[uint32]int32)// 被点豆的负番数
	anKong     = make(map[uint32]int32)// 明豆的番数
	buKong     = make(map[uint32]int32)// 拐弯豆的番数
	total      = make(map[uint32]int32)// 总番数
	var k uint32
	for k = 1; k <= 4; k++ {
		//牌型分
		if v, ok := HUS[k]; ok {
			var handCards []byte //手牌
			if v_hs, ok := HANDS[k]; ok {
				handCards = v_hs
			}
			f_t := HuType(v, handCards) //胡牌牌型,多个牌型时相乘
			f_w := HuWay(v)             //胡牌方式,多个方式时相乘
			f_tw := f_t * f_w                //牌型分
			//牌型分,t.seat=出牌(放冲)位置
			huFan = fanType(DEALER, PAOSEAT, k, f_tw, huFan)
		}
		//杠牌分
		var ks []uint32 //杠牌数据
		if v_ks, ok := KONGS[k]; ok {
			ks = v_ks
		}
		for _, v := range ks {  //杠牌分
			i, _, cy := DecodeKong(v) //解码杠值
			f_k := HuKong(cy) //杠
			if cy == MING_KONG {
				mingKong[k] += f_k       //收一家
				beMingKong[i] += 0 - f_k //被收一家
			} else if cy == BU_KONG {
				buKong = over3(buKong, k, f_k) //收三家
			} else if cy == AN_KONG {
				anKong = over3(anKong, k, f_k) //收三家
			}
		}
	}
	//总番数
	for k = 1; k <= 4; k++ {
		total[k] += huFan[k] + mingKong[k] + beMingKong[k] + anKong[k] + buKong[k]
	}
	return huFan, mingKong, beMingKong, anKong, buKong, total
}

//倍数(放冲,庄家 - 双倍) TODO:优化
//dealer=庄家位置,paoseat=放炮位置,huseat=胡牌位置
func fanType(dealer, paoseat, huseat uint32, f_tw int32,
hf map[uint32]int32) map[uint32]int32 {
	if paoseat == huseat { //自摸,收三家
		if huseat == dealer { //庄家自摸
			//6 //其它3家*2 (庄家*2倍)
			hf = over3(hf, dealer, f_tw * 2) //收三家
		} else { //闲家自摸
			//4 //庄家*2 + 其它2家*1
			var i uint32
			for i = 1; i <= 4; i++ {
				if i == huseat {
					continue
				}
				if i == dealer {
					hf = over1(hf, huseat, i, f_tw * 2) //收一家
				} else {
					hf = over1(hf, huseat, i, f_tw * 1) //收一家
				}
			}
		}
	} else { //炮胡,收一家
		if huseat == dealer { //庄家胡(肯定闲家放炮)
			//4//收放炮的*4 (放2倍*庄2倍)
			hf = over1(hf, huseat, paoseat, f_tw * 4) //收一家
		} else { //闲家胡
			if paoseat == dealer { //庄家放炮
				//4//收庄家的*4 (放炮2倍*庄家2倍)
				hf = over1(hf, huseat, dealer, f_tw * 4) //收一家
			} else { //闲家放炮
				//2//收闲家的*2 (放炮2倍)
				hf = over1(hf, huseat, paoseat, f_tw * 2) //收一家
			}
		}
	}
	return hf
}

//收三家,seat=收的位置,val=收的番数
func over3(m map[uint32]int32, seat uint32, val int32) map[uint32]int32 {
	var i uint32
	for i = 1; i <= 4; i++ {
		var value int32
		if i != seat {
			value = 0 - val //为负数
		} else {
			value = 3 * val //收三家
		}
		m[i] += value
	}
	return m
}

//收一家,s1=收的位置,s2=出的位置,val=收的番数
func over1(m map[uint32]int32, s1, s2 uint32, val int32) map[uint32]int32 {
	m[s1] += val
	m[s2] -= val
	return m
}

//-------------------

//牌型
func HuTypeStr(v uint32) string {
	var s string = ""
	if v&ZIMO > 0 {
		s += HU2STR[ZIMO] + ","
	}
	if v&PAOHU > 0 {
		s += HU2STR[PAOHU] + ","
	}
	if v&HU_PING > 0 {
		s += HU2STR[HU_PING] + ","
	}
	if v&HU_SINGLE > 0 {
		s += HU2STR[HU_SINGLE] + ","
	}
	if v&HU_SINGLE_ZI > 0 {
		s += HU2STR[HU_SINGLE_ZI] + ","
	}
	if v&HU_SEVEN_PAIR_BIG > 0 {
		s += HU2STR[HU_SEVEN_PAIR_BIG] + ","
	}
	if v&HU_SEVEN_PAIR > 0 {
		s += HU2STR[HU_SEVEN_PAIR] + ","
	}
	if v&HU_SEVEN_PAIR_KONG > 0 {
		s += HU2STR[HU_SEVEN_PAIR_KONG] + ","
	}
	if v&HU_ONE_SUIT > 0 {
		s += HU2STR[HU_ONE_SUIT] + ","
	}
	if v&HU_ALL_ZI > 0 {
		s += HU2STR[HU_ALL_ZI] + ","
	}
	return s
}

// 打牌检测,胡牌, 放炮胡检测
func DiscardHu(cs []byte) uint32 {
	cards := make([]byte, len(cs))
	copy(cards, cs)
	//胡
	var status uint32 = existHu(cards)
	if status > 0 {
		if status&HU_PING > 0 {
			status ^= HU_PING //去掉平胡
		}
		if status&HU_DANDIAO > 0 {
			status ^= HU_DANDIAO //放炮时不算单钓分数
		}
		status |= PAOHU
	}
	return status
}

// 摸牌检测,胡牌
func DrawDetect(cs []byte) uint32 {
	le := len(cs)
	cards := make([]byte, le)
	copy(cards, cs)
	var status uint32
	//自摸胡检测
	status = existHu(cards)
	if status > 0 {
		if status&HU_PING > 0 {
			status ^= HU_PING //去掉平胡
		}
		status |= ZIMO
	}
	if status > 0 && le == 14 {
		status |= HU_MENGQING
	}
	return status
}

// 胡牌牌型检测
func HuTypeDetect(hu, chow, kong bool, cs []byte) uint32 {
	Sort(cs, 0, len(cs)-1)
	return existHuType(hu, chow, kong, cs)
}

//大七对,清一色,字一色
func heType(val, seat uint32, hands []byte) uint32 {
	var cs []byte = []byte{}
	cs = append(cs, hands...)
	var c_s []uint32 //
	if v_cs, ok := CHOWS[seat]; ok {
		c_s = v_cs
	}
	for _, v1 := range c_s { //吃
		c1, c2, c3 := DecodeChow(v1) //解码
		cs = append(cs, c1, c2, c3)
	}
	var k_s []uint32 //
	if v_ks, ok := KONGS[seat]; ok {
		k_s = v_ks
	}
	for _, v2 := range k_s { //杠
		_, c, _ := DecodeKong(v2) //解码
		cs = append(cs, c, c, c, c)
	}
	var p_s []uint32 //
	if v_ps, ok := KONGS[seat]; ok {
		p_s = v_ps
	}
	for _, v3 := range p_s { //碰
		_, c := DecodePeng(v3) //解码
		cs = append(cs, c, c, c)
	}
	var hu bool = false //是否有胡牌牌型
	if val&HU > 0 {
		hu = true //有牌型
	}
	var chow bool = len(c_s) != 0 //有吃
	var kong bool = KongDetect(hands) //手牌中是否有杠
	//glog.Infof("cs -> %+x, chow -> %v", cs, chow)
	return HuTypeDetect(hu, chow, kong, cs)
}

// 检测手牌是否有杠
func KongDetect(cs []byte) bool {
	var n uint32 = havekongs(cs)
	return n > 0
}

// 碰杠吃数据
func EncodePeng(seat uint32, card byte) uint32 {
	seat = seat << 8
	seat |= uint32(card)
	return seat
}

func DecodePeng(value uint32) (seat uint32, card byte) {
	seat = value >> 8
	card = byte(value & 0xFF)
	return
}

func EncodeKong(seat uint32, card byte, value uint32) uint32 {
	value = value << 16
	value |= (seat << 8)
	value |= uint32(card)
	return value
}

func DecodeKong(value uint32) (seat uint32, card byte, v uint32) {
	v = value >> 16
	seat = (value >> 8) & 0xFF
	card = byte(value & 0xFF)
	return
}

func EncodeChow(c1, c2, c3 byte) (value uint32) {
	value =  uint32(c1) << 16
	value |= uint32(c2) << 8
	value |= uint32(c3)
	return
}

func DecodeChow(value uint32) (c1, c2, c3 byte) {
	c1 = byte(value >> 16)
	c2 = byte(value >> 8 & 0xFF)
	c3 = byte(value & 0xFF)
	return
}

func DecodeChow2(value uint32) (c1, c2 byte) {
	c1 = byte(value >> 8)
	c2 = byte(value & 0xFF)
	return
}

// 对牌值从小到大排序，采用快速排序算法
func Sort(arr []byte, start, end int) {
	if start < end {
		i, j := start, end
		key := arr[(start+end)/2]
		for i <= j {
			for arr[i] < key {
				i++
			}
			for arr[j] > key {
				j--
			}
			if i <= j {
				arr[i], arr[j] = arr[j], arr[i]
				i++
				j--
			}
		}
		if start < j {
			Sort(arr, start, j)
		}
		if end > i {
			Sort(arr, i, end)
		}
	}
}

//算番(牌型) TODO:优化
func HuType(v uint32, cs []byte) int32 {
	var f int32 = 1
	if v&HU_PING > 0 {
		f *= FAN[HU_PING]
	}
	if v&HU_SINGLE > 0 {
		f *= FAN[HU_SINGLE]
	}
	if v&HU_SINGLE_ZI > 0 {
		f *= FAN[HU_SINGLE_ZI]
	}
	if v&HU_SEVEN_PAIR_BIG > 0 {
		f *= FAN[HU_SEVEN_PAIR_BIG]
	}
	if v&HU_SEVEN_PAIR > 0 {
		f *= FAN[HU_SEVEN_PAIR]
	}
	if v&HU_SEVEN_PAIR_KONG > 0 {
		var n uint32 = havekongs(cs)
		f *= FAN[HU_SEVEN_PAIR_KONG] * (1 << n) //n=手牌中暗杠数:4*2^n
	}
	if v&HU_ONE_SUIT > 0 {
		f *= FAN[HU_ONE_SUIT]
	}
	if v&HU_ALL_ZI > 0 { //没牌型6番,有12番
		if v == (HU|HU_ALL_ZI|ZIMO) ||
		v == (HU|HU_ALL_ZI|PAOHU) {
			f *= FAN[HU_ALL_ZI] / 2
		} else {
			f *= FAN[HU_ALL_ZI]
		}
	}
	return f
}

//算番(胡牌方式) TODO:优化
func HuWay(v uint32) int32 {
	var f int32 = 1
	if v&QIANG_GANG > 0 {
		f *= FAN[QIANG_GANG]
	}
	if v&HU_KONG_FLOWER > 0 {
		f *= FAN[HU_KONG_FLOWER]
	}
	if v&HU_MENGQING > 0 {
		f *= FAN[HU_MENGQING]
	}
	if v&HU_DANDIAO > 0 {
		f *= FAN[HU_DANDIAO]
	}
	if v&TIAN_HU > 0 {
		f *= FAN[TIAN_HU]
	}
	if v&DI_HU > 0 {
		f *= FAN[DI_HU]
	}
	return f
}

//算番(杠)
func HuKong(v uint32) int32 {
	return FAN[v]
}

// 手牌有多少个杠牌
func havekongs(cards []byte) uint32 {
	var i uint32 = 0
	var m = make(map[byte]int)
	for _, v := range cards {
		if n, ok := m[v]; ok {
			m[v] = n + 1
		} else {
			m[v] = 1
		}
	}
	for _, v := range m {
		if v == 4 {
			i += 1
		}
	}
	return i
}

// 胡牌后检测
// 清一色(同一花色,必须有牌型)/字一色(全部是字牌)
// 大七对, 1对子(将)+(刻子,杠) + 不能吃 + 杠(非手牌中杠)
// 有序slices cs 包含吃，碰，杠, chow=false没有吃,=true有吃
// kong=false没有杠,=true有杠, hu=false没有牌型,=true有牌型
func existHuType(hu, chow, kong bool, cs []byte) uint32 {
	var all_zi bool = true
	var one_suit bool = true
	var seven_pair_big bool = true
	var b bool
	var c byte
	var huType uint32
	var m = make(map[byte]int)
	for _, v := range cs {
		if n, ok := m[v]; ok {
			m[v] = n + 1
		} else {
			m[v] = 1
		}
		if !one_suit && !all_zi {
			continue
		}
		if c == 0 {
			c = v
		} else if c >> 4 != v >> 4 {
			one_suit = false
		} else if uint32(v >> 4) >= FENG {
			one_suit = false
		}
		if uint32(v >> 4) < FENG {
			all_zi = false
		}
	}
	for _, v := range m {
		if v == 2 && !b {
			b = true
			continue
		}
		if v < 3 || chow || kong {
			seven_pair_big = false
			break
		}
	}
	if seven_pair_big {
		huType |= HU
		huType |= HU_SEVEN_PAIR_BIG
	}
	// (huType > 0 || hu) //有牌型
	if one_suit && (huType > 0 || hu) {
		huType |= HU
		huType |= HU_ONE_SUIT
	}
	if all_zi {
		huType |= HU
		huType |= HU_ALL_ZI
	}
	return huType
}

// 13烂(数牌相隔2个或以上,字牌不重复)/七星13烂(13烂基础上包含7个不同字牌)
// 有序slices cs && len(cs) == 14
func existThirteen(cs []byte) uint32 {
	var thirteen_single bool = true
	var thirteen_single_zi bool = true
	var c byte
	var n int
	var huType uint32
	for _, v := range cs {
		if uint32(v >> 4) < FENG {
			if c == 0 || c >> 4 != v >> 4 {
				c = v
				continue
			}
			if v - c < 0x03 {
				thirteen_single = false
				break
			}
			c = v
		} else {
			n++
			if v == c {
				thirteen_single_zi = false
				break
			}
			c = v
		}
	}
	if thirteen_single_zi && thirteen_single_zi && n == 7 {
		huType |= HU_SINGLE_ZI
	} else if thirteen_single && thirteen_single_zi {
		huType |= HU_SINGLE
	}
	return huType
}

// 判断是否小七对(7个对子),豪华小七对(7个对子,其中有杠)
// 有序slices cs && len(cs) == 14
func exist7pair(cs []byte) uint32 {
	var seven_pair bool = true
	var c byte
	var i int
	var huType uint32
	var le int = len(cs)
	for n := 0; n < le-1; n += 2 {
		if cs[n] != cs[n+1] {
			seven_pair = false
			break
		}
		if i != 4 && cs[n] == c {
			i += 2
		} else if i != 4 {
			c = cs[n]
			i = 2
		}
	}
	if seven_pair && i == 4 {
		huType |= HU_SEVEN_PAIR_KONG
	} else if seven_pair {
		huType |= HU_SEVEN_PAIR
	}
	return huType
}

// 判断是否胡牌,0表示不胡牌,非0用32位表示不同的胡牌牌型
func existHu(cs []byte) uint32 {
	le := len(cs)
	var huType uint32
	//单钓胡牌
	if le == 2 && cs[0] == cs[1] {
		huType |= HU
		huType |= HU_DANDIAO
		return huType
	}
	//排序slices
	Sort(cs, 0, len(cs)-1)
	// 14张牌型胡牌
	if le == 14 {
		// 七小对牌型胡牌
		t1 := exist7pair(cs)
		if t1 > 0 {
			huType |= HU
			huType |= t1
			return huType
		}
		// 十三烂牌型胡牌
		t2 := existThirteen(cs)
		if t2 > 0 {
			huType |= HU
			huType |= t2
			return huType
		}
	}
	// 3n +2 牌型胡牌
	if (le - 2) % 3 != 0 {
		return huType
	}
	return existHu3n2(cs, le)
}

// 3n +2 牌型胡牌
// 有序slices cs
func existHu3n2(cs []byte, le int) (huType uint32) {
	for n := 0; n < le-1; n++ {
		if cs[n] == cs[n+1] { //
			list := make([]byte, le)
			copy(list, cs)
			list[n] = 0x00
			list[n+1] = 0x00
			for i := 0; i < le-2; i++ {
				if list[i] > 0 {
					for j := i + 1; j < le-1; j++ {
						if list[j] > 0 && list[i] > 0 {
							for k := j + 1; k < le; k++ {
								if list[k] > 0 && list[i] > 0 && list[j] > 0 {
									//刻子
									if list[i] == list[j] && list[j] == list[k] {
										list[i], list[j], list[k] = 0x00, 0x00, 0x00
										break
									}
									//顺子
									if list[i]+1 == list[j] && list[j]+1 == list[k] {
										list[i], list[j], list[k] = 0x00, 0x00, 0x00
										break
									}
									if list[i] >= 0x41 && list[k] <= 0x44 &&
									list[i] != list[j] && list[j] != list[k] {
										list[i], list[j], list[k] = 0x00, 0x00, 0x00
										break
									}
								}
							}
						}
					}
				}
			}
			num := false
			for i := 0; i < le; i++ {
				if list[i] > 0 {
					num = true
					break
				}
			}
			if !num {
				huType |= HU
				huType |= HU_PING
				return huType
			}
		}
	}
	return huType
}
