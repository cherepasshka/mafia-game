package models

import (
	"bufio"
	"context"
	"fmt"
	"os"

	proto "soa.mafia-game/proto/mafia-game"
)

type Mafia struct {
	BaseUser
}

func (user *Mafia) GetRole() proto.Roles {
	return proto.Roles_Mafia
}

func (user *Mafia) MakeNightMove(client proto.MafiaServiceClient) error {
	if user.Status == Dead {
		fmt.Println("You are dead, so you skip this night")
		return nil
	}
	fmt.Printf("Enter victim's login: ")
	reader := bufio.NewReader(os.Stdin)
	victim, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	victim = victim[:len(victim)-1]
	ctx := context.Background()
	response, err := client.MakeMove(ctx, &proto.MoveRequest{Target: victim})
	if err != nil {
		return err
	}
	for !response.Accepted {
		fmt.Printf("%s\nEnter another login: ", response.Reason)
		victim, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		victim = victim[:len(victim)-1]
		response, err = client.MakeMove(ctx, &proto.MoveRequest{Target: victim})
		if err != nil {
			return err
		}
	}
	return nil
}
