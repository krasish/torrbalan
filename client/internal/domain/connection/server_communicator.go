package connection

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/krasish/torrbalan/client/pkg/eofutil"
)

const (
	GetOwnersPattern       = "GET_OWNERS %s\n"
	UploadPattern          = "UPLOAD %s %q\n"
	StopUploadPattern      = "STOP_UPLOAD %s\n"
	DisconnectRequest      = "DISCONNECT\n"
	UserAlreadyExists      = "UAE\n"
	RegisteredSuccessfully = "RS\n"
	ServerMessagesColour   = "\n\033[36m"
	resetColour            = "\033[0m>"
)

type ServerCommunicator struct {
	conn     net.Conn
	stopChan chan struct{}
}

func NewServerCommunicator(conn net.Conn, stopChan chan struct{}) *ServerCommunicator {
	return &ServerCommunicator{conn: conn, stopChan: stopChan}
}

func (c ServerCommunicator) Listen(stopChan <-chan struct{}) {
	for {
		select {
		case <-stopChan:
			break
		default:
			reader := bufio.NewReader(c.conn)
			readString, err := eofutil.ReadServerCheckEOF(reader, '\n', c.stopChan)
			if err != nil {
				log.Printf("cannot read from server: %v", err)
			}
			c.printServerMessage(readString)
		}
	}
}

func (c ServerCommunicator) Register(username string, port uint) error {
	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))

	if err := eofutil.WriteServerCheckEOF(rw.Writer, c.concatUsernamePort(username, port), c.stopChan); err != nil {
		return fmt.Errorf("while writing to server: %w", err)
	}

	resp, err := eofutil.ReadServerCheckEOF(rw.Reader, '\n', c.stopChan)
	if err != nil {
		return fmt.Errorf("while reading from server: %w", err)
	}
	if resp == UserAlreadyExists {
		return errors.New("username already exists")
	} else if resp == RegisteredSuccessfully {
		return nil
	} else {
		return fmt.Errorf("could not register. Server responded: %s", resp)
	}
}

func (c ServerCommunicator) GetOwners(filename string) {
	rw := bufio.NewWriter(c.conn)

	err := eofutil.WriteServerCheckEOF(rw, fmt.Sprintf(GetOwnersPattern, filename), c.stopChan)
	if err != nil {
		log.Printf("an error while getting owners from server: %v\n", err)
	}
}

func (c ServerCommunicator) StartUploading(fileName string, fileHash string) {
	rw := bufio.NewWriter(c.conn)

	err := eofutil.WriteServerCheckEOF(rw, fmt.Sprintf(UploadPattern, fileName, fileHash), c.stopChan)
	if err != nil {
		log.Printf("an error while writing an upload command to server: %v\n", err)
	}
}

func (c ServerCommunicator) StopUploading(fileName string) {
	rw := bufio.NewWriter(c.conn)

	err := eofutil.WriteServerCheckEOF(rw, fmt.Sprintf(StopUploadPattern, fileName), c.stopChan)
	if err != nil {
		log.Printf("an error while writing a stop-upload command to server: %v\n", err)
	}
}

func (c ServerCommunicator) Disconnect() {
	rw := bufio.NewWriter(c.conn)

	err := eofutil.WriteServerCheckEOF(rw, DisconnectRequest, c.stopChan)
	if err != nil {
		log.Printf("an error while disconnecting from server: %v\n", err)
	}
}

func (c ServerCommunicator) printServerMessage(msg string) {
	fmt.Println(ServerMessagesColour)
	fmt.Println(msg)
	fmt.Print(resetColour)
}

func (c ServerCommunicator) concatUsernamePort(username string, port uint) string {
	portString := strconv.Itoa(int(port))
	return username + "#" + portString + "\n"
}
