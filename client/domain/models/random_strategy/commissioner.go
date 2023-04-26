package random_strategy

import (
	"context"
	"fmt"

	"math/rand"

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

// User interaction below
// func (user *Commissioner) MakeNightMove(players []string, client proto.MafiaServiceClient) error {
// 	if user.Status == Dead {
// 		fmt.Println("You are dead, so you skip this night")
// 		return nil
// 	}
// 	fmt.Printf("Enter suspected's login: ")
// 	reader := bufio.NewReader(os.Stdin)
// 	suspected, err := reader.ReadString('\n')
// 	if err != nil {
// 		return err
// 	}
// 	suspected = suspected[:len(suspected)-1]
// 	ctx := context.Background()
// 	response, err := client.MakeMove(ctx, &proto.MoveRequest{Target: suspected})
// 	if err != nil {
// 		return err
// 	}
// 	for !response.Accepted {
// 		fmt.Printf("%s\nEnter another login: ", response.Reason)
// 		suspected, err = reader.ReadString('\n')
// 		if err != nil {
// 			return err
// 		}
// 		suspected = suspected[:len(suspected)-1]
// 		response, err = client.MakeMove(ctx, &proto.MoveRequest{Target: suspected})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (user *Commissioner) MakeNightMove(players []string, client proto.MafiaServiceClient) error {
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this night")
		return nil
	}
	suspected := user.Login
	for suspected == user.Login {
		suspected = players[rand.Intn(len(players))]
	}
	ctx := context.Background()
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

func (user *Commissioner) VoteForMafia(alive_players []string, client proto.MafiaServiceClient) error {
	if user.Status == models.Dead {
		fmt.Println("You are dead, so you skip this day vote")
		return nil
	}
	guess := user.Login
	if user.lastGuess == "" {
		for guess == user.Login {
			guess = alive_players[rand.Intn(len(alive_players))]
		}
	} else {
		guess = user.lastGuess
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
