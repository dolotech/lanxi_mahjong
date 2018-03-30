package resource

import (
	"game/data"
	"game/interfacer"
	"protocol"

	"code.google.com/p/goprotobuf/proto"
	"github.com/golang/glog"
	"game/players"
	"errors"
)

const (
	ROOM_CARD uint32 = 4
)

var RES_HASH = map[uint32]string{
	4:   "RoomCard",
}

func NotifyChangeRes(userid string, id uint32, count int32) error {
	player := players.Get(userid)
	if player != nil {
		if id == ROOM_CARD {
			if count < 0 {
				count = 0
			}
			player.SetRoomCard(uint32(count))
		}

		stoc := &protocol.SResource{}
		stoc.List = append(stoc.List, &protocol.ResVo{Id: &id, Count: &count})
		player.Send(stoc)
		return nil
	}
	return errors.New("player offline!")
}

// 更改单个资源
func ChangeRes(c interfacer.IPlayer, id uint32, count int32, Type uint32) error {
	m := make(map[uint32]int32)
	m[id] = count
	return ChangeMulti(c, m, Type)
}

// 更改多个资源
func ChangeMulti(userdata interfacer.IPlayer, res map[uint32]int32, Type uint32) error {
	list := make([]*protocol.ResVo, 0, len(res))
	var record data.DataResChanges
	updataValue := map[string]int32{}
	for id, count := range res {
		var current int32
		switch id {

		case ROOM_CARD:
			current = int32(userdata.GetRoomCard()) + count
			if current < 0 {
				current = 0
			}
			userdata.SetRoomCard(uint32(current))
		}

		if key, ok := RES_HASH[id]; ok {
			updataValue[key] = current
		}

		list = append(list, &protocol.ResVo{
			Id:    proto.Uint32(id),
			Count: proto.Int32(current),
		})

		ch := &data.DataResChange{
			Kind:     id,
			Channel:  Type,
			Count:    count,
			Residual: uint32(current),
		}
		record = append(record, ch)
	}
	var err error
	user := &data.User{Userid: userdata.GetUserid()}
	if len(updataValue) > 0 {
		err = user.UpdateResource(updataValue)
		if err != nil {
			glog.Errorln(err)
			return err
		}
	}

	err = record.Save(userdata.GetUserid())
	if err != nil {
		glog.Errorln(err)
	}

	stoc := &protocol.SResource{List: list}
	userdata.Send(stoc)

	return err
}
