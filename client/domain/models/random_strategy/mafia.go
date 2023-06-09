package random_strategy

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"soa.mafia-game/client/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

type Mafia struct {
	models.BaseUser
}

func (user *Mafia) GetRole() proto.Roles {
	return proto.Roles_Mafia
}

func (user *Mafia) MakeNightMove(ctx context.Context, alive_players []string, client proto.MafiaServiceClient) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second))
	defer cancel()
	SetRandom()
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this night")
		return nil
	}
	victim := user.Login
	for victim == user.Login {
		victim = alive_players[rand.Intn(len(alive_players))]
	}
	response, err := client.MakeMove(ctx, &proto.MoveRequest{Login: user.Login, Target: victim})
	if err != nil {
		return err
	}
	for !response.Accepted {
		for victim == user.Login {
			victim = alive_players[rand.Intn(len(alive_players))]
		}
		response, err = client.MakeMove(ctx, &proto.MoveRequest{Login: user.Login, Target: victim})
		if err != nil {
			return err
		}
	}
	fmt.Printf("You murder %s\n", victim)
	return nil
}

func (user *Mafia) VoteForMafia(ctx context.Context, alive_players []string, client proto.MafiaServiceClient) error {
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
