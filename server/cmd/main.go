package main

import (
	// "context"
	// "fmt"
	// "fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	proto "soa.mafia-game/proto/mafia-game"
	mafia_server "soa.mafia-game/server/adapter"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	proto.RegisterMafiaServiceServer(srv, mafia_server.New())
	log.Printf("Start serving")
	log.Fatalln(srv.Serve(lis))
}
