package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	flag "github.com/spf13/pflag"
	"soa.mafia-game/client/application"
)

func main() {
	var port int
	var host string
	flag.IntVar(&port, "port", 9000, "specifies server port")
	flag.StringVar(&host, "host", "", "specifies server host")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer cancel()

	log.Println("Client running ...")
	app := application.New()
	go func() {
		err := app.Start(host, port)
		if err != nil {
			log.Fatal(err)
		}
		app.Stop()
		os.Exit(0)
	}()
	<-ctx.Done()

	app.Stop()

}
