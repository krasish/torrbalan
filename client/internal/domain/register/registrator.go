package register

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

const UserAlreadyExists = "UAE"
const RegisteredSuccessfully = "RS"

type Registrator struct {
	conn net.Conn
}

func NewRegistrator(conn net.Conn) *Registrator {
	return &Registrator{conn: conn}
}

func (r Registrator) Register() error {
	rw := bufio.NewReadWriter(bufio.NewReader(r.conn), bufio.NewWriter(r.conn))

sendUsername:
	username := r.getUsername()
	if _, err := rw.WriteString(username + "\n"); err != nil {
		return fmt.Errorf("while writing to server: %w", err)
	}

	resp, err := rw.ReadString('\n')
	if err != nil {
		return fmt.Errorf("while reading from server: %w", err)
	}

	if resp == UserAlreadyExists {
		fmt.Println("That didn't work. Try again!")
		goto sendUsername
	} else if resp == RegisteredSuccessfully {
		return nil
	} else {
		return fmt.Errorf("could not register. Server responded: %s", resp)
	}
}

func (r Registrator) getUsername() (username string) {
	var (
		err    error
		reader = bufio.NewReader(os.Stdin)
	)

	for {
		fmt.Print("Please enter a username: ")
		username, err = reader.ReadString('\n')
		if err == nil {
			break
		}
		fmt.Println("That didn't work. Try again!")
	}

	return
}
