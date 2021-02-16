package command

import (
	"bufio"
	"fmt"
	"net"

	"github.com/krasish/torrbalan/client/pkg/eofutil"

	"github.com/krasish/torrbalan/server/internal/memory"
)

//UploadCommand represents the action of a client sending information(name, unique hash)
//about a file which it desires to upload. The command updates the state of the given
//*memory.FileManager to represent that.
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
	writer := bufio.NewWriter(c.conn)
	handler := eofutil.LoggingEOFHandler{DestName: c.conn.RemoteAddr().String()}

	err := c.fm.AddFileInfo(c.fileName, c.fileHash, c.user)
	if err != nil {
		errorMessage := fmt.Sprintf("Could not upload file %q", c.fileName)

		if fileError, ok := err.(memory.FileAlreadyExistsError); ok {
			errorMessage = fileError.Error()
			err = fileError.Wrapped
		}

		if err := eofutil.WriteCheckEOF(writer, errorMessage+"\n", handler); err != nil {
			return fmt.Errorf("while writing error message to getOwners: %w", err)
		}
	}
	return err
}
