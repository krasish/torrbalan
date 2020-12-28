package command

import (
	"fmt"
	"net"
)

type InvalidCommand struct {
	conn    net.Conn
}

func NewInvalidCommand(conn net.Conn) *InvalidCommand {
	return &InvalidCommand{conn: conn}
}

func (c *InvalidCommand) Do() error {
	if _, err := c.conn.Write([]byte("Invalid request")); err != nil {
		return fmt.Errorf("could not write to %s: %w", c.conn.RemoteAddr().String(), err)
	}
	return nil
}
