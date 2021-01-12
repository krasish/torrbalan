package command

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/krasish/torrbalan/server/internal/memory"
)

const UserAlreadyExists = "UAE\n"
const RegisteredSuccessfully = "RS\n"

type RegisterCommand struct {
	manager *memory.UserManager
	conn    net.Conn
}

func NewRegisterCommand(um *memory.UserManager, conn net.Conn) *RegisterCommand {
	return &RegisterCommand{manager: um, conn: conn}
}

func (rc *RegisterCommand) Do() (memory.User, error) {
	var (
		rw         *bufio.ReadWriter = bufio.NewReadWriter(bufio.NewReader(rc.conn), bufio.NewWriter(rc.conn))
		remoteAddr string            = rc.conn.RemoteAddr().String()
		username   string
		err        error
	)

askForUsername:
	for username, err = rw.ReadString('\n'); err != nil; username, err = rw.ReadString('\n') {
		if err == io.EOF {
			return memory.User{}, fmt.Errorf("while reading username: %w", err)
		}
		log.Printf("could not read username for %s: %v\n", remoteAddr, err)
	}
	username = strings.TrimSuffix(username, "\n")
	user, err := rc.manager.RegisterUser(username, remoteAddr)

	if err != nil { // User already exists. Write error to download and retry process
		_, err = rw.WriteString(UserAlreadyExists)
		if err != nil {
			return user, fmt.Errorf("while writing to %s: %w", remoteAddr, err)
		}
		goto askForUsername
	}
	_, err = rw.WriteString(RegisteredSuccessfully)
	if err != nil {
		return user, fmt.Errorf("while writing to %s: %w", remoteAddr, err)
	}
	return user, nil
}
