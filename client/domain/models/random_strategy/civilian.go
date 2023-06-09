package random_strategy

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"soa.mafia-game/client/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

type Civilian struct {
	models.BaseUser
}

func (user *Civilian) GetRole() proto.Roles {
	return proto.Roles_Civilian
}

func (user *Civilian) MakeNightMove(ctx context.Context, alive []string, client proto.MafiaServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second))
	defer cancel()
	if user.Status == models.Alive {
		client.MakeMove(ctx, &proto.MoveRequest{Login: user.Login})
	}
	return nil
}

func (user *Civilian) VoteForMafia(ctx context.Context, alive_players []string, client proto.MafiaServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second))
	defer cancel()
	SetRandom()
	guess := user.Login
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this day vote")
		guess = "None"
	} else {
		for guess == user.Login {
			guess = alive_players[rand.Intn(len(alive_players))]
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
