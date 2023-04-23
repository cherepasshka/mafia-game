package main

import (
	// "context"
	// "fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	// "time"

	// "google.golang.org/grpc"
	// domain_client "soa.mafia-game/client/domain/mafia-client"
	"soa.mafia-game/client/application"
	// proto "soa.mafia-game/proto/mafia-game"
	flag "github.com/spf13/pflag"
)

func main() {
	// ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	// defer cancel()
	var port int
	var host string
	flag.IntVar(&port, "port", 9000, "specifies server port")
	flag.StringVar(&host, "host", "", "specifies server host")
	flag.Parse()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT)

	log.Println("Client running ...")
	app := application.New()
	go func() {
		err := app.Start(host, port)
		if err != nil {
			log.Fatal(err)
		}
	}()

	<-ch

	app.Stop()

}
