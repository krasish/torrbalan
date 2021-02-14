package command

import (
	"bufio"
	"fmt"
	"net"

	"github.com/krasish/torrbalan/client/pkg/eofutil"

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
	writer := bufio.NewWriter(c.conn)
	handler := eofutil.LoggingEOFHandler{DestName: c.conn.RemoteAddr().String()}
	fileInfo, err := c.fm.GetFileInfo(c.filename)
	if err != nil {
		if err := eofutil.WriteCheckEOF(writer, fmt.Sprintf("No one has uploaded %q\n", c.filename), handler); err != nil {
			return fmt.Errorf("while writing error to %s: %w", c.conn.RemoteAddr().String(), err)
		}
		return fmt.Errorf("while getting file %s: %w", c.filename, err)
	}
	holders, err := fileInfo.GetHolders()
	if err != nil {
		if err := eofutil.WriteCheckEOF(writer, fmt.Sprintf("Could not fetch info for %q\n", c.filename), handler); err != nil {
			return fmt.Errorf("while writing error to %s: %w", c.conn.RemoteAddr().String(), err)
		}
	}
	if err := eofutil.WriteCheckEOF(writer, string(holders)+"\n", handler); err != nil {
		return fmt.Errorf("while writing holders to %s: %w", c.conn.RemoteAddr().String(), err)
	}
	return nil
}
