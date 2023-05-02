package main

import (
	// "fmt"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	proto "soa.mafia-game/proto/mafia-game"
	mafia_server "soa.mafia-game/server/adapter"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	address := fmt.Sprintf(":%v", port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	proto.RegisterMafiaServiceServer(srv, mafia_server.New())
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	go func() {
		log.Printf("Start serving at address %s", address)
		if err = srv.Serve(lis); err != nil {
			log.Fatal(err)
		}
		stop <- os.Interrupt
	}()
	<-stop
	srv.Stop()
}
