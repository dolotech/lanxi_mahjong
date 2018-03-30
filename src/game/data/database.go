package data

import (
	"config"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
	"time"
)

const (
	_DBNAME = "lanxi_db"

	LIMIT_MAX = 200 // 每页最大记录数量
	LIMIT_MIN = 20  // 每页最小记录数量
)

// 所以集合名字
const (
	_USER                 = "user"             // 用户集合
	_GEN_USER_ID          = "user_id_gen"      // 玩家ID自增
	_GEN_ROOM_ID          = "room_id_gen"      // 房间ID自增
	_GAMEOVER_RECORD      = "gameover_record"  // 私人局结算记录用于前端
	_CARD_RECORD          = "card_record"      // 打牌记录,以房间ID+当前局数组成字符串作为_id,TODO:_id这样组合会出现覆盖
	_GAMEOVER_PRIVATE     = "gameover_private" // 私人
	_RESOURCE_RECORD      = "resource_record"  //资源消耗记录
	_STATISTICS           = "statistics"       // 统计在线相关数据
	_ROOM_CREATE_RECORD   = "create_room"      // 开房记录
	_LOGIN_LOG_OUT_RECORD = "login_log"        // 玩家登录登出的record
)

var session *mgo.Session

func C(collection string) *mgo.Collection {
	if session == nil {
		var err error
		session, err = mgo.Dial(config.Opts().Db_addr)

		//defer session.Close()
		if err != nil {
			glog.Fatalln(err, config.Opts().Db_addr)
		}
		//session.SetPoolLimit(9)
		go func() {
			for {
				time.Sleep(time.Second * 9)
				err := session.Ping()
				if err != nil {
					glog.Errorln(err)
					session.Refresh()
				}
			}

		}()
	}
	//	ses := session.Clone()
	//	defer ses.Close()
	return session.DB(_DBNAME).C(collection)
}
