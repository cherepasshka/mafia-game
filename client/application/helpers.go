package application

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"

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
	reader := bufio.NewReader(os.Stdin)
	login, err := reader.ReadString('\n')
	if err != nil {
		return "", nil, err
	}
	login = login[:len(login)-1]
	response, err := app.trySetLogin(login)
	if err != nil {
		return "", nil, err
	}
	for !response.Success {
		fmt.Print("This login is busy. Please, take another: ")
		login, err = reader.ReadString('\n')
		if err != nil {
			return "", nil, err
		}
		login = login[:len(login)-1]
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
			return m.Readiness, nil
		}
		fmt.Printf("%v %v at %v\n", m.Login, m.State, m.Time.AsTime())
	}
	return nil, nil
}
