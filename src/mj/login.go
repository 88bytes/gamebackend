package main

import (
	"mj/utils"

	"github.com/88bytes/nano"
	"github.com/88bytes/nano/component"
	"github.com/88bytes/nano/session"
)

type (
	// Login is login service
	Login struct {
		component.Base
	}

	// LoginMsg with uid and NickName
	LoginMsg struct {
		UID      uint
		NickName string
	}

	// LoginResp means Login Response
	LoginResp struct {
		Success  bool
		NickName string
	}
)

// AfterInit 注册了OnSessionClosed函数
func (comp *Login) AfterInit() {
	nano.OnSessionClosed(func(s *session.Session) {
		UID := uint(s.UID())
		utils.Logger.Println("sessionClosed, uid: ", UID)
		// 如果用户没有登陆过，这个地方就要return，不要再继续操作了
		if UID != 0 {
			utils.PVPMgrInst.QuitPVP(UID)
			utils.RoomMatchMgrInst.QuitWaitForRoomReadyState(UID)
			utils.UserInfoUtilInst.RemoveUserInfo(UID)
		}
	})
}

// Login bind user id
func (comp *Login) Login(s *session.Session, msg *LoginMsg) error {
	utils.Logger.Println("userLogin, uid: ", msg.UID)
	utils.Logger.Println("userLogin, NickName: ", msg.NickName)
	s.Bind(int64(msg.UID))
	utils.UserInfoUtilInst.AddUserInfo(msg.UID, msg.NickName)
	return s.Response(&LoginResp{Success: true})
}
