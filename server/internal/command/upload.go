package command

import (
	"fmt"
	"net"

	"github.com/krasish/torrbalan/server/internal/memory"
)

type UploadCommand struct {
	conn     net.Conn
	user     memory.User
	fm       *memory.FileManager
	fileName string
	fileHash string
}

func NewUploadCommand(conn net.Conn, user memory.User, fm *memory.FileManager, filename string, fileHash string) *UploadCommand {
	return &UploadCommand{conn: conn, user: user, fm: fm, fileName: filename, fileHash: fileHash}
}

func (c *UploadCommand) Do() error {
	err := c.fm.AddFileInfo(c.fileName, c.fileHash, c.user)
	if err != nil {
		errorMessage := fmt.Sprintf("Could not upload file %q", c.fileName)

		if fileError, ok := err.(memory.FileAlreadyExistsError); ok {
			errorMessage = fileError.Error()
			err = fileError.Wrapped
		}

		if _, err := c.conn.Write([]byte(errorMessage)); err != nil {
			return fmt.Errorf("while writing error message to download: %w", err)
		}
	}
	return err
}
