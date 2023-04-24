package application

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	domain_client "soa.mafia-game/client/domain/mafia-client"
	"soa.mafia-game/client/domain/models"

	proto "soa.mafia-game/proto/mafia-game"
)

type User struct {
	Login string
	Role  proto.Roles
}

type mafiaApplication struct {
	grpcClient *domain_client.Client
	//user       User
	login   string
	player  models.User
	players []string
	alive   []string
}

func New() *mafiaApplication {
	return &mafiaApplication{}
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
		login, err := reader.ReadString('\n')
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
	// role := proto.Roles_Undefined
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
			// role = m.Readiness.Role
			return m.Readiness, nil
		}
		fmt.Printf("%v %v at %v\n", m.Login, m.State, m.Time.AsTime())
	}
	return nil, nil
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
	app.players = readiness.Players
	fmt.Printf("Your session is ready, you are %v\n", role)

	app.player = models.MakeUser(login, role)
	app.alive = readiness.Players
	for {
		fmt.Printf("Alive members of this session are: ")
		for _, player := range app.alive {
			fmt.Printf("%s ", player)
		}
		fmt.Println()
		err = app.player.MakeNightMove(app.alive, app.grpcClient)
		if err != nil {
			return err
		}
		rsp, err := app.grpcClient.StartDay(context.TODO(), &proto.DayRequest{Login: app.login})
		if err != nil {
			return err
		}
		if rsp.Victim == app.login {
			fmt.Print("\tYou were killed this night!\n")
			app.player.SetStatus(models.Dead)
		} else {
			fmt.Printf("This night %s was murdured\n", rsp.Victim)
		}
		app.alive = rsp.Alive
		if len(app.alive) <= 1 {
			break
		}

	}
	fmt.Print("Game over\n")
	return nil
}

func (app *mafiaApplication) Stop() {
	app.grpcClient.LeaveSession(context.Background(), &proto.LeaveSessionRequest{User: &proto.User{Name: app.login}})
	app.grpcClient.Stop()
}
