package user_strategy

import (
	"context"
	"fmt"
	"math/rand"
	// "time"

	"soa.mafia-game/client/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

type Commissioner struct {
	models.BaseUser
	lastGuess string
}

func (user *Commissioner) GetRole() proto.Roles {
	return proto.Roles_Commissioner
}

func (user *Commissioner) MakeNightMove(ctx context.Context, players []string, client proto.MafiaServiceClient) error {
	// ctx, cancel := context.WithTimeout(ctx, time.Duration(time.Second))
	// defer cancel()
	// SetRandom()
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this night")
		return nil
	}
	suspected := user.Login
	for suspected == user.Login {
		suspected = players[rand.Intn(len(players))]
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
	
	// SetRandom()
	guess := user.Login
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this day vote")
		guess = "None"
	} else {
		if user.lastGuess == "" {
			for guess == user.Login {
				guess = alive_players[rand.Intn(len(alive_players))]
			}
		} else {
			guess = user.lastGuess
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
