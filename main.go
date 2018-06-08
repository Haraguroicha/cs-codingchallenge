package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Haraguroicha/cs-codingchallenge/Goddit"
	_ "github.com/heroku/x/hmetrics/onload"
)

type service interface {
	Start()
	Stop()
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var services = []service{}

	services = append(services, Goddit.NewService(false, port, "conf/config.yaml"))

	// Start the services
	for _, s := range services {
		s.Start()
	}

	// Handle SIGINT and SIGTERM
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	log.Println(<-ch)

	// Stop the services gracefully
	for _, s := range services {
		s.Stop()
	}

	// Wait exit
	log.Println("Bye.")
}
