package application

import (
	"context"
	"fmt"

	"soa.mafia-game/client/domain/game"
	domain_client "soa.mafia-game/client/domain/mafia-client"
	"soa.mafia-game/client/domain/models"
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
	fmt.Print("Hello! Welcome to Mafia game.\nEnter your login: ")
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
		role = readiness.Role
	}
	fmt.Printf("Your session is ready, you are %v\n", role)
	app.game = game.New(models.MakeUser(login, role), readiness.Players)
	return app.game.Start(context.Background(), app.grpcClient)
}

func (app *mafiaApplication) Stop() {
	app.grpcClient.LeaveSession(context.Background(), &proto.LeaveSessionRequest{User: &proto.User{Name: app.login}})
	app.grpcClient.Stop()
}
