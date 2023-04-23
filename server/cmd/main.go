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
	// c := make(chan int, 100)
	// c <- 23
	// fmt.Printf("%v, %v\n", c, <- c)
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	proto.RegisterMafiaServiceServer(srv, mafia_server.New())
	log.Printf("Start serving")
	log.Fatalln(srv.Serve(lis))
}
