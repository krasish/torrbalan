package command

import (
	"github.com/krasish/torrbalan/server/internal/memory"
	"net"
)

type StopUploadCommand struct {
	conn     net.Conn
	user     memory.User
	fm 		 *memory.FileManager
	fileName string
}

func NewStopUploadCommand(conn net.Conn, user memory.User, fm *memory.FileManager, filename string) *StopUploadCommand {
	return &StopUploadCommand{conn: conn, user: user, fm: fm, fileName: filename}
}

func (c *StopUploadCommand) Do() error {
	if err := c.fm.DeleteUserFromFileInfo(c.fileName, c.user); err != nil {
		return err
	}
	return nil
}
