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
		OpenID     string
		NickName   string
		HeadImgURL string
		UID        uint
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
	utils.Logger.Println("userLogin, OpenID: ", msg.OpenID)
	utils.Logger.Println("userLogin, NickName: ", msg.NickName)
	utils.Logger.Println("userLogin, HeadImgURL: ", msg.HeadImgURL)
	utils.Logger.Println("userLogin, UID: ", msg.UID)
	s.Bind(int64(msg.UID))
	utils.UserInfoUtilInst.AddUserInfo(msg.OpenID, msg.NickName, msg.HeadImgURL, msg.UID)
	return s.Response(&LoginResp{Success: true})
}
