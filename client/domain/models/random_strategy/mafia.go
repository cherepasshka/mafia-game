package random_strategy

import (
	// "bufio"
	// "os"
	"context"
	"fmt"
	"math/rand"

	"soa.mafia-game/client/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

type Mafia struct {
	models.BaseUser
}

func (user *Mafia) GetRole() proto.Roles {
	return proto.Roles_Mafia
}

func (user *Mafia) MakeNightMove(alive_players []string, client proto.MafiaServiceClient) error {
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this night")
		return nil
	}
	victim := user.Login
	for victim == user.Login {
		victim = alive_players[rand.Intn(len(alive_players))]
	}
	ctx := context.Background()
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

func (user *Mafia) VoteForMafia(alive_players []string, client proto.MafiaServiceClient) error {
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
