package random_strategy

import (
	"context"

	"soa.mafia-game/client/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

type Civilian struct {
	models.BaseUser
}

func (user *Civilian) GetRole() proto.Roles {
	return proto.Roles_Civilian
}

func (user *Civilian) MakeNightMove(alive []string, client proto.MafiaServiceClient) error {
	if user.Status == models.Alive {
		client.MakeMove(context.Background(), &proto.MoveRequest{Login: user.Login})
	}
	return nil
}
