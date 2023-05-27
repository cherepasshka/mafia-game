package main

import (
	"log"
	"os"
	"os/signal"

	"soa.mafia-game/server/application"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	app, err := application.New()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		app.Start()
		stop <- os.Interrupt
	}()
	<-stop
	app.Stop()
}
