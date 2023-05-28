package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"soa.mafia-game/server/application"
)

func main() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	app, err := application.New()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		app.Start()
		stop <- syscall.SIGINT
	}()
	<-stop
	app.Stop()
}
