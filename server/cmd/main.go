package main

import (
	"log"

	"github.com/krasish/torrbalan/server/internal/config"
	"github.com/krasish/torrbalan/server/internal/server"
)

func main() {
	port := "8080"
	serverConfig, err := config.NewServer(port, 100)
	if err != nil {
		log.Fatal(err)
	}
	s := server.NewServer(serverConfig)
	err = s.Run()
	if err != nil {
		log.Fatalf("could not start upload: %v", err)
	}
}
