package models

import (
	"soa.mafia-game/client/domain/models/random_strategy"
	"soa.mafia-game/client/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

func MakeUser(login string, role proto.Roles) models.User {
	base := models.BaseUser{
		Status: models.Alive,
		Login:  login,
	}
	if role == proto.Roles_Civilian {
		return &random_strategy.Civilian{
			BaseUser: base,
		}
	} else if role == proto.Roles_Mafia {
		return &random_strategy.Mafia{
			BaseUser: base,
		}
	} else if role == proto.Roles_Commissioner {
		return &random_strategy.Commissioner{
			BaseUser: base,
		}
	}
	return &base
}
