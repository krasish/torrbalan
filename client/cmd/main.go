package main

import (
	"log"
	"os"

	"github.com/krasish/torrbalan/client/internal/client"
	"github.com/krasish/torrbalan/client/internal/config"
	"github.com/phayes/freeport"
)

func main() {
	serverAddress := os.Args[1]
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatalf("cannot start client: while getting free port: %v", err)
	}
	cfg, err := config.NewClient(serverAddress, 2, 2, uint(port))
	if err != nil {
		log.Fatalf("cannot start client: while processing configuration: %v", err)
	}
	log.Printf("Starting torrbalan client listening on port %d. Server is at %s", port, serverAddress)
	cli := client.NewClient(cfg)
	if err = cli.Start(); err != nil {
		log.Fatalf("while starting client: %v", err)
	}
}
