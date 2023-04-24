package models

import (
	"context"

	proto "soa.mafia-game/proto/mafia-game"
)

type Civilian struct {
	BaseUser
}

func (user *Civilian) GetRole() proto.Roles {
	return proto.Roles_Civilian
}

func (user *Civilian) MakeNightMove(alive []string, client proto.MafiaServiceClient) error {
	if user.Status == Alive {
		client.MakeMove(context.Background(), &proto.MoveRequest{Login: user.Login})
	}
	return nil
}
