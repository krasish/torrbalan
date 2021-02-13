package client

import (
	"fmt"
	"log"
	"net"

	"github.com/krasish/torrbalan/client/internal/command"

	"github.com/krasish/torrbalan/client/internal/domain/connection"

	"github.com/krasish/torrbalan/client/internal/domain/upload"

	"github.com/krasish/torrbalan/client/internal/config"

	"github.com/krasish/torrbalan/client/internal/domain/download"
)

type Client struct {
	config.Client
	d download.Downloader
	u upload.Uploader
	c *connection.ServerCommunicator
	p *command.Processor
}

func NewClient(cfg config.Client) Client {
	return Client{Client: cfg}
}

func (c Client) Start() error {
	conn, err := net.Dial("tcp", c.ServerAddress)
	if err != nil {
		return fmt.Errorf("while dialing server: %w", err)
	}
	stopChan := make(chan struct{})
	c.c = connection.NewServerCommunicator(conn, stopChan)
	c.p = command.NewProcessor(c.c)
	c.p.Register()

	c.d = download.NewDownloader(c.ConcurrentDownloads, conn)
	c.u = upload.NewUploader(c.ConcurrentUploads, c.Port, stopChan)

	go c.c.Listen()
	go c.d.Start()
	go c.u.Start()

	<-stopChan
	log.Println("Stop signal received. Shutting down...")
	return nil
}
