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
		matchMgrTool *MatchMgrTool
		playerCount  int
		sessions     []*session.Session
		timerInst    *nano.Timer
		rand         *rand.Rand
		pvpRoomInfos map[int]*StartBattleInfo
	}

	// StartBattleInfo contains battle info
	StartBattleInfo struct {
		RoomID          int
		RandomSeed      int
		BankerPosition  int
		SittingPosition int
		AIOwnerPosition int
		MaxBattleCount  int
		ZhuaNiaoCount   int
		PlayerInfos     []*StartBattlePlayerInfo
	}

	// StartBattlePlayerInfo contains a battlePlayerInfo
	StartBattlePlayerInfo struct {
		UID          uint
		NickName     string
		ControlType  int
		IsBanker     bool
		Position     int
		IsRoomMaster bool
	}
)

var (
	// FreeMatchMgrInst 是用来处理玩家自由匹配的
	FreeMatchMgrInst *FreeMatchMgr
)

// NewFreeMatchMgr 会创建一个FreeMatchMgr出来
func NewFreeMatchMgr() *FreeMatchMgr {
	mgr := new(FreeMatchMgr)
	mgr.matchMgrTool = NewMatchMgrTool()
	mgr.pvpRoomInfos = make(map[int]*StartBattleInfo)
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
		PVPMgrInst.CurrentUsedRoomID = PVPMgrInst.CurrentUsedRoomID + 1
		if PVPMgrInst.CurrentUsedRoomID >= MaxRoomID {
			PVPMgrInst.CurrentUsedRoomID = MinRoomID
		}
		PVPMgrInst.CurrentFreeMatchRoomID = PVPMgrInst.CurrentUsedRoomID
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

	startBattleInfo := new(StartBattleInfo)

	roomID := PVPMgrInst.CurrentFreeMatchRoomID

	startBattleInfo.RoomID = roomID

	startBattleInfo.RandomSeed = mgr.rand.Intn(999999)

	startBattleInfo.BankerPosition = mgr.matchMgrTool.randAPosition(mgr.playerCount)

	// SittingPosition 会在客户端计算，服务端没有计算
	startBattleInfo.SittingPosition = playerposition.Dong

	startBattleInfo.AIOwnerPosition = mgr.matchMgrTool.randAPosition(mgr.playerCount)

	mgr.matchMgrTool.fillPlayerInfoOfFreeMatch(startBattleInfo, mgr.sessions)

	PVPMgrInst.RegisterPVPSessionInfo(roomID, mgr.sessions)
	for _, s := range mgr.sessions {
		s.Push("OnQueryStartBattleInfo", startBattleInfo)
		Logger.Println(fmt.Sprintf("push msg -> onQueryStartBattleInfo, uid: %d", s.UID()))
	}

	mgr.pvpRoomInfos[roomID] = startBattleInfo
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
