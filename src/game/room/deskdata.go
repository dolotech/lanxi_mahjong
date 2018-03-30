package room

import (
	"game/interfacer"
	"sync"
	"time"
	"encoding/json"
)

//房间牌桌数据结构
type Desk struct {
	id        uint32                        //房间id
	data      *DeskData                     //房间类型基础数据
	dealer    uint32                        //庄家的座位
	dice      uint32                        //骰子
	cards     []byte                        //没摸起的海底牌
	players   map[uint32]interfacer.IPlayer //房间玩家
	lianCount uint32                        // 连庄数
	round     uint32                        // 打牌局数
	vote      uint32                        //投票发起者座位号
	votes     map[uint32]uint32             //投票同意解散de玩家座位号列表
	voteT     *time.Timer                   //投票定时器
	skip      map[uint32]uint32             //过圈
	opt       map[uint32]int64             // 玩家吃碰杠胡掩码
	luckyCard byte                          //本局的财神牌值
	discard   byte                          //出牌
	draw      byte                          //摸牌
	qiangKongCard      byte                          // 抢杠牌
	seat      uint32                        //当前摸牌|出牌位置
	operate   int                           //操作状态
	sync.RWMutex                            //房间锁
	state     bool                          //房间状态
	kong      bool                          //是否杠牌出牌
	closeCh   chan bool                     //关闭通道
	ready     map[uint32]bool               //是否准备
	outCards  map[uint32][]byte             //海底牌
	pongCards map[uint32][]uint32           //碰牌
	kongCards map[uint32][]uint32           //杠牌
	chowCards map[uint32][]uint32           //吃牌(8bit-8-8)
	handCards map[uint32][]byte             //手牌
	maizi     map[uint32]uint32             //买子
	offline   map[uint32]bool               // true:离线,false:上线

	cheatHandCards [][]byte
	cheatLeftCards []byte
	cheatwildcard  byte
}

type RoomInfoResp struct {
	RoomId    string      `json:"roomid"`                 // 房间号(邀请码)
	Ids       map[uint32]string     `json:"ids"`          // 玩家的id
	Roomcards map[uint32]uint32     `json:"roomcards"`    // 玩家的房卡
	OutCards  map[uint32][]uint32     `json:"outCards"`   //海底牌
	PongCards map[uint32][]uint32      `json:"pongCards"` //碰牌
	KongCards map[uint32][]uint32      `json:"kongCards"` //杠牌
	ChowCards map[uint32][]uint32      `json:"chowCards"` //吃牌(8bit-8-8)
	HandCards map[uint32][]uint32      `json:"handCards"` //手牌
	Maizi     map[uint32]uint32      `json:"maizi"`       //买子
	Offline   map[uint32]bool     `json:"offline"`        // true:离线,false:上线

	Data      *DeskData                 `json:"data"`            //房间类型基础数据
	Dealer    uint32                       `json:"dealer"`       //庄家的座位
	Dice      uint32                       `json:"dice"`         //骰子
	Cards     []uint32                       `json:"cards"`      //没摸起的海底牌
	LianCount uint32                         `json:"lianCount"`  // 连庄数
	Round     uint32                        `json:"round"`       // 打牌局数
	Opt       map[uint32]int64             `json:"opt"`         // 玩家吃碰杠胡掩码
	LuckyCard uint32                          `json:"luckyCard"` //本局的财神牌值
	Discard   uint32                           `json:"discard"`  //出牌
	Draw      uint32                           `json:"draw"`     //摸牌
	Seat      uint32                         `json:"seat"`       //当前摸牌|出牌位置
	Operate   int                           `json:"operate"`     //操作状态
}

func (t *Desk) ToString() string {
	t.Lock()
	defer t.Unlock()
	roomInfoResp := &RoomInfoResp{
		RoomId:    t.data.Code,
		Data:      t.data,
		PongCards: t.pongCards,
		KongCards: t.kongCards,
		ChowCards: t.chowCards,
		Maizi:     t.maizi,
		Offline:   t.offline,
		Dealer:    t.dealer,
		Dice:      t.dice,
		LianCount: t.lianCount,
		Round:     t.round,
		Opt:       t.opt,
		LuckyCard: uint32(t.luckyCard),
		Discard:   uint32(t.discard),
		Draw:      uint32(t.draw),
		Seat:      t.seat,
		Operate:   t.operate,
	}
	for _, v := range t.cards {
		roomInfoResp.Cards = append(roomInfoResp.Cards, uint32(v))
	}

	roomInfoResp.HandCards = make(map[uint32][]uint32, 4)
	for k, v := range t.handCards {
		for _, card := range v {
			roomInfoResp.HandCards[k] = append(roomInfoResp.HandCards[k], uint32(card))
		}
	}

	roomInfoResp.OutCards = make(map[uint32][]uint32, 4)
	for k, v := range t.outCards {
		for _, card := range v {
			roomInfoResp.OutCards[k] = append(roomInfoResp.OutCards[k], uint32(card))
		}
	}

	roomInfoResp.Ids = make(map[uint32]string, 4)
	for k, v := range t.players {
		roomInfoResp.Ids[k] = v.GetUserid()
	}

	roomInfoResp.Roomcards = make(map[uint32]uint32, 4)
	for k, v := range t.players {
		roomInfoResp.Roomcards[k] = v.GetRoomCard()
	}

	str, err := json.Marshal(roomInfoResp)
	if err != nil {
		return ""
	}
	return string(str)
}
