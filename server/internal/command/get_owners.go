package command

import (
	"fmt"
	"net"

	"github.com/krasish/torrbalan/server/internal/memory"
)

type GetOwnersCommand struct {
	conn     net.Conn
	fm       *memory.FileManager
	filename string
}

func NewGetOwnersCommand(conn net.Conn, fm *memory.FileManager, filename string) *GetOwnersCommand {
	return &GetOwnersCommand{conn: conn, fm: fm, filename: filename}
}

func (c GetOwnersCommand) Do() error {
	fileInfo, err := c.fm.GetFileInfo(c.filename)
	if err != nil {
		if _, err := c.conn.Write([]byte(fmt.Sprintf("Could not find %q", c.filename))); err != nil {
			return fmt.Errorf("while writing error to %s: %w", c.conn.RemoteAddr().String(), err)
		}
		return fmt.Errorf("while getting file %s: %w", c.filename, err)
	}
	holders, err := fileInfo.GetHolders()
	if err != nil {
		if _, err := c.conn.Write([]byte(fmt.Sprintf("Could not fetch info for %q", c.filename))); err != nil {
			return fmt.Errorf("while writing error to %s: %w", c.conn.RemoteAddr().String(), err)
		}
	}
	if _, err := c.conn.Write(holders); err != nil {
		return fmt.Errorf("while writing holders to %s: %w", c.conn.RemoteAddr().String(), err)
	}
	return nil
}