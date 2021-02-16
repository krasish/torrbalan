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

//Uploader is responsible for uploading files to other clients. It also stores information
//about all the files which are currently possible to be downloaded from this client.
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

//Start should be started in a separate goroutine. It starts listening on the configured port
//and loops accepting clients. A new goroutine which serves clients is started for each client.
func (u Uploader) Start() {
	listener, err := net.Listen("tcp", ":"+u.port)
	if err != nil {
		log.Printf("an error occured while starting listener for uploader: %v", err)
		eofutil.TryWrite(u.stopChan)
		return
	}
	log.Printf("Started listeling on %s\n", listener.Addr().String())

	for {
		u.acceptPeers(listener)
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

//RemoveFile removes a file in current uploader.
func (u Uploader) RemoveFile(fileName string) {
	u.rw.Lock()
	defer u.rw.Unlock()
	delete(u.files, fileName)
}

func (u Uploader) acceptPeers(listener net.Listener) {
	u.q <- struct{}{}
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("An error ocurred while acceppting a connection: %v\n", err)
		u.q <- struct{}{}
		return
	}

	go u.processUploading(conn)
}

//initialContract should be called before processUploading.
//initialContract sends/reads messages to/form client on conn by which both clients
//can determine which file is about to get transferred and determine whether
//that is possible, returning an error if not.
func (u Uploader) initialContract(conn net.Conn) (string, error) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	h := eofutil.LoggingEOFHandler{DestName: conn.RemoteAddr().String()}

	fileName, err := eofutil.ReadCheckEOF(readWriter.Reader, '$', h)
	if err != nil {
		return "", fmt.Errorf("while waiting for first response from peer: %v", err)
	}
	fileName = strings.TrimSuffix(fileName, "$")

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

//processUploading loops reading from file and writing to the client.
//Both the net.Conn and *file.File are closed when downloading is done or an error occurred.
func (u Uploader) processUploading(conn net.Conn) {
	defer func() { <-u.q }()
	defer logutil.LogOnErr(conn.Close)
	remoteAddr := conn.RemoteAddr().String()

	filePath, err := u.initialContract(conn)
	if err != nil {
		log.Printf("An error occurred while establishinng initial contract with %s: %v", remoteAddr, err)
		return
	}
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("An error occurred while opening file %q: %v", filePath, err)
		return
	}
	defer logutil.LogOnErr(file.Close)

	reader, writer := bufio.NewReader(file), bufio.NewWriter(conn)
	errorMessages := [5]string{
		fmt.Sprintf("Stopping upload of file %q to %s, EOF read.", file.Name(), remoteAddr),
		fmt.Sprintf("An error occurred while reading from file %q: %%v", file.Name()),
		fmt.Sprintf("Stopping upload of file %q from %s. EOF read through connection.", file.Name(), remoteAddr),
		fmt.Sprintf("An error occurred while writing to %s: %%v", remoteAddr),
		fmt.Sprintf("final flush to peer %s failed", remoteAddr),
	}

	download.ReadWriteLoop(reader, writer, errorMessages)
}
