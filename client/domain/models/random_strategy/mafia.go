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

// User interaction below
// func (user *Mafia) MakeNightMove(client proto.MafiaServiceClient) error {
// 	if user.Status == Dead {
// 		fmt.Println("You are dead, so you skip this night")
// 		return nil
// 	}
// 	fmt.Printf("Enter victim's login: ")
// 	reader := bufio.NewReader(os.Stdin)
// 	victim, err := reader.ReadString('\n')
// 	if err != nil {
// 		return err
// 	}
// 	victim = victim[:len(victim)-1]
// 	ctx := context.Background()
// 	response, err := client.MakeMove(ctx, &proto.MoveRequest{Target: victim})
// 	if err != nil {
// 		return err
// 	}
// 	for !response.Accepted {
// 		fmt.Printf("%s\nEnter another login: ", response.Reason)
// 		victim, err = reader.ReadString('\n')
// 		if err != nil {
// 			return err
// 		}
// 		victim = victim[:len(victim)-1]
// 		response, err = client.MakeMove(ctx, &proto.MoveRequest{Target: victim})
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// Random below
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
