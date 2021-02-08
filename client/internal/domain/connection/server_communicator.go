package connection

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	GetOwnersPattern       = "GET_OWNERS %s\n"
	UploadPattern          = "UPLOAD %s %s\n"
	UserAlreadyExists      = "UAE"
	RegisteredSuccessfully = "RS"
)

type ServerCommunicator struct {
	conn     net.Conn
	stopChan chan struct{}
}

func (c ServerCommunicator) Register(username string) error {
	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))

	if _, err := rw.WriteString(username + "\n"); err != nil {
		return fmt.Errorf("while writing to server: %w", err)
	}

	resp, err := rw.ReadString('\n')
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

	_, err := rw.WriteString(fmt.Sprintf(GetOwnersPattern, filename))
	if err != nil {
		fmt.Printf("While writing to server: %v\n", err)
	}
}

func (c ServerCommunicator) StartUploading(fileName string, fileHash string) {
	rw := bufio.NewWriter(c.conn)

	_, err := rw.WriteString(fmt.Sprintf(UploadPattern, fileName, fileHash))
	if err != nil {
		fmt.Printf("While writing to server: %v\n", err)
	}
}

func (c ServerCommunicator) Listen() {
	for {
		reader := bufio.NewReader(c.conn)
		readString, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				log.Println("Connection to server dropped")
				//TODO: Process proper server shutdown
				close(c.stopChan)
			}
			log.Printf("while reading form server: %v", err)
		}
		fmt.Println(readString)
	}
}
