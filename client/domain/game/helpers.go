package game

import (
	"context"
	"fmt"

	domain_client "soa.mafia-game/client/domain/grpc-client"
	"soa.mafia-game/client/domain/models/user"
	proto "soa.mafia-game/proto/mafia-game"
)

var ErrSessionInterrupted = fmt.Errorf("Some players left game session")

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
		isActive, err := game.player.MakeNightMove(ctx, game.alive, grpcClient)
		if err != nil {
			return err
		}
		if !isActive {
			return ErrSessionInterrupted
		}
		rsp, err := grpcClient.StartDay(ctx, &proto.DefaultRequest{Login: game.player.GetLogin()})
		if err != nil {
			return err
		}
		if !rsp.SessionStatus.AllConnected {
			return ErrSessionInterrupted
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
		isActive, err = game.player.VoteForMafia(ctx, game.alive, grpcClient)
		if err != nil {
			return err
		}
		if !isActive {
			return ErrSessionInterrupted
		}
		status_rsp, err := grpcClient.GetStatus(ctx, &proto.DefaultRequest{Login: game.player.GetLogin()})
		if err != nil {
			return err
		}
		if !status_rsp.SessionStatus.AllConnected {
			return ErrSessionInterrupted
		}
		game.alive = status_rsp.Alive
		if !status_rsp.GameStatus.Active {
			if status_rsp.GameStatus.Winner == proto.Roles_Civilian {
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

func (game *Game) Stop() {
	game.player.Stop()
}
