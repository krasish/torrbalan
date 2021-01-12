package download

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

type Registrator struct {
	conn net.Conn
}

func (r Registrator) Register() error {
	rw := bufio.NewReadWriter(bufio.NewReader(r.conn), bufio.NewWriter(r.conn))

	username := r.getUsername()
	if _, err := rw.WriteString(username + "\n"); err != nil {
		return fmt.Errorf("while writing to server: %w", err)
	}

	//TODO: Finish implementation
	_, err := rw.ReadString('\n')
	if err != nil {
		return fmt.Errorf("while rading from server: %w", err)
	}

	return nil
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
