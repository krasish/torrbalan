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

func NewServerCommunicator(conn net.Conn, stopChan chan struct{}) *ServerCommunicator {
	return &ServerCommunicator{conn: conn, stopChan: stopChan}
}

func (c ServerCommunicator) Register(username string) error {
	rw := bufio.NewReadWriter(bufio.NewReader(c.conn), bufio.NewWriter(c.conn))

	if err := WriteCheckEOF(rw.Writer, username+"\n", c.stopChan); err != nil {
		return fmt.Errorf("while writing to server: %w", err)
	}

	resp, err := ReadCheckEOF(rw.Reader, '\n', c.stopChan)
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

	err := WriteCheckEOF(rw, fmt.Sprintf(GetOwnersPattern, filename), c.stopChan)
	if err != nil {
		fmt.Printf("While writing to server: %v\n", err)
	}
}

func (c ServerCommunicator) StartUploading(fileName string, fileHash string) {
	rw := bufio.NewWriter(c.conn)

	err := WriteCheckEOF(rw, fmt.Sprintf(UploadPattern, fileName, fileHash), c.stopChan)
	if err != nil {
		fmt.Printf("While writing to server: %v\n", err)
	}
}

func (c ServerCommunicator) Listen() {
	for {
		reader := bufio.NewReader(c.conn)
		readString, err := ReadCheckEOF(reader, '\n', c.stopChan)
		if err != nil {
			log.Printf("while reading form server: %v", err)
		}
		fmt.Println(readString)
	}
}

func WriteCheckEOF(writer *bufio.Writer, s string, stopChan chan<- struct{}) error {
	if _, err := writer.WriteString(s); err != nil {
		if err == io.EOF {
			log.Println("EOF while writing to server.")
			stopChan <- struct{}{}
		} else {
			return err
		}
	}
	return nil
}

func ReadCheckEOF(reader *bufio.Reader, delim byte, stopChan chan<- struct{}) (string, error) {
	read, err := reader.ReadString(delim)
	if err != nil {
		if err == io.EOF {
			log.Println("EOF while reading from server.")
			stopChan <- struct{}{}
			return "", nil
		}
		return "", err
	}
	return read, nil
}
