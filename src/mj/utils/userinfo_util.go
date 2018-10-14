package utils

import (
	"log"
	"os"
)

type (
	// UserInfoUtil 用来管理登陆进来的用户信息
	UserInfoUtil struct {
		loginCount uint
		userInfos  map[uint]*UserInfo
	}

	// UserInfo 存了每个用户的基础信息
	// 说实话，Go在内存的控制上比JS舒服多了，用Pomelo就是觉得虚，分布式的功能我觉得我也有能力开发出来
	UserInfo struct {
		OpenID     string
		NickName   string
		HeadImgURL string
		UID        uint
		RoomID     int
	}
)

var (
	// UserInfoUtilInst 存了所有在线用户的信息
	UserInfoUtilInst *UserInfoUtil
	logger           ILogger = log.New(os.Stderr, "", log.LstdFlags|log.Llongfile)
)

// NewUserInfoUtil 用来创建一个UserInfoUtil出来
func NewUserInfoUtil() *UserInfoUtil {
	var util *UserInfoUtil
	util = new(UserInfoUtil)
	util.userInfos = make(map[uint]*UserInfo)
	util.loginCount = 0
	return util
}

// AddUserInfo 在用户登陆的时候，把他的用户信息记录下来
func (util *UserInfoUtil) AddUserInfo(openID string, nickName string, headImgURL string, UID uint) {
	_, ok := util.userInfos[UID]
	if ok {
		return
	}
	var info UserInfo
	info = UserInfo{openID, nickName, headImgURL, UID, 0}
	util.userInfos[UID] = &info
	util.loginCount = util.loginCount + 1
	logger.Println("AddUserInfo when he login, uid:", UID, ", nickName:", nickName, ", loginCount:", util.loginCount)
}

// RemoveUserInfo 在网络链接断开的时候会清理掉用户信息
func (util *UserInfoUtil) RemoveUserInfo(UID uint) {
	_, ok := util.userInfos[UID]
	if !ok {
		return
	}
	info := util.userInfos[UID]
	delete(util.userInfos, UID)
	util.loginCount = util.loginCount - 1
	logger.Println("RemoveUserInfo, uid:", UID, ", nickName:", info.NickName, ", loginCount:", util.loginCount)
}

// GetUserInfo 用来查询在线的用户信息，查不到就是nil
func (util *UserInfoUtil) GetUserInfo(UID uint) *UserInfo {
	value, ok := util.userInfos[UID]
	if ok {
		return value
	}
	return nil
}
