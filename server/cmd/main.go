package main

import (
	"os"
	"os/signal"

	"soa.mafia-game/server/application"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	app := application.New()
	go func() {
		app.Start()
		stop <- os.Interrupt
	}()
	<-stop
	app.Stop()
}
