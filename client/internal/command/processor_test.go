package command_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/krasish/torrbalan/client/internal/domain/connection"

	"github.com/krasish/torrbalan/client/internal/domain/download"

	"github.com/krasish/torrbalan/client/internal/command/mocks"

	"github.com/krasish/torrbalan/client/internal/command"
)

func TestProcessor(t *testing.T) {
	dw := &mocks.DownloaderMock{}
	up := &mocks.UploaderMock{}
	//TODO: Create a connection mock

	connection.NewServerCommunicator()
	processor := command.NewProcessor(nil, dw, up)
	info := download.Info{
		Filename:    "test.txt",
		PeerAddress: "[::]:12345",
		PathToSave:  "/a/test/path",
	}
	dw.On("Download", info).Return()
	reader := bytes.NewReader([]byte(fmt.Sprintf("download %s %s %s\nexit", info.Filename, info.PathToSave, info.PeerAddress)))
	processor.Process(reader)

	dw.AssertCalled(t, "Download", info.Filename, info.PathToSave, info.PeerAddress)
}
