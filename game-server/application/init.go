package application

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"google.golang.org/grpc"

	mafia_server "soa.mafia-game/game-server/adapter/grpc"
	http_mafia "soa.mafia-game/game-server/adapter/http"
	usersdb "soa.mafia-game/game-server/domain/models/storage"
	proto "soa.mafia-game/proto/mafia-game"
)

type MafiaApplication struct {
	mafia_server *grpc.Server
	http_mafia   *http_mafia.HttpHandler
	users        *usersdb.Storage
	http_server  *http.Server
}

func New() (*MafiaApplication, error) {
	app := &MafiaApplication{
		mafia_server: grpc.NewServer(),
		users:        usersdb.New(),
	}
	app.http_mafia = http_mafia.New(app.users)
	brokerServers := os.Getenv("KAFKA_BROKER_URL")
	server, err := mafia_server.New(app.users, brokerServers)
	if err != nil {
		return nil, err
	}
	proto.RegisterMafiaServiceServer(app.mafia_server, server)
	return app, nil
}

func (app *MafiaApplication) Start() {
	http_address := os.Getenv("HTTP_ADDRESS")
	router := chi.NewRouter()
	app.http_server = &http.Server{
		Addr:    http_address,
		Handler: router,
	}

	app.http_mafia.AddHandlers(router)
	go func() {
		log.Printf("Starting http server at %v", http_address)
		if err := app.http_server.ListenAndServe(); err != nil {
			log.Printf("Failed to start http: %v", err)
		}
	}()

	game_address := os.Getenv("MAFIA_GAME_ADDRESS")
	lis, err := net.Listen("tcp", game_address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("Start serving at address %s", game_address)
	if err = app.mafia_server.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func (app *MafiaApplication) Stop() {
	app.mafia_server.Stop()
}
