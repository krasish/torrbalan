package connection

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const GetOwnersPattern = "GET_OWNERS %s\n"
const UploadPattern = "UPLOAD %s\n"

type ServerCommunicator struct {
	conn     net.Conn
	stopChan chan struct{}
}

func (c ServerCommunicator) GetOwners(filename string) {
	rw := bufio.NewWriter(c.conn)

	_, err := rw.WriteString(fmt.Sprintf(GetOwnersPattern, filename))
	if err != nil {
		fmt.Printf("While writing to server: %v\n", err)
	}
}

func (c ServerCommunicator) StartUploading(filename string) {
	rw := bufio.NewWriter(c.conn)

	_, err := rw.WriteString(fmt.Sprintf(UploadPattern, filename))
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
				close(c.stopChan)
			}
			log.Printf("while reading form server: %v", err)
		}
		fmt.Println(readString)
	}
}
