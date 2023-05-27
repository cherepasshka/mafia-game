package application

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"time"

	"soa.mafia-game/client/internal/utils/console"
	proto "soa.mafia-game/proto/mafia-game"
)

func validLogin(login string) bool {
	isAlpha := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	return isAlpha(login)
}

func (app *mafiaApplication) trySetLogin(login string) (*proto.ConnectToSessionResponse, error) {
	if !validLogin(login) {
		return &proto.ConnectToSessionResponse{Success: false}, nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	response, err := app.grpcClient.ConnectToSession(ctx, &proto.DefaultRequest{
		Login: login,
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
	login, err := console.Ask("Hello you! Welcome to Mafia game! Enter your login, please use only alphanumeric symbols")
	if err != nil {
		return "", nil, err
	}
	response, err := app.trySetLogin(login)
	if err != nil {
		return "", nil, err
	}
	for !response.Success {
		login, err = console.Ask("This login is busy or invalid (use only alphanumeric symbols). Please, take another")
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
	rsp, err := app.grpcClient.ListConnections(ctx, &proto.DefaultRequest{Login: login})
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
