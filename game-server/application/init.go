package application

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	mafia_server "soa.mafia-game/game-server/adapter"
	proto "soa.mafia-game/proto/mafia-game"
)

type MafiaApplication struct {
	mafia_server *grpc.Server
}

func New() (*MafiaApplication, error) {
	app := &MafiaApplication{
		mafia_server: grpc.NewServer(),
	}
	brokerServers := os.Getenv("KAFKA_BROKER_URL")
	server, err := mafia_server.New(brokerServers)
	if err != nil {
		return nil, err
	}
	proto.RegisterMafiaServiceServer(app.mafia_server, server)
	return app, nil
}

func (app *MafiaApplication) Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	address := fmt.Sprintf(":%v", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Start serving at address %s", address)
	if err = app.mafia_server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func (app *MafiaApplication) Stop() {
	app.mafia_server.Stop()
}
