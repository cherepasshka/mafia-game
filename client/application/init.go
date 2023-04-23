package application

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	// "time"

	domain_client "soa.mafia-game/client/domain/mafia-client"

	empty "github.com/golang/protobuf/ptypes/empty"
	proto "soa.mafia-game/proto/mafia-game"
)

type User struct {
	Login string
}

type mafiaApplication struct {
	grpcClient *domain_client.Client
	user       User
}

func New() *mafiaApplication {
	return &mafiaApplication{}
}

func (app *mafiaApplication) Start(host string, port int) error {
	fmt.Print("Hello! Welcome to Mafia game.\nEnter your login: ")
	reader := bufio.NewReader(os.Stdin)
	login, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	login = login[:len(login)-1]
	app.grpcClient, err = domain_client.New(host, port)
	if err != nil {
		return err
	}
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
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	// defer cancel()
	ctx := context.Background()
	rsp, err := app.grpcClient.ListConnections(ctx, &empty.Empty{})
	if err != nil {
		return err
	}
	for {
		m, e := rsp.Recv()
		if e == io.EOF {
			break
		}
		if e != nil {
			return e
		}
		fmt.Printf("%v %v at %v\n", m.Login, m.State, m.Time.AsTime())
	}
	return nil
}

func (app *mafiaApplication) Stop() {
	app.grpcClient.LeaveSession(context.Background(), &proto.LeaveSessionRequest{User: &proto.User{Name: app.user.Login}})
	app.grpcClient.Stop()
}
