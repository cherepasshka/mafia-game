package game

import (
	"context"
	"fmt"

	domain_client "soa.mafia-game/client/domain/grpc-client"
	"soa.mafia-game/client/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

func (game *Game) PrintAlive() {
	fmt.Printf("Alive members of this session are: ")
	for _, player := range game.alive {
		fmt.Printf("%s ", player)
	}
	fmt.Println()
}

func (game *Game) Start(ctx context.Context, grpcClient *domain_client.Client) error {
	for {
		game.PrintAlive()
		err := game.player.MakeNightMove(ctx, game.alive, grpcClient)
		if err != nil {
			return err
		}
		rsp, err := grpcClient.StartDay(ctx, &proto.DefaultRequest{Login: game.player.GetLogin()})
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
		fmt.Println("Start day")
		game.PrintAlive()
		err = game.player.VoteForMafia(ctx, game.alive, grpcClient)
		if err != nil {
			return err
		}
		rsp1, err := grpcClient.GetStatus(ctx, &proto.DefaultRequest{Login: game.player.GetLogin()})
		if err != nil {
			return err
		}
		game.alive = rsp1.Alive
		if !rsp1.GameStatus.Active {
			if rsp1.GameStatus.Winner == proto.Roles_Civilian {
				fmt.Printf("Civilians won!\n")
			} else {
				fmt.Printf("Mafia won =(\n")
			}
			grpcClient.ExitGameSession(ctx, &proto.DefaultRequest{Login: game.player.GetLogin()})
			break
		}
	}
	fmt.Print("Game over\n")
	return nil
}

// func (game *Game) Stop() {
// 	// game.player.Stop()
// }
