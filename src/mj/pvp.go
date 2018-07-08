package main

import (
	"mj/utils"

	"github.com/88bytes/nano/component"
	"github.com/88bytes/nano/session"
)

type (
	// PVP is the pvp service
	PVP struct {
		component.Base
	}

	// EmptyMsg is an empty msg
	EmptyMsg struct {
	}
)

// StartPVP will start the pvp competition
func (comp *PVP) StartPVP(s *session.Session, msg *EmptyMsg) error {
	utils.PVPMgrInst.StartPVP(s)
	return nil
}

// InsertLockStepCmd 在帧同步的流程中插入玩家的操作指令
func (comp *PVP) InsertLockStepCmd(s *session.Session, msg *utils.LockStepCmd) error {
	userInfo := utils.UserInfoUtilInst.GetUserInfo(uint(s.UID()))
	utils.PVPMgrInst.InsertLockStepCmd(userInfo.RoomID, msg)
	return nil
}

// QuitPVP 在客户端主动退出PVP的游戏的时候会被调用
func (comp *PVP) QuitPVP(s *session.Session, msg *EmptyMsg) error {
	utils.PVPMgrInst.QuitPVP(uint(s.UID()))
	return nil
}
