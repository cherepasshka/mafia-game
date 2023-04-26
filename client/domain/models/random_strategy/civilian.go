package random_strategy

import (
	"context"
	"fmt"
	"math/rand"

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

func (user *Civilian) VoteForMafia(alive_players []string, client proto.MafiaServiceClient) error {
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this day vote")
		return nil
	}
	guess := user.Login
	for guess == user.Login {
		guess = alive_players[rand.Intn(len(alive_players))]
	}
	rsp, err := client.VoteForMafia(context.Background(), &proto.VoteForMafiaRequest{Login: user.Login, MafiaGuess: guess})
	if err != nil {
		return err
	}
	fmt.Printf("You voted for %s\n", guess)
	if rsp.KilledUser == user.Login {
		fmt.Println("Most voted for you")
	} else {
		fmt.Printf("Most voted for %s, this user had role: %s\n", rsp.KilledUser, rsp.KilledUserRole)
	}
	return nil
}
