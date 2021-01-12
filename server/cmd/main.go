package main

import (
	"log"

	"github.com/krasish/torrbalan/server/internal/config"
	"github.com/krasish/torrbalan/server/internal/server"
)

func main() {
	serverConfig, err := config.NewServer("8080", 100)
	if err != nil {
		log.Fatal(err)
	}
	s := server.NewServer(serverConfig)
	err = s.Run()
	if err != nil {
		log.Fatalf("could not start uploader: %v", err)
	}
}
