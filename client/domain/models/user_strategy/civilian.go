package user_strategy

import (
	"context"
	"fmt"

	"soa.mafia-game/client/domain/chat"
	"soa.mafia-game/client/domain/models/user"
	"soa.mafia-game/client/internal/utils/console"
	proto "soa.mafia-game/proto/mafia-game"
)

type Civilian struct {
	models.BaseUser
	ChatService *chat.ChatService
}

func (user *Civilian) GetRole() proto.Roles {
	return proto.Roles_Civilian
}

func (user *Civilian) MakeNightMove(ctx context.Context, alive []string, client proto.MafiaServiceClient) error {
	if user.Status == models.Alive {
		client.MakeMove(ctx, &proto.MoveRequest{Login: user.Login})
	}
	return nil
}

func (user *Civilian) VoteForMafia(ctx context.Context, alive_players []string, client proto.MafiaServiceClient) error {
	guess := user.Login
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this day vote")
		guess = "None"
	} else {
		user.ChatService.Start(user.Login, user.Session, user.Partition)
		for guess == user.Login {
			guess, _ = console.AskPrompt("Select your mafia guess", alive_players)
		}
		fmt.Printf("You voted for %s\n", guess)
	}
	rsp, err := client.VoteForMafia(ctx, &proto.VoteForMafiaRequest{Login: user.Login, MafiaGuess: guess})
	if err != nil {
		return err
	}
	if rsp.KilledUser == user.Login {
		fmt.Println("Most voted for you")
	} else {
		fmt.Printf("Most voted for %s, this user had role: %s\n", rsp.KilledUser, rsp.KilledUserRole)
	}
	return nil
}
