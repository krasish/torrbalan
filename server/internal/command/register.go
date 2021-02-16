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

//RegisterCommand is the first command executed after a net.Conn is accepted.
//Its purpose is to read a username from the given net.Conn which will be
//associated with the client reading on that net.Conn.
//
//The command will keep on reading for username and port in format `<username>#<port>\n`
//at which client listens until the client sends an appropriate
//username or an error preventing further reading happens.
//
//On success RegisteredSuccessfully message is send to client.
//If the client sends a username which is already in use a UserAlreadyExists
//message is send.
type RegisterCommand struct {
	manager *memory.UserManager
	conn    net.Conn
}

func NewRegisterCommand(um *memory.UserManager, conn net.Conn) *RegisterCommand {
	return &RegisterCommand{manager: um, conn: conn}
}

func (rc *RegisterCommand) Do() (memory.User, error) {
	var (
		rw           = bufio.NewReadWriter(bufio.NewReader(rc.conn), bufio.NewWriter(rc.conn))
		remoteAddr   = rc.conn.RemoteAddr().String()
		userNamePort string
		err          error
		h            = eofutil.LoggingEOFHandler{DestName: rc.conn.RemoteAddr().String()}
	)

askForUsername:
	for userNamePort, err = rw.ReadString('\n'); err != nil; userNamePort, err = rw.ReadString('\n') {
		if err == io.EOF {
			return memory.User{}, fmt.Errorf("while reading userNamePort: %w", err)
		}
		log.Printf("could not read userNamePort for %s: %v\n", remoteAddr, err)
	}
	username, port := separateUsernamePort(userNamePort)
	remoteAddr = replacePortInAddress(remoteAddr, port)
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
	log.Printf("Client at %s registered succesfully with username %s", remoteAddr, username)
	return user, nil
}

func separateUsernamePort(concated string) (username string, port string) {
	split := strings.Split(concated, "#")
	username = split[0]
	port = strings.TrimSuffix(split[1], "\n")
	return
}

func replacePortInAddress(addr, newPort string) string {
	cleanAddress := addr[:strings.LastIndex(addr, ":")]
	return cleanAddress + ":" + newPort
}
