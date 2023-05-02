package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	flag "github.com/spf13/pflag"
	"soa.mafia-game/client/application"
)

func main() {
	var port int
	var host string
	flag.IntVar(&port, "port", 9000, "specifies server port")
	flag.StringVar(&host, "host", "", "specifies server host")
	flag.Parse()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	log.Println("Client running ...")
	app := application.New()
	go func() {
		log.Printf("Connecting to server %s:%v", host, port)
		err := app.Start(host, port)
		if err != nil {
			log.Fatal(err)
		}
		stop <- os.Interrupt
	}()
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Stop(ctx)

}
