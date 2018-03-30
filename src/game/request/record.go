package request

import (
	"lib/socket"
	"game/data"

	"game/interfacer"
	"sort"
	"protocol"

	"code.google.com/p/goprotobuf/proto"
)

func init() {

	socket.Regist(&protocol.CPrivateRecord{}, records)

	socket.Regist(&protocol.CPRecordByRid{}, getrecordsbyroomid)
}

// get private room record by room id   TODO:为什么要去数据库里面取?
func getrecordsbyroomid(ctos *protocol.CPRecordByRid, c interfacer.IConn) {
	stoc := &protocol.SPRecordByRid{}

	if ctos.GetRid() == 0 {
		stoc.Error = proto.Uint32(uint32(protocol.Error_PrivateRecordEmpty))
	} else {
		d := data.GameOverRecord{RoomId: ctos.GetRid()}

		if err := d.Get(c.GetUserid()); err != nil {
			stoc.Error = proto.Uint32(uint32(protocol.Error_PrivateRecordEmpty))
			//glog.Errorln(err)
		} else {
			stoc.Rname = proto.String("")
			stoc.Invitecode = &d.Invitecode
			//stoc.Ma = &d.MaCount
			stoc.TotalRound = &d.TotalRound
			stoc.Rid = &d.RoomId
			stoc.Ante = &d.Ante
			stoc.Time = &d.Ctime

			coinMap := make(map[string]int32)

			for _, v := range d.Rounds {
				for _, v1 := range v.Users {
					//coinMap[v1.Userid] = v1.Coin
					coinMap[v1.Userid] += v1.Coin
				}
			}

			for _, userRecord := range d.Rounds[0].Users {
				user := &data.User{Userid: userRecord.Userid}
				user.Get()
				var coin int32 = coinMap[userRecord.Userid]
				l := &protocol.PrivateRecordDetails{
					Seat:   &userRecord.Seat,
					Userid: &userRecord.Userid,
					Sex:    &user.Sex,
					Photo:  &user.Photo,
					Uname:  &user.Nickname,
					Coin:   &coin,
				}
				stoc.List = append(stoc.List, l)
			}
		}
	}

	c.Send(stoc)
}

// 私人局记录
func records(ctos *protocol.CPrivateRecord, c interfacer.IConn) {
	stoc := &protocol.SPrivateRecord{}

	list := data.GameOverRecords{}
	err := list.Get(c.GetUserid(), 1, 100)

	if err != nil || len(list) == 0 {
		stoc.Error = proto.Uint32(uint32(protocol.Error_PrivateRecordEmpty))
	} else {
		users := make(map[string]*data.User)
		for _, d := range list {

			coinMap := make(map[string]int32)
			p := &protocol.PrivateRecord{
				Rname:      proto.String(""),
				Invitecode: &d.Invitecode,
				TotalRound: &d.TotalRound,
				Rid:        proto.Uint32(d.RoomId),
				Ante:       &d.Ante,
				Time:       &d.Ctime,
				Round:      &d.TotalRound,
			}

			//展示第一局
			for _, userRecord := range d.Rounds {
				l := &protocol.PrivateRecords{
					Round: proto.Uint32(userRecord.Round),
				}
				sort.Sort(StructSlice(userRecord.Users))
				for _, userRecords := range userRecord.Users {
					var name string
					if user, ok := users[userRecords.Userid]; ok {
						name = user.Nickname
					} else {
						user := &data.User{Userid: userRecords.Userid}
						err := user.GetPhotoSexName()
						if err == nil {
							users[userRecords.Userid] = user
							name = user.Nickname
						}
					}
					h := &protocol.PrivateDetails{
						Uname: proto.String(name),
						Coin:  proto.Int32(userRecords.Coin),
					}
					l.List = append(l.List, h)
					coinMap[userRecords.Userid] += userRecords.Coin
				}
				p.Lists = append(p.Lists, l)
			}

			for _, userRecord := range d.Rounds[0].Users {
				l := &protocol.PrivateRecordDetails{
					Seat:   &userRecord.Seat,
					Userid: &userRecord.Userid,
					Coin:   proto.Int32(coinMap[userRecord.Userid]),
				}
				var user *data.User
				ok := false
				if user, ok = users[userRecord.Userid]; ok {
					l.Sex = &user.Sex
					l.Photo = &user.Photo
					l.Uname = &user.Nickname
				} else {
					user := &data.User{Userid: userRecord.Userid}
					err := user.GetPhotoSexName()
					if err == nil {
						users[userRecord.Userid] = user
					}
				}
				if user != nil {
					l.Sex = &user.Sex
					l.Photo = &user.Photo
					l.Uname = &user.Nickname
				}
				p.List = append(p.List, l)
			}

			p.Coin = proto.Int32(coinMap[c.GetUserid()])
			stoc.List = append(stoc.List, p)
		}
	}
	c.Send(stoc)
}

type StructSlice []*data.GameOverUserRecord

func (b StructSlice) Len() int {
	return len(b)
}

func (b StructSlice) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b StructSlice) Less(i, j int) bool {
	return b[i].Seat < b[j].Seat
}
