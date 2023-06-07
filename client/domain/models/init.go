package models

import (
	"soa.mafia-game/client/domain/chat"
	"soa.mafia-game/client/domain/models/user"
	"soa.mafia-game/client/domain/models/user_strategy"
	proto "soa.mafia-game/proto/mafia-game"
)

func MakeUser(login string, role proto.Roles, session string, partition int32) models.User {
	base := models.CommunicatorUser{
		BaseUser: models.BaseUser{
			Status:    models.Alive,
			Login:     login,
			Session:   session,
			Partition: partition,
		},
		ExitedChat: true,
	}
	chatService, _ := chat.New("localhost:9092")
	if role == proto.Roles_Civilian {
		return &user_strategy.Civilian{
			CommunicatorUser: base,
			ChatService:      chatService,
		}
	} else if role == proto.Roles_Mafia {
		return &user_strategy.Mafia{
			CommunicatorUser: base,
			ChatService:      chatService,
		}
	} else if role == proto.Roles_Commissioner {
		return &user_strategy.Commissioner{
			CommunicatorUser: base,
			ChatService:      chatService,
		}
	}
	return &base
}
