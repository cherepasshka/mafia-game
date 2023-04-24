package models

import (
	"bufio"
	"context"
	"fmt"
	"os"

	proto "soa.mafia-game/proto/mafia-game"
)

type Commissioner struct {
	BaseUser
}

func (user *Commissioner) GetRole() proto.Roles {
	return proto.Roles_Commissioner
}

func (user *Commissioner) MakeNightMove(client proto.MafiaServiceClient) error {
	if user.Status == Dead {
		fmt.Println("You are dead, so you skip this night")
		return nil
	}
	fmt.Printf("Enter suspected's login: ")
	reader := bufio.NewReader(os.Stdin)
	suspected, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	suspected = suspected[:len(suspected)-1]
	ctx := context.Background()
	response, err := client.MakeMove(ctx, &proto.MoveRequest{Target: suspected})
	if err != nil {
		return err
	}
	for !response.Accepted {
		fmt.Printf("%s\nEnter another login: ", response.Reason)
		suspected, err = reader.ReadString('\n')
		if err != nil {
			return err
		}
		suspected = suspected[:len(suspected)-1]
		response, err = client.MakeMove(ctx, &proto.MoveRequest{Target: suspected})
		if err != nil {
			return err
		}
	}
	return nil
}
