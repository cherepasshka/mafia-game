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
	models.CommunicatorUser
	ChatService *chat.ChatService
}

func (user *Civilian) GetRole() proto.Roles {
	return proto.Roles_Civilian
}

func (user *Civilian) MakeNightMove(ctx context.Context, alive []string, client proto.MafiaServiceClient) (isValid bool, err error) {
	if user.Status == models.Alive {
		rsp, err := client.MakeMove(ctx, &proto.MoveRequest{Login: user.Login})
		return rsp.SessionStatus.AllConnected, err
	}
	return true, nil
}

func (user *Civilian) VoteForMafia(ctx context.Context, alive_players []string, client proto.MafiaServiceClient) (isValid bool, err error) {
	var guess string
	user.ExitedChat = false
	user.ChatService.Start(user.Login, user.Session, user.Partition, user.Status == models.Dead)
	user.ExitedChat = true
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this day vote, please wait until alive finish discussion")
		guess = "None"
	} else {
		guess, _ = console.AskPrompt("Select your mafia guess", user.ExcludeFromAliveList(alive_players))
		fmt.Printf("You voted for %s\n", guess)
	}
	rsp, err := client.VoteForMafia(ctx, &proto.VoteForMafiaRequest{Login: user.Login, MafiaGuess: guess})
	if err != nil {
		return false, err
	}
	if rsp.KilledUser == user.Login {
		fmt.Println("Most voted for you")
	} else {
		fmt.Printf("Most voted for %s, this user had role: %s\n", rsp.KilledUser, rsp.KilledUserRole)
	}
	return rsp.SessionStatus.AllConnected, nil
}
