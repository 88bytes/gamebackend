package main

import (
	"mj/utils"

	"github.com/lonnng/nano/component"
	"github.com/lonnng/nano/session"
)

type (
	// GetRoomIDMsg has not been used
	GetRoomIDMsg struct {
	}
	// Match is home service
	Match struct {
		component.Base
	}
)

// GetPVPRoomInfo tells client the info of battle room
func (comp *Match) GetPVPRoomInfo(s *session.Session, msg *GetRoomIDMsg) error {
	utils.FreeMatchMgrInst.AddMatchPlayer(s)
	return nil
}
