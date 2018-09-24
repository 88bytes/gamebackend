package utils

import (
	"fmt"
	"math/rand"
	"time"

	"mj/player_position"

	"github.com/88bytes/nano"
	"github.com/88bytes/nano/session"
)

type (
	// FreeMatchMgr 用来管理自由匹配的数据，还有定时器
	// pvpRoomInfos 这个 map 要注意没有用的要回收，要不然内存就会越来越大
	FreeMatchMgr struct {
		matchMgrTool  *MatchMgrTool
		currentRoomID int
		playerCount   int
		sessions      []*session.Session
		timerInst     *nano.Timer
		rand          *rand.Rand
		pvpRoomInfos  map[int]*StartBattleRoomInfo
	}

	// StartBattleRoomInfo contains pvp room info
	StartBattleRoomInfo struct {
		RoomID          int              `json:"RoomID"`
		StartBattleInfo *StartBattleInfo `json:"StartBattleInfo"`
	}

	// StartBattleInfo contains battle info
	StartBattleInfo struct {
		RoomtID         int                     `json:"RoomID"`
		RandomSeed      int                     `json:"RandomSeed"`
		BankerPosition  int                     `json:"BankerPosition"`
		SittingPosition int                     `json:"SittingPosition"`
		AIOwnerPosition int                     `json:"AIOwnerPosition"`
		MaxBattleCount  int                     `json:"MaxBattleCount"`
		ZhuaNiaoCount   int                     `json:"ZhuaNiaoCount"`
		PlayerInfos     []StartBattlePlayerInfo `json:"PlayerInfos"`
	}

	// StartBattlePlayerInfo contains a battlePlayerInfo
	StartBattlePlayerInfo struct {
		UID          uint   `json:"uid"`
		NickName     string `json:"NickName"`
		ControlType  int    `json:"ControlType"`
		IsBanker     bool   `json:"IsBanker"`
		Position     int    `json:"Position"`
		IsRoomMaster bool   `json:"IsRoomMaster"`
	}
)

var (
	// FreeMatchMgrInst 是用来处理玩家自由匹配的
	FreeMatchMgrInst *FreeMatchMgr
)

const (
	// AutoMatchMaxRoomID 是自由匹配的最大房间ID，非自由匹配用6位的房间号码
	AutoMatchMaxRoomID = 99999
)

// NewFreeMatchMgr 会创建一个FreeMatchMgr出来
func NewFreeMatchMgr() *FreeMatchMgr {
	mgr := new(FreeMatchMgr)
	mgr.matchMgrTool = NewMatchMgrTool()
	mgr.pvpRoomInfos = make(map[int]*StartBattleRoomInfo)
	mgr.sessions = make([]*session.Session, 0)
	source := rand.NewSource(time.Now().UnixNano())
	mgr.rand = rand.New(source)
	return mgr
}

// IsRoomEmpty 用来判定当前的房间是不是空的
func (mgr *FreeMatchMgr) IsRoomEmpty() bool {
	if mgr.playerCount == 0 {
		return true
	}
	return false
}

// AddMatchPlayer 会添加一个匹配游戏的玩家进来
func (mgr *FreeMatchMgr) AddMatchPlayer(s *session.Session) {
	if mgr.playerCount == 0 {
		mgr.currentRoomID = mgr.currentRoomID + 1
		if mgr.currentRoomID >= AutoMatchMaxRoomID {
			mgr.currentRoomID = 0
		}
		mgr.timerInst = nano.NewTimer(time.Second*3, mgr.onMatchTimeout)
	}

	mgr.sessions = append(mgr.sessions, s)
	mgr.playerCount = mgr.playerCount + 1

	if mgr.playerCount == 4 {
		mgr.onMatchFinished()
	}
}

func (mgr *FreeMatchMgr) onMatchTimeout() {
	mgr.onMatchFinished()
}

func (mgr *FreeMatchMgr) onMatchFinished() {

	if mgr.timerInst != nil {
		mgr.timerInst.Stop()
		mgr.timerInst = nil
	}

	startBattleInfo := &StartBattleInfo{}

	startBattleInfo.RandomSeed = mgr.rand.Intn(999999)

	startBattleInfo.BankerPosition = mgr.matchMgrTool.randAPosition(mgr.playerCount)
	startBattleInfo.AIOwnerPosition = mgr.matchMgrTool.randAPosition(mgr.playerCount)

	// SittingPosition 会在客户端计算，服务端没有计算
	startBattleInfo.SittingPosition = playerposition.Dong

	startBattleRoomInfo := &StartBattleRoomInfo{}

	startBattleRoomInfo.RoomID = mgr.currentRoomID
	startBattleRoomInfo.StartBattleInfo = startBattleInfo
	mgr.matchMgrTool.fillPlayerInfo(startBattleInfo, mgr.sessions)

	PVPMgrInst.RegisterPVPSessionInfo(mgr.currentRoomID, mgr.sessions)
	for _, s := range mgr.sessions {
		s.Push("OnQueryStartBattleInfo", &startBattleRoomInfo)
		Logger.Println(fmt.Sprintf("push msg -> onQueryStartBattleInfo, uid: %d", s.UID()))
	}

	mgr.pvpRoomInfos[mgr.currentRoomID] = startBattleRoomInfo
	mgr.clearMatchInfo()

	for _, v := range startBattleInfo.PlayerInfos {
		Logger.Println("pvp players, UID:", v.UID, ", NickName:", v.NickName)
	}
}

func (mgr *FreeMatchMgr) clearMatchInfo() {
	// currentRoomID int
	mgr.playerCount = 0
	mgr.sessions = make([]*session.Session, 0)
	if mgr.timerInst != nil {
		mgr.timerInst.Stop()
		mgr.timerInst = nil
	}
	// rand          *rand.Rand
}

// RemoveBattleRoomInfo 在PVP结束的时候，把这个房间的初始信息释放掉
func (mgr *FreeMatchMgr) RemoveBattleRoomInfo(roomID int) {
	info, ok := mgr.pvpRoomInfos[roomID]
	if !ok {
		return
	}
	Logger.Println(fmt.Sprintf("remove battle room info of room %d", info.RoomID))
	delete(mgr.pvpRoomInfos, roomID)
}
