package command

import (
	"bufio"
	"fmt"
	"net"

	"github.com/krasish/torrbalan/client/pkg/eofutil"
)

//InvalidCommand represents a command which was not matched by any
//of the regexes of Parser. It writes an appropriate message to the client
//reading on the given net.Conn.
type InvalidCommand struct {
	conn net.Conn
}

func NewInvalidCommand(conn net.Conn) *InvalidCommand {
	return &InvalidCommand{conn: conn}
}

func (c *InvalidCommand) Do() error {
	writer := bufio.NewWriter(c.conn)
	handler := eofutil.LoggingEOFHandler{DestName: c.conn.RemoteAddr().String()}

	if err := eofutil.WriteCheckEOF(writer, "invalid request\n", handler); err != nil {
		return fmt.Errorf("could not write to %s: %w", c.conn.RemoteAddr().String(), err)
	}
	return nil
}
