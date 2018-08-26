package utils

import (
	"fmt"
	"math/rand"
	"time"

	"mj/control_type"
	"mj/match_battle_type"
	"mj/player_position"

	"github.com/88bytes/nano"
	"github.com/88bytes/nano/session"
)

type (
	// FreeMatchMgr 用来管理自由匹配的数据，还有定时器
	// pvpRoomInfos 这个 map 要注意没有用的要回收，要不然内存就会越来越大
	FreeMatchMgr struct {
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
		MatchBattleType int              `json:"MatchBattleType"`
		MaxBattleCount  int              `json:"MaxBattleCount"`
		ZhuaNiaoCount   int              `json:"ZhuaNiaoCount"`
		StartBattleInfo *StartBattleInfo `json:"StartBattleInfo"`
	}

	// StartBattleInfo contains battle info
	StartBattleInfo struct {
		RoomtID         int                     `json:"RoomID"`
		RandomSeed      int                     `json:"RandomSeed"`
		BankerPosition  int                     `json:"BankerPosition"`
		SittingPosition int                     `json:"SittingPosition"`
		AIOwnerPosition int                     `json:"AIOwnerPosition"`
		PlayerInfos     []StartBattlePlayerInfo `json:"PlayerInfos"`
	}

	// StartBattlePlayerInfo contains a battlePlayerInfo
	StartBattlePlayerInfo struct {
		UID         uint   `json:"uid"`
		NickName    string `json:"NickName"`
		ControlType int    `json:"ControlType"`
		IsBanker    bool   `json:"IsBanker"`
		Position    int    `json:"Position"`
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

	startBattleInfo.BankerPosition = mgr.randAPosition()
	startBattleInfo.AIOwnerPosition = mgr.randAPosition()
	// SittingPosition 会在客户端计算，服务端没有计算
	startBattleInfo.SittingPosition = playerposition.Dong

	startBattleRoomInfo := &StartBattleRoomInfo{}

	startBattleRoomInfo.RoomID = mgr.currentRoomID
	startBattleRoomInfo.MatchBattleType = matchbattletype.QuickBattle
	startBattleRoomInfo.MaxBattleCount = 4
	startBattleRoomInfo.ZhuaNiaoCount = 4
	startBattleRoomInfo.StartBattleInfo = startBattleInfo
	mgr.fillPlayerInfo(startBattleInfo)

	PVPMgrInst.RegisterPVPSessionInfo(mgr.currentRoomID, mgr.sessions)
	for _, s := range mgr.sessions {
		s.Push("OnGetPVPRoomInfo", &startBattleRoomInfo)
		Logger.Println(fmt.Sprintf("push msg -> onGetPVPRoomInfo, uid: %d", s.UID()))
	}

	mgr.pvpRoomInfos[mgr.currentRoomID] = startBattleRoomInfo
	mgr.clearMatchInfo()

	for _, v := range startBattleInfo.PlayerInfos {
		Logger.Println("pvp players, UID:", v.UID, ", NickName:", v.NickName)
	}
}

func (mgr *FreeMatchMgr) randAPosition() int {
	playerCount := mgr.playerCount
	randInt := int(mgr.rand.Float32() * float32(playerCount))
	randPosition := playerposition.Dong + randInt
	if randPosition >= playerposition.Dong+4 {
		randPosition = playerposition.Dong + 3
	}
	return randPosition
}

func (mgr *FreeMatchMgr) fillPlayerInfo(battleInfo *StartBattleInfo) {
	battleInfo.PlayerInfos = make([]StartBattlePlayerInfo, 0)
	var computerIndex uint
	computerIndex = 1
	for index := 0; index < 4; index++ {
		playerInfo := StartBattlePlayerInfo{}
		length := len(mgr.sessions)
		if index < length {
			UID := uint(mgr.sessions[index].UID())
			userInfo := UserInfoUtilInst.GetUserInfo(UID)
			playerInfo.UID = userInfo.UID
			playerInfo.NickName = userInfo.NickName
		} else {
			playerInfo.UID = computerIndex
			playerInfo.NickName = fmt.Sprintf("COM%d", computerIndex)
			computerIndex = computerIndex + 1
		}
		position := playerposition.Dong + index
		if playerInfo.UID < 100 {
			playerInfo.ControlType = controltype.ByAi
		} else {
			playerInfo.ControlType = controltype.ByPlayer
		}
		if position == battleInfo.BankerPosition {
			playerInfo.IsBanker = true
		} else {
			playerInfo.IsBanker = false
		}
		playerInfo.Position = position
		battleInfo.PlayerInfos = append(battleInfo.PlayerInfos, playerInfo)
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
