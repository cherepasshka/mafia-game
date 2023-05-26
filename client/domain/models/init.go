package models

import (
	"soa.mafia-game/client/domain/models/user"
	"soa.mafia-game/client/domain/models/user_strategy"
	proto "soa.mafia-game/proto/mafia-game"
)

func MakeUser(login string, role proto.Roles) models.User {
	base := models.BaseUser{
		Status: models.Alive,
		Login:  login,
	}
	if role == proto.Roles_Civilian {
		return &user_strategy.Civilian{
			BaseUser: base,
		}
	} else if role == proto.Roles_Mafia {
		return &user_strategy.Mafia{
			BaseUser: base,
		}
	} else if role == proto.Roles_Commissioner {
		return &user_strategy.Commissioner{
			BaseUser: base,
		}
	}
	return &base
}
