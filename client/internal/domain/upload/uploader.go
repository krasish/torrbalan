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
	"sync"

	"github.com/krasish/torrbalan/client/pkg/eofutil"

	"github.com/krasish/torrbalan/client/internal/domain/download"

	"github.com/krasish/torrbalan/client/internal/logutil"
)

type Uploader struct {
	port     string
	q        chan struct{}
	files    map[string]string
	rw       *sync.RWMutex
	stopChan chan<- struct{}
}

func NewUploader(concurrentUploads, port uint, stopChan chan<- struct{}) Uploader {
	return Uploader{
		q:        make(chan struct{}, concurrentUploads),
		port:     strconv.Itoa(int(port)),
		files:    make(map[string]string),
		rw:       &sync.RWMutex{},
		stopChan: stopChan,
	}
}

func (u Uploader) Start(stopChan <-chan struct{}) {
	listener, err := net.Listen("tcp", ":"+u.port)
	if err != nil {
		log.Printf("an error occured while starting listener for uploader: %v", err)
		eofutil.TryWrite(u.stopChan)
		return
	}
	defer logutil.LogOnErr(listener.Close)
	log.Printf("Started listeling on %s\n", listener.Addr().String())

	for {
		select {
		case <-stopChan:
			fmt.Println("STOP")
			break
		default:
			u.acceptPeers(listener)
		}
	}
}

//AddFile adds a file in current uploader and returns a SHA256 calculated for the file added and its name.
func (u Uploader) AddFile(filePath string) (name string, hash string, err error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", "", fmt.Errorf("while opening file: %v", err)
	}
	defer logutil.LogOnErr(f.Close)
	fileInfo, err := f.Stat()
	if err != nil {
		return "", "", fmt.Errorf("while getting file info: %v", err)
	}
	u.rw.Lock()
	u.files[fileInfo.Name()] = filePath
	u.rw.Unlock()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Printf("an error occured while calculating hash for file: %v", err)
	}
	return fileInfo.Name(), string(h.Sum(nil)), err
}

//AddFile adds a file in current uploader and returns a SHA256 calculated for the file added and its name.
func (u Uploader) RemoveFile(fileName string) {
	u.rw.Lock()
	defer u.rw.Unlock()
	delete(u.files, fileName)
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

func (u Uploader) initialContract(conn net.Conn) (string, error) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	fileName, err := readWriter.ReadString('$')
	fileName = strings.TrimSuffix(fileName, "$")
	if err != nil {
		return "", fmt.Errorf("while waiting for first response from peer: %v", err)
	}
	h := eofutil.LoggingEOFHandler{conn.RemoteAddr().String()}

	u.rw.RLock()
	if _, ok := u.files[fileName]; !ok {
		u.rw.RUnlock()
		if err := eofutil.WriteCheckEOF(readWriter.Writer, "BAD$", h); err != nil {
			logutil.LogOnErr(conn.Close)
			return "", fmt.Errorf("while writing bad response to client: %w", err)
		}
		return "", fmt.Errorf("%s asked for file %s which was not found", conn.RemoteAddr().String(), fileName)
	}
	u.rw.RUnlock()

	if err := eofutil.WriteCheckEOF(readWriter.Writer, "OK$", h); err != nil {
		logutil.LogOnErr(conn.Close)
		return "", fmt.Errorf("while writing ok response to client: %w", err)
	}
	return u.files[fileName], nil
}

func (u Uploader) processUploading(conn net.Conn) {
	defer func() { u.q <- struct{}{} }()
	defer logutil.LogOnErr(conn.Close)

	filePath, err := u.initialContract(conn)
	if err != nil {
		log.Printf("An error occurred while establishinng initial contract with %s: %v", conn.RemoteAddr().String(), err)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("An error occurred while opening file %q: %v", filePath, err)
		return
	}
	defer logutil.LogOnErr(file.Close)

	reader, writer := bufio.NewReader(file), bufio.NewWriter(conn)
	errorMessages := [3]string{
		fmt.Sprintf("Stopping upload of file %q to %s", file.Name(), conn.RemoteAddr().String()),
		fmt.Sprintf("An error occurred while reading from file %q: %%v", file.Name()),
		fmt.Sprintf("An error occurred while writing to %s: %%v", conn.RemoteAddr().String()),
	}

	download.ReadWriteLoop(reader, writer, errorMessages)
}
