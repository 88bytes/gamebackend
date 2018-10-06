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
		matchMgrTool     *MatchMgrTool
		currentRoomID    int
		startBattleInfos map[int]*RoomMatchInfoOnServer
		rand             *rand.Rand
	}

	// RoomMatchInfoOnServer 里面装了私人房的匹配信息
	RoomMatchInfoOnServer struct {
		sessions        []*session.Session
		startBattleInfo *StartBattleInfo
	}
)

var (
	// RoomMatchMgrInst 是用来处理玩家自由匹配的
	RoomMatchMgrInst *RoomMatchMgr
)

const (
	// RoomMatchDefaultRoomID 是私人房的初始房号
	RoomMatchDefaultRoomID = 1024
	// RoomMatchMaxRoomID 是自由匹配的最大房间ID，非自由匹配用6位的房间号码
	RoomMatchMaxRoomID = 9024
)

// NewRoomMatchMgr 会创建一个FreeMatchMgr出来
func NewRoomMatchMgr() *RoomMatchMgr {
	mgr := new(RoomMatchMgr)
	mgr.matchMgrTool = NewMatchMgrTool()
	mgr.currentRoomID = RoomMatchDefaultRoomID
	mgr.startBattleInfos = make(map[int]*RoomMatchInfoOnServer)
	source := rand.NewSource(time.Now().UnixNano())
	mgr.rand = rand.New(source)
	return mgr
}

// CreateRoom 创建一个私人房
func (mgr *RoomMatchMgr) CreateRoom(ses *session.Session, maxBattleCount int, zhuaNiaoCount int) {
	// 计算RoomID
	roomID := mgr.currentRoomID
	mgr.currentRoomID = mgr.currentRoomID + 1
	if mgr.currentRoomID >= RoomMatchMaxRoomID {
		mgr.currentRoomID = RoomMatchDefaultRoomID
	}

	// 申请房间信息的内存空间
	roomInfoOnServer := new(RoomMatchInfoOnServer)
	mgr.startBattleInfos[roomID] = roomInfoOnServer

	txt := fmt.Sprintf("createRoom %d.", roomID)
	logger.Println(txt)

	// 产生SessionList，把RoomMaster放进去
	roomInfoOnServer.sessions = make([]*session.Session, 0)
	roomInfoOnServer.sessions = append(roomInfoOnServer.sessions, ses)

	// 这里是真正的房间信息
	startBattleInfo := new(StartBattleInfo)

	// RoomID
	startBattleInfo.RoomtID = roomID

	// 局数和抓鸟数量
	startBattleInfo.MaxBattleCount = maxBattleCount
	startBattleInfo.ZhuaNiaoCount = zhuaNiaoCount

	// 产生单个Player的基础信息
	playerInfo := StartBattlePlayerInfo{}
	UID := uint(ses.UID())
	playerInfo.UID = UID
	userInfo := UserInfoUtilInst.GetUserInfo(UID)
	playerInfo.NickName = userInfo.NickName
	playerInfo.ControlType = controltype.ByPlayer
	playerInfo.IsRoomMaster = true

	// 产生Player信息的列表，把Player的信息插进去
	startBattleInfo.PlayerInfos = make([]StartBattlePlayerInfo, 0)
	startBattleInfo.PlayerInfos = append(startBattleInfo.PlayerInfos, playerInfo)

	roomInfoOnServer.startBattleInfo = startBattleInfo

	mgr.BroadcastOnUpdateRoomInfoMsg(roomID)
}

// JoinRoom 玩家加入房间
func (mgr *RoomMatchMgr) JoinRoom(s *session.Session, roomID int) {
	roomInfoOnServer, ok := mgr.startBattleInfos[roomID]
	if !ok {
		txt := fmt.Sprintf("room %d not exist.", roomID)
		logger.Println(txt)
		return
	}

	txt := fmt.Sprintf("joinRoom %d.", roomID)
	logger.Println(txt)

	roomInfoOnServer.sessions = append(roomInfoOnServer.sessions, s)

	startBattleInfo := roomInfoOnServer.startBattleInfo

	// 产生单个Player的基础信息
	playerInfo := StartBattlePlayerInfo{}
	UID := uint(s.UID())
	playerInfo.UID = UID
	userInfo := UserInfoUtilInst.GetUserInfo(UID)
	playerInfo.NickName = userInfo.NickName
	playerInfo.ControlType = controltype.ByPlayer
	playerInfo.IsRoomMaster = true

	// 把Player的信息插进去
	startBattleInfo.PlayerInfos = append(startBattleInfo.PlayerInfos, playerInfo)

	mgr.BroadcastOnUpdateRoomInfoMsg(roomID)
}

// BroadcastOnUpdateRoomInfoMsg 会把房间的信息广播给房间中所有的玩家
func (mgr *RoomMatchMgr) BroadcastOnUpdateRoomInfoMsg(roomID int) {
	roomInfoOnServer, ok := mgr.startBattleInfos[roomID]
	if !ok {
		return
	}

	startBattleInfo := roomInfoOnServer.startBattleInfo
	for index, playerInfo := range startBattleInfo.PlayerInfos {
		txt := fmt.Sprintf("OnUpdateRoomInfo, nickName: %s", playerInfo.NickName)
		Logger.Println(txt)
		s := roomInfoOnServer.sessions[index]
		s.Push("OnUpdateRoomInfo", startBattleInfo)
	}
}

// GetStartBattleRoomInfo 返回启动战斗的信息，收到这个信息之后，就可以开始接受PVP同步的信息了
func (mgr *RoomMatchMgr) GetStartBattleRoomInfo(roomID int) {
	roomInfoOnServer, ok := mgr.startBattleInfos[roomID]
	if !ok {
		return
	}

	startBattleInfo := roomInfoOnServer.startBattleInfo
	for _, playerInfo := range startBattleInfo.PlayerInfos {
		txt := fmt.Sprintf("player, nickName: %s", playerInfo.NickName)
		Logger.Println(txt)
		// s.Push("OnJoinPlayer", &startBattleRoomInfo)
	}

	startBattleInfo.RandomSeed = mgr.rand.Intn(999999)

	startBattleInfo.BankerPosition = mgr.matchMgrTool.randAPosition(4)
	startBattleInfo.AIOwnerPosition = mgr.matchMgrTool.randAPosition(4)
	// SittingPosition 会在客户端计算，服务端没有计算
	startBattleInfo.SittingPosition = playerposition.Dong

	startBattleRoomInfo := &StartBattleRoomInfo{}

	startBattleRoomInfo.RoomID = mgr.currentRoomID
	startBattleRoomInfo.StartBattleInfo = startBattleInfo
	mgr.matchMgrTool.fillPlayerInfo(startBattleInfo, roomInfoOnServer.sessions)

	PVPMgrInst.RegisterPVPSessionInfo(mgr.currentRoomID, roomInfoOnServer.sessions)
	for _, s := range roomInfoOnServer.sessions {
		s.Push("OnQueryStartBattleInfo", &startBattleRoomInfo)
		Logger.Println(fmt.Sprintf("push msg -> onQueryStartBattleInfo, uid: %d", s.UID()))
	}

	for _, v := range startBattleInfo.PlayerInfos {
		Logger.Println("pvp players, UID:", v.UID, ", NickName:", v.NickName)
	}
}
