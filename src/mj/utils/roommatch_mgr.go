package utils

import (
	"math/rand"
	"time"

	"github.com/88bytes/nano/session"
)

type (
	// RoomMatchMgr 用来管理私人房匹配的数据
	RoomMatchMgr struct {
		currentRoomID int
		roomInfos     map[int]*RoomInfo
		rand          *rand.Rand
	}

	// RoomInfo 里面装了私人房的匹配信息
	RoomInfo struct {
		roomID       int
		sessions     []*session.Session
		pvpRoomInfos map[int]*StartBattleRoomInfo
	}
)

var (
	// RoomMatchMgrInst 是用来处理玩家自由匹配的
	RoomMatchMgrInst *RoomMatchMgr
)

const (
	// RoomMatchMaxRoomID 是自由匹配的最大房间ID，非自由匹配用6位的房间号码
	RoomMatchMaxRoomID = 99999
)

// NewRoomMatchMgr 会创建一个FreeMatchMgr出来
func NewRoomMatchMgr() *RoomMatchMgr {
	mgr := new(RoomMatchMgr)
	mgr.currentRoomID = 0
	mgr.roomInfos = make(map[int]*RoomInfo)
	source := rand.NewSource(time.Now().UnixNano())
	mgr.rand = rand.New(source)
	return mgr
}
