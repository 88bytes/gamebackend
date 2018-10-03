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

	// CreateRoomMsg contains info to create room
	CreateRoomMsg struct {
		MaxBattleCount int
		ZhuaNiaoCount  int
	}

	// JoinRoomMsg contains info to join room
	JoinRoomMsg struct {
		RoomNumber int
	}

	// RoomBattleQuickStartMsg contains info to make the battle quick start
	RoomBattleQuickStartMsg struct {
		RoomNumber int
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
func (comp *RoomMatch) JoinRoom(s *session.Session, msg *JoinRoomMsg) error {
	utils.Logger.Println(fmt.Sprintf("joinRoom -> uid: %d, roomNumber: %d", s.UID(), msg.RoomNumber))
	return nil
}

// RoomBattleQuickStart will request server start the battle now, AI players will quick join the game.
func (comp *RoomMatch) RoomBattleQuickStart(s *session.Session, msg *RoomBattleQuickStartMsg) error {
	utils.Logger.Println(fmt.Sprintf("roomBattleQuickStart -> uid: %d, roomNumber: %d", s.UID(), msg.RoomNumber))
	return nil
}
