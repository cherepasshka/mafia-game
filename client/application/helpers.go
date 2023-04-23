package application

import (
	"context"
	// "fmt"
	"time"

	proto "soa.mafia-game/proto/mafia-game"
)

func (app *mafiaApplication) setLogin(login string) (*proto.ConnectToSessionResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := app.grpcClient.ConnectToSession(ctx, &proto.User{
		Name: login,
		Role: proto.Roles_Undefined,
	})
	if err != nil {
		return nil, err
	}
	if response.Success {
		app.user.Login = login
	}
	return response, nil
}
