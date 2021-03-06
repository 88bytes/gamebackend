package main

import (
	"mj/utils"

	"github.com/88bytes/nano/component"
	"github.com/88bytes/nano/session"
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

// QueryStartBattleInfo tells client the info of battle room
func (comp *Match) QueryStartBattleInfo(s *session.Session, msg *GetRoomIDMsg) error {
	utils.FreeMatchMgrInst.AddMatchPlayer(s)
	return nil
}
