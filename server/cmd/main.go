package main

import (
	"github.com/krasish/torrbalan/server/internal/config"
	"github.com/krasish/torrbalan/server/internal/server"
	"log"
)

func main(){
	serverConfig := config.NewServer("8080", 100)
	s := server.NewServer(serverConfig)
	err := s.Run()
	if err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
