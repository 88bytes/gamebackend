package utils

import (
	"fmt"
	"math/rand"
	"time"

	"mj/control_type"
	"mj/player_position"

	"github.com/88bytes/nano/session"
)

type (
	// RoomMatchMgr 用来管理私人房匹配的数据
	RoomMatchMgr struct {
		matchMgrTool           *MatchMgrTool
		roomMatchInfosOnServer map[int]*RoomMatchInfoOnServer
		rand                   *rand.Rand
	}

	// RoomMatchInfoOnServer 里面装了私人房的匹配信息
	RoomMatchInfoOnServer struct {
		sessions        []*session.Session
		startBattleInfo *StartBattleInfo
	}

	// S2CDestroyRoomMsg 是服务器用来通知客户端，这个房间销毁了
	S2CDestroyRoomMsg struct {
		RoomID int
	}
)

var (
	// RoomMatchMgrInst 是用来处理玩家自由匹配的
	RoomMatchMgrInst *RoomMatchMgr
)

// NewRoomMatchMgr 会创建一个FreeMatchMgr出来
func NewRoomMatchMgr() *RoomMatchMgr {
	mgr := new(RoomMatchMgr)
	mgr.matchMgrTool = NewMatchMgrTool()
	mgr.roomMatchInfosOnServer = make(map[int]*RoomMatchInfoOnServer)
	source := rand.NewSource(time.Now().UnixNano())
	mgr.rand = rand.New(source)
	return mgr
}

// CreateRoom 创建一个私人房
func (mgr *RoomMatchMgr) CreateRoom(ses *session.Session, maxBattleCount int, zhuaNiaoCount int) {
	// 计算RoomID
	PVPMgrInst.CurrentUsedRoomID = PVPMgrInst.CurrentUsedRoomID + 1
	if PVPMgrInst.CurrentUsedRoomID >= MaxRoomID {
		PVPMgrInst.CurrentUsedRoomID = MinRoomID
	}
	roomID := PVPMgrInst.CurrentUsedRoomID

	// 申请房间信息的内存空间
	roomInfoOnServer := new(RoomMatchInfoOnServer)
	mgr.roomMatchInfosOnServer[roomID] = roomInfoOnServer

	txt := fmt.Sprintf("createRoom %d.", roomID)
	logger.Println(txt)

	// 产生SessionList，把RoomMaster放进去
	roomInfoOnServer.sessions = make([]*session.Session, 0)
	roomInfoOnServer.sessions = append(roomInfoOnServer.sessions, ses)

	// 这里是真正的房间信息
	startBattleInfo := new(StartBattleInfo)

	// RoomID
	startBattleInfo.RoomID = roomID

	// 局数和抓鸟数量
	startBattleInfo.MaxBattleCount = maxBattleCount
	startBattleInfo.ZhuaNiaoCount = zhuaNiaoCount

	// 产生单个Player的基础信息
	playerInfo := new(StartBattlePlayerInfo)
	UID := uint(ses.UID())
	playerInfo.UID = UID
	userInfo := UserInfoUtilInst.GetUserInfo(UID)
	playerInfo.NickName = userInfo.NickName
	playerInfo.ControlType = controltype.ByPlayer
	playerInfo.IsBanker = false
	playerInfo.IsRoomMaster = true

	// 产生Player信息的列表，把Player的信息插进去
	startBattleInfo.PlayerInfos = make([]*StartBattlePlayerInfo, 0)
	startBattleInfo.PlayerInfos = append(startBattleInfo.PlayerInfos, playerInfo)

	roomInfoOnServer.startBattleInfo = startBattleInfo

	mgr.BroadcastOnUpdateRoomInfoMsg(roomID)
}

// JoinRoom 玩家加入房间
func (mgr *RoomMatchMgr) JoinRoom(ses *session.Session, roomID int) {
	roomInfoOnServer, ok := mgr.roomMatchInfosOnServer[roomID]
	if !ok {
		txt := fmt.Sprintf("room %d not exist.", roomID)
		logger.Fatal(txt)
		return
	}

	txt := fmt.Sprintf("joinRoom %d.", roomID)
	logger.Println(txt)

	roomInfoOnServer.sessions = append(roomInfoOnServer.sessions, ses)

	startBattleInfo := roomInfoOnServer.startBattleInfo

	// 产生单个Player的基础信息
	playerInfo := new(StartBattlePlayerInfo)
	UID := uint(ses.UID())
	playerInfo.UID = UID
	userInfo := UserInfoUtilInst.GetUserInfo(UID)
	playerInfo.NickName = userInfo.NickName
	playerInfo.ControlType = controltype.ByPlayer
	playerInfo.IsBanker = false
	playerInfo.IsRoomMaster = false

	// 把Player的信息插进去
	startBattleInfo.PlayerInfos = append(startBattleInfo.PlayerInfos, playerInfo)

	mgr.BroadcastOnUpdateRoomInfoMsg(roomID)
}

// BroadcastOnUpdateRoomInfoMsg 会把房间的信息广播给房间中所有的玩家
func (mgr *RoomMatchMgr) BroadcastOnUpdateRoomInfoMsg(roomID int) {
	roomInfoOnServer, ok := mgr.roomMatchInfosOnServer[roomID]
	if !ok {
		return
	}

	startBattleInfo := roomInfoOnServer.startBattleInfo
	for index, playerInfo := range startBattleInfo.PlayerInfos {
		position := playerposition.Dong + index
		playerInfo.Position = position
		s := roomInfoOnServer.sessions[index]
		s.Push("OnUpdateRoomInfo", startBattleInfo)
		txt := fmt.Sprintf("OnUpdateRoomInfo, index: %d, nickName: %s, position: %d", index, playerInfo.NickName, position)
		Logger.Println(txt)
	}
}

// QuitWaitForRoomReadyState 是一个客户端在等待的阶段就退出了房间
// 如果房主退出了，就解散房间，如果是 非房主 用户，那就单独把这个用户清理出去
func (mgr *RoomMatchMgr) QuitWaitForRoomReadyState(UID uint) {
	roomID := mgr.findRoomID(UID)

	// 检查能不能找到对应的RoomID
	Logger.Println(fmt.Sprintf("user %d is in room %d.", UID, roomID))
	if roomID == -1 {
		Logger.Println(fmt.Sprintf("room %d is invalid.", roomID))
		return
	}

	// 检查能不能找到对应的RoomID
	roomInfoOnServer, ok := mgr.roomMatchInfosOnServer[roomID]
	if !ok {
		Logger.Println(fmt.Sprintf("room %d is invalid.", roomID))
		return
	}

	startBattleInfo := roomInfoOnServer.startBattleInfo

	isRoomMasterLeaveThisRoom := false

	for _, playerInfo := range startBattleInfo.PlayerInfos {
		if playerInfo.UID == UID {
			if playerInfo.IsRoomMaster {
				isRoomMasterLeaveThisRoom = true
			}
		}
	}

	// 房主离开了，房间就解散了
	// 普通成员离开，就刷新房间信息
	if isRoomMasterLeaveThisRoom {
		mgr.destroyTheRoom(roomID)
	} else {
		mgr.playerLeaveTheRoom(roomID, UID)
		mgr.BroadcastOnUpdateRoomInfoMsg(roomID)
	}
}

func (mgr *RoomMatchMgr) findRoomID(UID uint) int {
	roomMatchInfosOnServer := mgr.roomMatchInfosOnServer
	for roomID, roomMatchInfoOnServer := range roomMatchInfosOnServer {
		playerInfos := roomMatchInfoOnServer.startBattleInfo.PlayerInfos
		for _, playerInfo := range playerInfos {
			if playerInfo.UID == UID {
				return roomID
			}
		}
	}
	return -1
}

func (mgr *RoomMatchMgr) destroyTheRoom(roomID int) {
	roomInfoOnServer, ok := mgr.roomMatchInfosOnServer[roomID]
	if !ok {
		Logger.Println(fmt.Sprintf("room %d is invalid.", roomID))
		return
	}

	sessions := roomInfoOnServer.sessions
	for _, session := range sessions {
		session.Push("OnDestroyTheRoom", S2CDestroyRoomMsg{roomID})
	}

	txt := fmt.Sprintf("destroyTheRoom -> %d", roomID)
	Logger.Println(txt)

	delete(mgr.roomMatchInfosOnServer, roomID)
}

func (mgr *RoomMatchMgr) playerLeaveTheRoom(roomID int, UID uint) {
	roomInfoOnServer, ok := mgr.roomMatchInfosOnServer[roomID]
	if !ok {
		Logger.Println(fmt.Sprintf("room %d is invalid.", roomID))
		return
	}

	// 删掉 sessions 里面对应的 session
	newSessions := make([]*session.Session, 0)
	sessions := roomInfoOnServer.sessions
	for _, session := range sessions {
		if session.UID() != int64(UID) {
			newSessions = append(newSessions, session)
		}
	}
	sessions = newSessions

	// 从 playerInfos 里面删除对应的 playerInfo
	newInfos := make([]*StartBattlePlayerInfo, 0)
	startBattleInfo := roomInfoOnServer.startBattleInfo
	for _, playerInfo := range startBattleInfo.PlayerInfos {
		if playerInfo.UID != UID {
			newInfos = append(newInfos, playerInfo)
		}
	}
	startBattleInfo.PlayerInfos = newInfos

	txt := fmt.Sprintf("player %d, leave room %d", UID, roomID)
	Logger.Println(txt)
}

// GetStartBattleInfo 返回启动战斗的信息，收到这个信息之后，就可以开始接受PVP同步的信息了
// TODO
func (mgr *RoomMatchMgr) GetStartBattleInfo(roomID int) {
	roomInfoOnServer, ok := mgr.roomMatchInfosOnServer[roomID]
	if !ok {
		return
	}

	startBattleInfo := roomInfoOnServer.startBattleInfo
	for _, playerInfo := range startBattleInfo.PlayerInfos {
		txt := fmt.Sprintf("player, nickName: %s", playerInfo.NickName)
		Logger.Println(txt)
		// s.Push("OnJoinPlayer", &startBattleRoomInfo)
	}

	mgr.matchMgrTool.fillPlayerInfo(startBattleInfo, roomInfoOnServer.sessions)

	PVPMgrInst.RegisterPVPSessionInfo(roomID, roomInfoOnServer.sessions)
	for _, s := range roomInfoOnServer.sessions {
		s.Push("OnQueryStartBattleInfo", startBattleInfo)
		Logger.Println(fmt.Sprintf("push msg -> onQueryStartBattleInfo, uid: %d", s.UID()))
	}

	for _, v := range startBattleInfo.PlayerInfos {
		Logger.Println("pvp players, UID:", v.UID, ", NickName:", v.NickName)
	}
}
