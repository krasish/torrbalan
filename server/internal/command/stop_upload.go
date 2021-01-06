package command

import (
	"fmt"
	"net"

	"github.com/krasish/torrbalan/server/internal/memory"
)

type StopUploadCommand struct {
	conn     net.Conn
	user     memory.User
	fm       *memory.FileManager
	fileName string
}

func NewStopUploadCommand(conn net.Conn, user memory.User, fm *memory.FileManager, filename string) *StopUploadCommand {
	return &StopUploadCommand{conn: conn, user: user, fm: fm, fileName: filename}
}

func (c *StopUploadCommand) Do() error {
	if err := c.fm.DeleteUserFromFileInfo(c.fileName, c.user); err != nil {
		//TODO: Create branching for different errors
		if _, err := c.conn.Write([]byte("You are not uploading this file.")); err != nil {
			return fmt.Errorf("while writing error message to client: %w", err)
		}
		return err
	}
	return nil
}
