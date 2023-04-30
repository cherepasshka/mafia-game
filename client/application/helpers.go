package application

import (
	// "bufio"
	"context"
	"fmt"
	"io"
	// "os"
	"time"

	"soa.mafia-game/client/internal/utils/console"
	proto "soa.mafia-game/proto/mafia-game"
)

func (app *mafiaApplication) trySetLogin(login string) (*proto.ConnectToSessionResponse, error) {
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
		app.login = login
	}
	return response, nil
}

func (app *mafiaApplication) SetLogin() (string, *proto.SessionReadiness, error) {
	login, err := console.Ask("Hello you! Welcome to Mafia game! Enter your login")
	if err != nil {
		return "", nil, err
	}
	response, err := app.trySetLogin(login)
	if err != nil {
		return "", nil, err
	}
	for !response.Success {
		login, err = console.Ask("This login is busy. Please, take another")
		if err != nil {
			return "", nil, err
		}
		response, err = app.trySetLogin(login)
		if err != nil {
			return "", nil, err
		}
	}
	return login, response.Readiness, nil
}

func (app *mafiaApplication) WaitForSession(login string) (*proto.SessionReadiness, error) {
	ctx := context.Background()
	rsp, err := app.grpcClient.ListConnections(ctx, &proto.ListConnectionsRequest{Login: login})
	if err != nil {
		return nil, err
	}
	for {
		m, err := rsp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if m.Readiness.SessionReady {
			rsp.CloseSend()
			return m.Readiness, nil
		}
		fmt.Printf("%v %v at %v\n", m.Login, m.State, m.Time.AsTime())
	}
	return nil, nil
}
