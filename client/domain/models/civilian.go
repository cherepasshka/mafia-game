package models

import (
	proto "soa.mafia-game/proto/mafia-game"
)

type Civilian struct {
	BaseUser
}

func (user *Civilian) GetRole() proto.Roles {
	return proto.Roles_Civilian
}

func (user *Civilian) MakeNightMove(proto.MafiaServiceClient) error {
	return nil
}
