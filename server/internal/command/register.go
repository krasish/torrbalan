package command

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

	"github.com/krasish/torrbalan/client/pkg/eofutil"

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
		rw         = bufio.NewReadWriter(bufio.NewReader(rc.conn), bufio.NewWriter(rc.conn))
		remoteAddr = rc.conn.RemoteAddr().String()
		username   string
		err        error
		h          = eofutil.LoggingEOFHandler{DestName: rc.conn.RemoteAddr().String()}
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

	if err != nil { // User already exists. Write error to getOwners and retry process
		err := eofutil.WriteCheckEOF(rw.Writer, UserAlreadyExists, h)
		if err != nil {
			return user, fmt.Errorf("while writing to %s: %w", remoteAddr, err)
		}
		goto askForUsername
	}
	err = eofutil.WriteCheckEOF(rw.Writer, RegisteredSuccessfully, h)
	if err != nil {
		return user, fmt.Errorf("while writing to %s: %w", remoteAddr, err)
	}
	log.Printf("Client at %s registered succesfully with username %s", rc.conn.RemoteAddr().String(), username)
	return user, nil
}
