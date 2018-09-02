package main

import (
	"mj/utils"

	"github.com/88bytes/nano/component"
	"github.com/88bytes/nano/session"
)

type (
	// RoomMatch is for room match service
	RoomMatch struct {
		component.Base
	}
)

// QueryStartBattleInfo tells client the info of battle room
func (comp *RoomMatch) QueryStartBattleInfo(s *session.Session, msg *GetRoomIDMsg) error {
	utils.FreeMatchMgrInst.AddMatchPlayer(s)
	return nil
}
