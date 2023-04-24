package application

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	// "time"

	domain_client "soa.mafia-game/client/domain/mafia-client"

	proto "soa.mafia-game/proto/mafia-game"
)

type User struct {
	Login string
	Role  proto.Roles
}

type mafiaApplication struct {
	grpcClient *domain_client.Client
	user       User
}

func New() *mafiaApplication {
	return &mafiaApplication{}
}

func (app *mafiaApplication) SetLogin() {

}

func (app *mafiaApplication) Start(host string, port int) error {
	fmt.Print("Hello! Welcome to Mafia game.\nEnter your login: ")
	var err error
	app.grpcClient, err = domain_client.New(host, port)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	login, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	login = login[:len(login)-1]
	response, err := app.setLogin(login)
	if err != nil {
		return err
	}
	for !response.Success {
		fmt.Print("This login is busy. Please, take another: ")
		login, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		response, err = app.setLogin(login)
		if err != nil {
			return err
		}
	}
	app.user.Role = proto.Roles_Undefined
	if response.Readiness.SessionReady {
		app.user.Role = response.Readiness.Role
	} else {
		ctx := context.Background()
		rsp, err := app.grpcClient.ListConnections(ctx, &proto.ListConnectionsRequest{Login: login})
		if err != nil {
			return err
		}
		for {
			m, err := rsp.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			if m.Readiness.SessionReady {
				app.user.Role = m.Readiness.Role
				break
			}
			fmt.Printf("%v %v at %v\n", m.Login, m.State, m.Time.AsTime())
		}
	}
	fmt.Printf("Your session is ready, you are %v\n", app.user.Role)
	return nil
}

func (app *mafiaApplication) Stop() {
	app.grpcClient.LeaveSession(context.Background(), &proto.LeaveSessionRequest{User: &proto.User{Name: app.user.Login}})
	app.grpcClient.Stop()
}
