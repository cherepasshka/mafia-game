package models

import (
	"fmt"

	// "soa.mafia-game/client/domain/chat"
	"soa.mafia-game/client/domain/models/user"
	"soa.mafia-game/client/domain/models/user_strategy"
	proto "soa.mafia-game/proto/mafia-game"
)

func MakeUser(login string, role proto.Roles, session string, partition int32) models.User {
	base := models.BaseUser{
		Status: models.Alive,
		Login:  login,
		Session: session,
		Partition: partition,
	}
	//chatService, _ := chat.New("kafka1:9092")
	fmt.Printf("After chat.New\n")
	if role == proto.Roles_Civilian {
		return &user_strategy.Civilian{
			BaseUser: base,
		}
	} else if role == proto.Roles_Mafia {
		return &user_strategy.Mafia{
			BaseUser: base,
			//ChatService: chatService,
		}
	} else if role == proto.Roles_Commissioner {
		return &user_strategy.Commissioner{
			BaseUser: base,
		}
	}
	return &base
}
