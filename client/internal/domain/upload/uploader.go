package upload

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/krasish/torrbalan/client/internal/logutil"
)

type Uploader struct {
	port  string
	q     chan struct{}
	files map[string]os.FileInfo
}

func NewUploader(concurrentUploads, port uint) Uploader {
	return Uploader{
		q:    make(chan struct{}, concurrentUploads),
		port: strconv.Itoa(int(port)),
	}
}

func (u Uploader) Start() error {
	listener, err := net.Listen("tcp", ":"+u.port)
	if err != nil {
		return fmt.Errorf("while getting a listener: %w", err)
	}
	log.Printf("Started listeling on %s\n", listener.Addr().String())

	for {
		u.acceptPeers(listener)
	}
}

func (u Uploader) acceptPeers(listener net.Listener) {
	<-u.q
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("An error ocurred while acceppting a connection: %v\n", err)
		u.q <- struct{}{}
		return
	}

	go u.processUploading(conn)
}

func (u Uploader) initialContract(conn net.Conn) (os.FileInfo, error) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	fileName, err := readWriter.ReadString('$')
	fileName = strings.TrimSuffix(fileName, "$")
	if err != nil {
		return nil, fmt.Errorf("while waiting for first response from peer: %v", err)
	}

	if _, ok := u.files[fileName]; !ok {
		if _, err := readWriter.WriteString("BAD$"); err != nil {
			logutil.LogOnErr(conn.Close)
			return nil, fmt.Errorf("while writing bad response to client: %w", err)
		}
		return nil, fmt.Errorf("%s asked for file %s which was not found", conn.RemoteAddr().String(), fileName)
	}

	if _, err := readWriter.WriteString("OK$"); err != nil {
		logutil.LogOnErr(conn.Close)
		return nil, fmt.Errorf("while writing ok response to client: %w", err)
	}
	return u.files[fileName], nil
}

//AddFile adds a file in current uploader and returns its name a SHA256 calculated for the file added.
func (u Uploader) AddFile(filePath string) (hash string, name string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", "", fmt.Errorf("while opening file: %v", err)
	}
	defer logutil.LogOnErr(f.Close)
	fileInfo, err := f.Stat()
	if err != nil {
		return "", "", fmt.Errorf("while getting file info: %v", err)
	}
	u.files[fileInfo.Name()] = fileInfo

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Printf("an error occured while calculating hash for file: %v", err)
	}
	return string(h.Sum(nil)), fileInfo.Name(), err
}

func (u Uploader) processUploading(conn net.Conn) {
	defer func() { u.q <- struct{}{} }()

	fileInfo, err := u.initialContract(conn)
	if err != nil {
		log.Printf("An error occurred while establishinng initial contract with %s: %v", conn.RemoteAddr().String(), err)
		return
	}

	os.OpenFile(fi)
}
