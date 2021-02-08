package download

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path"

	"github.com/krasish/torrbalan/client/internal/logutil"
)

const BUFFER_SIZE = 1024

type Info struct {
	Filename    string
	PeerAddress string
	PathToSave  string
}

type Downloader struct {
	q chan Info
}

func NewDownloader(concurrentDownloads uint, conn net.Conn) Downloader {
	return Downloader{
		q: make(chan Info, concurrentDownloads),
	}
}

func (d Downloader) Start() {
	for {
		info := <-d.q
		rootPath := path.Clean(info.PathToSave)
		filePath := rootPath + "/" + info.Filename
		file, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL, 0644)
		if err != nil {
			fmt.Printf("Could not create file: %v\n", err)
			continue
		}
		conn, err := d.connectToPeer(info)
		if err != nil {
			fmt.Printf("Could not connect: %v\n", err)
			logutil.LogOnErr(file.Close)
			continue
		}
		d.processDownloading(file, conn)
	}
}

func (d Downloader) Download(info Info) {
	d.q <- info
}

func (d Downloader) connectToPeer(info Info) (net.Conn, error) {
	conn, err := net.Dial("tcp", info.PeerAddress)
	if err != nil {
		return nil, fmt.Errorf("while trying to connecto to %s: %w", info.PeerAddress, err)
	}
	return conn, nil
}

func (d Downloader) initialContract(conn net.Conn, fileName string) error {
	readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	_, err := readWriter.WriteString(fileName + "$")
	if err != nil {
		return fmt.Errorf("while writing filename: %v", err)
	}
	resp, err := readWriter.ReadString('$')
	if err != nil {
		return fmt.Errorf("while waiting for first response from peer: %v", err)
	}
	if resp != "OK$" {
		return fmt.Errorf("file %q cannot be downloaded from %s", fileName, conn.RemoteAddr().String())
	}
	return nil
}

func (d Downloader) processDownloading(file *os.File, conn net.Conn) {
	defer logutil.LogAllOnErr(file.Close, conn.Close)

	if err := d.initialContract(conn, file.Name()); err != nil {
		log.Printf("An error occurred while establishinng initial contract with %s: %v", conn.RemoteAddr().String(), err)
		return
	}

	bytes := make([]byte, BUFFER_SIZE)
	reader, writer := bufio.NewReader(conn), bufio.NewWriter(file)

	for {
		n, err := reader.Read(bytes)
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("An error occurred while reading from %s: %v", conn.RemoteAddr().String(), err)
			return
		}
		_, err = writer.Write(bytes[:n])
		if err != nil {
			log.Printf("An error occurred while writing in file %s: %v", file.Name(), err)
			return
		}
	}
}
