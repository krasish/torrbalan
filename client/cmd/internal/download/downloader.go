package download

import (
	"fmt"
	"net"

	"github.com/krasish/torrbalan/client/cmd/internal/config"
)

type Downloader struct {
	serverAddress string
	q             chan string
	conn          net.Conn
}

func NewDownloader(c config.Downloader) (Downloader, error) {
	return Downloader{
		serverAddress: c.ServerAddress,
		q:             make(chan string, c.ConcurrentDownloads),
	}, nil
}

func (d Downloader) Start() error {
	var err error
	d.conn, err = net.Dial("tcp", d.serverAddress)
	if err != nil {
		return fmt.Errorf("while dialing server: %w", err)
	}

	d.registerToServer()
	for {
		filename := <-d.q
		d.processDownloading(filename)
	}
}

func (d Downloader) Download(filename string) {
	d.q <- filename
}

func (d Downloader) processDownloading(filename string) {
	//TODO: Consider what happens when two simultaneous
	// downloads for the same file are started.
}

func (d Downloader) registerToServer() {

	d.conn.Write()
}
