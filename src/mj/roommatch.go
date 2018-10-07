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
		RoomID int
	}

	// RoomBattleQuickStartMsg contains info to make the battle quick start
	RoomBattleQuickStartMsg struct {
		RoomID int
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
	utils.Logger.Println(fmt.Sprintf("joinRoom -> uid: %d, roomID: %d", s.UID(), msg.RoomID))
	utils.RoomMatchMgrInst.JoinRoom(s, msg.RoomID)
	return nil
}

// QuitWaitForRoomReadyState 是一个客户端在等待的阶段就退出了房间
// 如果房主退出了，就解散房间，如果是 非房主 用户，那就单独把这个用户清理出去
func (comp *RoomMatch) QuitWaitForRoomReadyState(s *session.Session, msg *EmptyMsg) error {
	UID := s.UID()
	uUID := uint(UID)
	utils.Logger.Println(fmt.Sprintf("QuitWaitForRoomReadyState -> uid: %d", UID))
	utils.RoomMatchMgrInst.QuitWaitForRoomReadyState(uUID)
	return nil
}

// RoomBattleQuickStart will request server start the battle now, AI players will quick join the game.
func (comp *RoomMatch) RoomBattleQuickStart(s *session.Session, msg *RoomBattleQuickStartMsg) error {
	utils.Logger.Println(fmt.Sprintf("roomBattleQuickStart -> uid: %d, roomID: %d", s.UID(), msg.RoomID))
	return nil
}
