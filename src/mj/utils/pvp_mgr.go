package utils

import (
	"fmt"
	"time"

	"github.com/lonnng/nano"
	"github.com/lonnng/nano/session"
)

const (
	// ServerTickInterval 表示 200ms 的tick间隔
	ServerTickInterval = 200
)

type (
	// PVPMgr 用来管理PVP流程，驱动PVP的运行
	PVPMgr struct {
		RoomBattleDatas map[int]*BattleData
	}

	// BattleData 里面存了一个对战房间会用到的全部信息
	BattleData struct {
		Sessions       []*session.Session
		PlayerReadyMap map[uint]bool
		Timer          *nano.Timer
		IsRunning      bool
		CurFrameCmd    FrameCmd
	}

	// LockStepCmd 是玩家发送过来的每一帧的操作指令
	LockStepCmd struct {
		MJType              int `json:"MJType"`
		MJNumber            int `json:"MJNumber"`
		Cmd                 int `json:"Cmd"`
		OwenrPosition       int `json:"OwnerPosition"`
		ShouPaiIndex        int `json:"ShouPaiIndex"`
		ChiPaiIndex1        int `json:"ChiPaiIndex1"`
		ChiPaiIndex2        int `json:"ChiPaiIndex2"`
		WinType             int `json:"WinType"`
		OtherPlayerPosition int `json:"OtherPlayerPosition"`
	}

	// FrameCmd 包含了每一个同步帧里面所有玩家的指令
	FrameCmd struct {
		ServerFrameIndex int            `json:"ServerFrameIndex"`
		LockStepCmds     []*LockStepCmd `json:"LockStepCmds"`
	}
)

var (
	// PVPMgrInst 存了PVP过程中，每个Room里面的Session
	PVPMgrInst *PVPMgr
)

// NewPVPMgr 会创建一个PVPMgr出来
func NewPVPMgr() *PVPMgr {
	mgr := new(PVPMgr)
	mgr.RoomBattleDatas = make(map[int]*BattleData)
	return mgr
}

// RegisterPVPSessionInfo 在开始PVP的时候把所有的Session信息记录下来，PVP帧同步的时候会给这些Session广播消息
func (mgr *PVPMgr) RegisterPVPSessionInfo(roomID int, sessions []*session.Session) {
	battleData := new(BattleData)
	battleData.Sessions = sessions
	battleData.PlayerReadyMap = make(map[uint]bool)
	battleData.IsRunning = false
	for _, s := range sessions {
		uid := uint(s.UID())
		userInfo := UserInfoUtilInst.GetUserInfo(uid)
		userInfo.RoomID = roomID
	}
	mgr.RoomBattleDatas[roomID] = battleData
}

// StartPVP 申请启动PVP，当一个Room里面的所有玩家都发送了这个请求之后，就会启动PVP的同步
func (mgr *PVPMgr) StartPVP(session *session.Session) {
	uid := uint(session.UID())
	userInfo := UserInfoUtilInst.GetUserInfo(uid)
	roomID := userInfo.RoomID

	battleData, ok := mgr.RoomBattleDatas[roomID]
	if !ok {
		Logger.Println(fmt.Sprintf("room %d battleData not found", roomID))
		return
	}

	battleData.PlayerReadyMap[uid] = true

	Logger.Println(fmt.Sprintf("player %d of room %d is ready", uid, roomID))
	if mgr.isAllPlayerReady(roomID) {
		mgr.startPVPSync(roomID)
	}
}

func (mgr *PVPMgr) isAllPlayerReady(roomID int) bool {
	battleData := mgr.RoomBattleDatas[roomID]
	for _, s := range battleData.Sessions {
		uid := uint(s.UID())
		_, ok := battleData.PlayerReadyMap[uid]
		if !ok {
			return false
		}
	}
	return true
}

// 这个房间正式开始PVP的帧同步
func (mgr *PVPMgr) startPVPSync(roomID int) {
	Logger.Println(fmt.Sprintf("players of room %d is all ready", roomID))

	battleData := mgr.RoomBattleDatas[roomID]
	battleData.IsRunning = true
	battleData.CurFrameCmd = FrameCmd{}
	battleData.CurFrameCmd.ServerFrameIndex = 0

	// 先把启动比赛的指令插进去
	battleData.CurFrameCmd.LockStepCmds = make([]*LockStepCmd, 0)
	cmd := new(LockStepCmd)
	// PlayerPosition.CompetitionMgr
	cmd.OwenrPosition = 1
	// LockStepCmdType.StartBattle
	cmd.Cmd = 1
	battleData.CurFrameCmd.LockStepCmds = append(battleData.CurFrameCmd.LockStepCmds, cmd)

	battleData.Timer = nano.NewTimer(time.Millisecond*ServerTickInterval, func() {
		pvpTick(roomID)
	})
}

func pvpTick(roomID int) {
	data := PVPMgrInst.RoomBattleDatas[roomID]

	userSessions := data.Sessions
	for _, session := range userSessions {
		session.Push("OnReceiveFrameCmd", data.CurFrameCmd)
	}

	data.CurFrameCmd.ServerFrameIndex = data.CurFrameCmd.ServerFrameIndex + 1
	data.CurFrameCmd.LockStepCmds = make([]*LockStepCmd, 0)
}

// InsertLockStepCmd 会把指令存下来，下一帧同步的时候会发送给客户端
func (mgr *PVPMgr) InsertLockStepCmd(roomID int, cmd *LockStepCmd) {
	Logger.Println(fmt.Sprintf("roomID %d receive LockStepCmd.", roomID))
	battleData, ok := PVPMgrInst.RoomBattleDatas[roomID]
	if !ok {
		return
	}
	battleData.CurFrameCmd.LockStepCmds = append(battleData.CurFrameCmd.LockStepCmds, cmd)
}

// QuitPVP 客户端退出PVP对战
func (mgr *PVPMgr) QuitPVP(userID uint) {
	userInfo := UserInfoUtilInst.GetUserInfo(uint(userID))
	roomID := userInfo.RoomID
	Logger.Println(fmt.Sprintf("user %d of room %d try to quitPVP", userID, roomID))
	battleData, ok := mgr.RoomBattleDatas[roomID]
	if !ok {
		Logger.Println(fmt.Sprintf("user %d of room %d is not in pvp game", userID, roomID))
		return
	}
	battleData.PlayerReadyMap[userID] = false
	if mgr.isAllPlayerQuitPVP(roomID) {
		battleData.IsRunning = false
		battleData.Timer.Stop()
		delete(mgr.RoomBattleDatas, roomID)
		FreeMatchMgrInst.RemoveBattleRoomInfo(roomID)
		Logger.Println(fmt.Sprintf("all players in room %d quit the pvp game", roomID))
	}
}

func (mgr *PVPMgr) isAllPlayerQuitPVP(roomID int) bool {
	battleData, ok := mgr.RoomBattleDatas[roomID]

	if !ok {
		return true
	}

	allPlayerNotReady := true
	for _, s := range battleData.Sessions {
		uid := uint(s.UID())
		bReady, ok := battleData.PlayerReadyMap[uid]
		if !ok {
			return true
		}
		// 如果所有的人，都不在Ready状态了，那么就认为全部都退出了
		allPlayerNotReady = allPlayerNotReady && !bReady
	}

	return allPlayerNotReady
}
