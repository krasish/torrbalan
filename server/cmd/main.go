package main

import (
	"log"
	"os"
	"strconv"

	"github.com/krasish/torrbalan/server/internal/config"
	"github.com/krasish/torrbalan/server/internal/server"
)

func main() {
	port := getPort()
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

func getPort() (port string) {
	if len(os.Args) < 2 {
		log.Fatalf("port was not passed")
	}
	port = os.Args[1]
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1024 {
		log.Fatalf("port is invalid")
	}
	return
}
