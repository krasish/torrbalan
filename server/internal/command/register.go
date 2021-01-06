package command

import (
	"bufio"
	"fmt"
	"github.com/krasish/torrbalan/server/internal/memory"
	"io"
	"log"
	"net"
	"strings"
)

const UserAlreadyExists = "UAE"
const RegisteredSuccessfully = "RS"

type RegisterCommand struct {
	manager *memory.UserManager
	conn    net.Conn
}

func NewRegisterCommand(um *memory.UserManager, conn net.Conn) *RegisterCommand {
	return &RegisterCommand{manager: um, conn: conn}
}

func (rc *RegisterCommand) Do() (memory.User, error) {
	var (
		r          *bufio.Reader = bufio.NewReader(rc.conn)
		remoteAddr string        = rc.conn.RemoteAddr().String()
		username   string
		err        error
	)

askForUsername:
	for username, err = r.ReadString('\n'); err != nil; username, err = r.ReadString('\n') {
		if err == io.EOF {
			return memory.User{}, fmt.Errorf("while reading username: %w", err)
		}
		log.Printf("could not read username for %s: %v\n", remoteAddr, err)
	}
	username = strings.TrimSuffix(username, "\n")
	user, err := rc.manager.RegisterUser(username, remoteAddr)

	if err != nil { // User already exists. Write error to client and retry process
		_, err = rc.conn.Write([]byte(UserAlreadyExists))
		if err != nil {
			return user, fmt.Errorf("while writing to %s: %w", remoteAddr, err)
		}
		goto askForUsername
	}
	_, err = rc.conn.Write([]byte(RegisteredSuccessfully))
	if err != nil {
		return user, fmt.Errorf("while writing to %s: %w", remoteAddr, err)
	}
	return user, nil
}
