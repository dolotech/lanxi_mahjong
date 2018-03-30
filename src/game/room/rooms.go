package room

import (
	"game/interfacer"
	"lib/utils"
	"strconv"
	"sync"
	"time"

	"github.com/golang/glog"
)

//全局
var rooms *SocialRooms

//牌桌列表结构
type SocialRooms struct {
	sync.RWMutex                             //读写锁
	list         map[string]interfacer.IDesk //牌桌列表
	closeCh      chan bool                   //关闭通道
}

//初始化列表
func init() {
	rooms = &SocialRooms{
		list:    make(map[string]interfacer.IDesk),
		closeCh: make(chan bool, 1),
	}
	go rooms_ticker() //goroutine,定时清理
}

//生成一个牌桌邀请码,全列表中唯一
func GenInvitecode(count int) (s string) {
	if count > 0 {
		count--
		s = strconv.Itoa(int(utils.RandInt32N(900000)) + 100000)
		if Exist(s) { //是否已经存在
			return GenInvitecode(count) //重复尝试,TODO:一定次数后放弃尝试
		}
	}
	return s
}

//添加一个私人局房间
func Add(key string, rdata interfacer.IDesk) {
	rooms.Lock()
	defer rooms.Unlock()
	rooms.list[key] = rdata
}

//删除一个牌桌
func Del(key string) {
	rooms.Lock()
	defer rooms.Unlock()
	delete(rooms.list, key)
}

//获取牌桌接口
func Get(key string) interfacer.IDesk {
	rooms.RLock()
	defer rooms.RUnlock()
	if room, ok := rooms.list[key]; ok {
		return room
	}
	return nil
}

//是否存在
func Exist(key string) bool {
	rooms.RLock()
	defer rooms.RUnlock()
	_, ok := rooms.list[key]
	return ok
}

//关闭列表
func Close() {
	rooms.Lock()
	defer rooms.Unlock()
	close(rooms.closeCh) //关闭
}

//计时器
func rooms_ticker() {
	tick := time.Tick(10 * time.Minute) //每局限制了最高10分钟
	glog.Infof("rooms ticker started -> %d", 1)
	for {
		select {
		case <-tick:
			//逻辑处理
			rooms_expire(false) //过期清理
		case <-rooms.closeCh:
			glog.Infof("rooms close -> %d", len(rooms.list))
			//TODO:rooms_expire(true) //强制关闭
			return
		}
	}
}

//定期清理or关闭清除
func rooms_expire(ok bool) {
	rooms.Lock()
	defer rooms.Unlock()
	for _, r := range rooms.list {
		go func(room interfacer.IDesk) {
			room.Closed(ok) //goroutine处理,避免删除时死锁
		}(r)
	}
}

func Len() int {
	rooms.Lock()
	defer rooms.Unlock()
	return len(rooms.list)
}

//func LenInPlaying() int {
//	rooms.Lock()
//	defer rooms.Unlock()
//	var count int
//	for _, v := range rooms.list {
//		if ok := v.IsExpire(); ok {
//			count++
//		}
//	}
//	return count
//}
