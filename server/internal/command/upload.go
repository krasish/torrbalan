package command

import (
	"fmt"
	"github.com/krasish/torrbalan/server/internal/memory"
	"net"
)

type UploadCommand struct {
	conn     net.Conn
	user     memory.User
	fm 		 *memory.FileManager
	fileName string
	fileHash string
}

func NewUploadCommand(conn net.Conn, user memory.User, fm *memory.FileManager, filename string, fileHash string) *UploadCommand {
	return &UploadCommand{conn: conn, user: user, fm: fm, fileName: filename, fileHash: fileHash}
}

func (c *UploadCommand) Do() error {
	err := c.fm.AddFileInfo(c.fileName, c.fileHash, c.user)
	if err != nil {
		if _, err := c.conn.Write([]byte(err.Error())); err != nil {
			return fmt.Errorf("whiler writing error message to client: %w", err)
		}
	}
	return nil
}
