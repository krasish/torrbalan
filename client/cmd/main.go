package main

import (
	"log"

	"github.com/krasish/torrbalan/client/internal/client"
	"github.com/krasish/torrbalan/client/internal/config"
)

func main() {
	cfg, err := config.NewClient("localhost:8080", 2)
	if err != nil {
		log.Fatalf("while processing configuration: %v", err)
	}
	cli := client.NewClient(cfg)
	if err = cli.Start(); err != nil {
		log.Fatalf("while starting client: %v", err)
	}
}
