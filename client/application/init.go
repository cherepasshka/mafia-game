package application

import (
	"context"
	"errors"
	"fmt"

	"soa.mafia-game/client/domain/game"
	domain_client "soa.mafia-game/client/domain/grpc-client"
	"soa.mafia-game/client/domain/models"
	"soa.mafia-game/client/internal/utils/console"
	proto "soa.mafia-game/proto/mafia-game"
)

type mafiaApplication struct {
	grpcClient *domain_client.Client
	login      string
	game       *game.Game
}

func New() *mafiaApplication {
	return &mafiaApplication{}
}

func (app *mafiaApplication) Start(host string, port int) error {
	var err error
	app.grpcClient, err = domain_client.New(host, port)
	if err != nil {
		return err
	}
	login, readiness, err := app.SetLogin()
	if err != nil {
		return err
	}
	app.login = login

	role := proto.Roles_Undefined
	if readiness.SessionReady {
		role = readiness.Role
	} else {
		if readiness, err = app.WaitForSession(login); err != nil {
			return err
		}
		if readiness == nil {
			return nil
		}
		role = readiness.Role
	}
	app.game = game.New(models.MakeUser(login, role, readiness.SessionId, readiness.Partition), readiness.Players)
	for {
		fmt.Printf("Your session is ready, you are %v\n", role)
		if err = app.game.Start(context.Background(), app.grpcClient); err != nil {
			if errors.Is(game.ErrSessionInterrupted, err) {
				fmt.Println("You was disconnected because some players left your game session")
				app.grpcClient.ExitGameSession(context.Background(), &proto.DefaultRequest{Login: login})
			} else {
				return err
			}
		}
		proceed, err := console.AskPrompt("Do you want to continue?", []string{"yes", "no"})
		if err != nil {
			return err
		}
		if proceed == "no" {
			break
		}
		readiness, err = app.WaitForSession(app.login)
		if err != nil {
			return err
		}
		if readiness == nil {
			break
		}
		role = readiness.Role
		app.game = game.New(models.MakeUser(login, readiness.Role, readiness.SessionId, readiness.Partition), readiness.Players)
	}
	return nil
}

func (app *mafiaApplication) Stop(ctx context.Context) {
	if app.game != nil {
		app.game.Stop()
	}
	if app.grpcClient != nil {
		app.grpcClient.LeaveSession(ctx, &proto.DefaultRequest{Login: app.login})
		app.grpcClient.Stop()
	}
}
