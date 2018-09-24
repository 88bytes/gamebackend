package main

import (
	"fmt"
	"mj/utils"

	"github.com/88bytes/nano/component"
	"github.com/88bytes/nano/session"
)

type (
	// RoomMatch is for room match service
	RoomMatch struct {
		component.Base
	}

	// CreateRoomMsg has not been used
	CreateRoomMsg struct {
		MaxBattleCount int
		ZhuaNiaoCount  int
	}
)

// CreateRoom means create a battle room
// 私人房的房间号和自由匹配的房间号是分开计数的，还是共享呢？
// 暂时我先做成分开计数吧，我开心这么多，就这么做
func (comp *RoomMatch) CreateRoom(s *session.Session, msg *CreateRoomMsg) error {
	uid := s.UID()

	maxBattleCount := msg.MaxBattleCount
	zhuaNiaoCount := msg.ZhuaNiaoCount

	txt := fmt.Sprintf("createRoom -> uid: %d, maxBattleCount: %d, zhuaNiaoCount: %d", uid, maxBattleCount, zhuaNiaoCount)
	utils.Logger.Println(txt)

	utils.RoomMatchMgrInst.CreateRoom(s, maxBattleCount, zhuaNiaoCount)

	return nil
}

// JoinRoom means create a battle room
func (comp *RoomMatch) JoinRoom(s *session.Session, msg *GetRoomIDMsg) error {
	utils.Logger.Println(fmt.Sprintf("joinRoom -> uid: %d", s.UID()))
	return nil
}
