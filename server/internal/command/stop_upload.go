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
	err := c.fm.DeleteUserFromFileInfo(c.fileName, c.user)
	if err != nil {
		errorMessage := fmt.Sprintf("File %q does not exist in db or is corrupted.", c.fileName)

		if ownerError, ok := err.(memory.UserIsNotOwnerError); ok {
			errorMessage = ownerError.Error()
			err = ownerError.Wrapped
		}

		if _, err := c.conn.Write([]byte(errorMessage)); err != nil {
			return fmt.Errorf("while writing error message to download: %w", err)
		}
	}
	return err
}
