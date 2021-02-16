package command

import (
	"bufio"
	"fmt"
	"net"

	"github.com/krasish/torrbalan/client/pkg/eofutil"

	"github.com/krasish/torrbalan/server/internal/memory"
)

//StopUploadCommand represents the action of a client sending information that it no
//longer uploads a file which it previously desired to upload. The command updates the state of the given
//*memory.FileManager to represent that.
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
	writer := bufio.NewWriter(c.conn)
	handler := eofutil.LoggingEOFHandler{DestName: c.conn.RemoteAddr().String()}

	err := c.fm.DeleteUserFromFileInfo(c.fileName, c.user)
	if err != nil {
		errorMessage := fmt.Sprintf("File %q does not exist in db or is corrupted.", c.fileName)

		if ownerError, ok := err.(memory.UserIsNotOwnerError); ok {
			errorMessage = ownerError.Error()
			err = ownerError.Wrapped
		}

		if err := eofutil.WriteCheckEOF(writer, errorMessage+"\n", handler); err != nil {
			return fmt.Errorf("while writing error message to client: %w", err)
		}
	}
	return err
}
