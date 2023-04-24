package game

import (
	"context"
	"fmt"

	domain_client "soa.mafia-game/client/domain/mafia-client"
	"soa.mafia-game/client/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

func (game *Game) Start(grpcClient *domain_client.Client) error {
	for {
		fmt.Printf("Alive members of this session are: ")
		for _, player := range game.alive {
			fmt.Printf("%s ", player)
		}
		fmt.Println()
		err := game.player.MakeNightMove(game.alive, grpcClient)
		if err != nil {
			return err
		}
		rsp, err := grpcClient.StartDay(context.TODO(), &proto.DayRequest{Login: game.player.GetLogin()})
		if err != nil {
			return err
		}
		if rsp.Victim == game.player.GetLogin() {
			fmt.Print("\tYou were killed this night!\n")
			game.player.SetStatus(models.Dead)
		} else {
			fmt.Printf("This night %s was murdured\n", rsp.Victim)
		}
		game.alive = rsp.Alive
		if !rsp.GameStatus.Active {
			if rsp.GameStatus.Winner == proto.Roles_Civilian {
				fmt.Printf("Civilians won!\n")
			} else {
				fmt.Printf("Mafia won =(\n")
			}
			break
		}
	}
	fmt.Print("Game over\n")
	return nil
}
