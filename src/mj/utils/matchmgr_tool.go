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
	// MatchMgrTool 这个工具里面有 FreeMatchMgr RoomMatchMgr 共享的工具代码
	MatchMgrTool struct {
		rand *rand.Rand
	}
)

// NewMatchMgrTool 会创建一个MatchMgrTool出来
func NewMatchMgrTool() *MatchMgrTool {
	mgr := new(MatchMgrTool)
	source := rand.NewSource(time.Now().UnixNano())
	mgr.rand = rand.New(source)
	return mgr
}

func (mgr *MatchMgrTool) randAPosition(playerCount int) int {
	randInt := int(mgr.rand.Float32() * float32(playerCount))
	randPosition := playerposition.Dong + randInt
	if randPosition >= playerposition.Dong+4 {
		randPosition = playerposition.Dong + 3
	}
	return randPosition
}

func (mgr *MatchMgrTool) fillPlayerInfoOfFreeMatch(battleInfo *StartBattleInfo, sessions []*session.Session) {
	battleInfo.PlayerInfos = make([]*StartBattlePlayerInfo, 0)
	var computerIndex uint
	computerIndex = 1
	for index := 0; index < 4; index++ {
		playerInfo := new(StartBattlePlayerInfo)
		length := len(sessions)
		if index < length {
			UID := uint(sessions[index].UID())
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

func (mgr *MatchMgrTool) fillPlayerInfoOfRoomMatch(battleInfo *StartBattleInfo, roomMatchInfoOnServer *RoomMatchInfoOnServer) {
	sessions := roomMatchInfoOnServer.sessions
	preSavedPlayerInfos := roomMatchInfoOnServer.startBattleInfo.PlayerInfos

	if len(sessions) != len(preSavedPlayerInfos) {
		txt := fmt.Sprintf("fatal error, playerInfo len(%d) != sessions len(%d)", len(preSavedPlayerInfos), len(sessions))
		Logger.Fatal(txt)
		return
	}

	battleInfo.PlayerInfos = make([]*StartBattlePlayerInfo, 0)
	var computerIndex uint
	computerIndex = 1
	for index := 0; index < 4; index++ {
		length := len(sessions)

		// 如果已经有余存的数据，就直接取出来
		var playerInfo *StartBattlePlayerInfo
		if index < length {
			playerInfo = preSavedPlayerInfos[index]
		} else {
			playerInfo = new(StartBattlePlayerInfo)
		}

		if index < length {
			// 根据预存的数据，设置房主为 AiOwner
			if playerInfo.IsRoomMaster {
				battleInfo.AIOwnerPosition = playerInfo.Position
			}
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
