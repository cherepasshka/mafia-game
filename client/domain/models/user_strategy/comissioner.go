package user_strategy

import (
	"context"
	"fmt"

	"soa.mafia-game/client/domain/chat"
	"soa.mafia-game/client/domain/models/user"
	"soa.mafia-game/client/internal/utils/console"
	proto "soa.mafia-game/proto/mafia-game"
)

type Commissioner struct {
	models.CommunicatorUser
	lastGuess   string
	ChatService *chat.ChatService
}

func (user *Commissioner) GetRole() proto.Roles {
	return proto.Roles_Commissioner
}

func (user *Commissioner) MakeNightMove(ctx context.Context, alive_players []string, client proto.MafiaServiceClient) error {
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this night")
		return nil
	}
	suspected := user.Login
	for suspected == user.Login {
		suspected, _ = console.AskPrompt("Select suspect", user.ExcludeFromAliveList(alive_players))
	}
	response, err := client.MakeMove(ctx, &proto.MoveRequest{Login: user.Login, Target: suspected})
	if err != nil {
		return err
	}
	if response.Accepted {
		fmt.Printf("You suspected %s correct, this user is mafia\n", suspected)
		user.lastGuess = suspected
	} else {
		fmt.Printf("You suspected %s wrong, this user is not mafia\n", suspected)
		user.lastGuess = ""
	}
	return nil
}

func (user *Commissioner) VoteForMafia(ctx context.Context, alive_players []string, client proto.MafiaServiceClient) error {
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
		return err
	}
	if rsp.KilledUser == user.Login {
		fmt.Println("Most voted for you")
	} else {
		fmt.Printf("Most voted for %s, this user had role: %s\n", rsp.KilledUser, rsp.KilledUserRole)
	}
	return nil
}
