package command_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/krasish/torrbalan/client/internal/domain/download"

	"github.com/stretchr/testify/mock"

	"github.com/krasish/torrbalan/client/internal/domain/connection"

	"github.com/krasish/torrbalan/client/internal/command/mocks"

	"github.com/krasish/torrbalan/client/internal/command"
)

type MockArgsResponse struct {
	FunctionName string
	Arguments    []interface{}
	ReturnValues []interface{}
}

func TestProcessor(t *testing.T) {
	var (
		fileName = "test.txt"
		fileHash = "test1234test1234test1234test1234"
		peerAddr = "[::]:12345"
		filePath = "/a/sample/path"
	)

	table := []struct {
		message             string
		readerInput         []byte
		connResponses       []MockArgsResponse
		uploaderResponses   []MockArgsResponse
		downloaderResponses []MockArgsResponse
	}{
		{
			message:     "processes upload command",
			readerInput: []byte(fmt.Sprintf("  upload %s\nexit", fileName)),
			connResponses: []MockArgsResponse{
				{
					FunctionName: "Write",
					Arguments:    []interface{}{[]byte(fmt.Sprintf(connection.UploadPattern, fileName, fileHash))},
					ReturnValues: []interface{}{1, nil},
				},
				{
					FunctionName: "Write",
					Arguments:    []interface{}{[]byte(connection.DisconnectRequest)},
					ReturnValues: []interface{}{1, nil},
				},
			},
			uploaderResponses: []MockArgsResponse{
				{
					FunctionName: "AddFile",
					Arguments:    []interface{}{mock.AnythingOfType("string")},
					ReturnValues: []interface{}{fileName, fileHash, nil},
				},
			},
		},
		{
			message:     "processes stop-upload command",
			readerInput: []byte(fmt.Sprintf("  stop-upload %s\nexit", fileName)),
			connResponses: []MockArgsResponse{
				{
					FunctionName: "Write",
					Arguments:    []interface{}{[]byte(fmt.Sprintf(connection.StopUploadPattern, fileName))},
					ReturnValues: []interface{}{1, nil},
				},
				{
					FunctionName: "Write",
					Arguments:    []interface{}{[]byte(connection.DisconnectRequest)},
					ReturnValues: []interface{}{1, nil},
				},
			},
			uploaderResponses: []MockArgsResponse{
				{
					FunctionName: "RemoveFile",
					Arguments:    []interface{}{fileName},
					ReturnValues: []interface{}{},
				},
			},
		},
		{
			message:     "processes download command",
			readerInput: []byte(fmt.Sprintf("  download %s %s %s\nexit", fileName, filePath, peerAddr)),
			downloaderResponses: []MockArgsResponse{
				{
					FunctionName: "Download",
					Arguments: []interface{}{download.Info{
						Filename:    fileName,
						PeerAddress: peerAddr,
						PathToSave:  filePath,
					}},
					ReturnValues: []interface{}{},
				},
			},
			connResponses: []MockArgsResponse{
				{
					FunctionName: "Write",
					Arguments:    []interface{}{[]byte(connection.DisconnectRequest)},
					ReturnValues: []interface{}{1, nil},
				},
			},
		},
		{
			message:     "processes get-owners command",
			readerInput: []byte(fmt.Sprintf("  get-owners %s\nexit", fileName)),

			connResponses: []MockArgsResponse{
				{
					FunctionName: "Write",
					Arguments:    []interface{}{[]byte(fmt.Sprintf(connection.GetOwnersPattern, fileName))},
					ReturnValues: []interface{}{1, nil},
				},
				{
					FunctionName: "Write",
					Arguments:    []interface{}{[]byte(connection.DisconnectRequest)},
					ReturnValues: []interface{}{1, nil},
				},
			},
		},
	}

	for _, test := range table {
		t.Run(test.message, func(t *testing.T) {
			var (
				dw           = &mocks.DownloaderMock{}
				up           = &mocks.UploaderMock{}
				conn         = &mocks.ConnMock{}
				ch           = make(chan struct{}, 1)
				communicator = connection.NewServerCommunicator(conn, ch)
			)
			for _, mar := range test.connResponses {
				conn.On(mar.FunctionName, mar.Arguments...).Return(mar.ReturnValues...)
			}
			for _, mar := range test.uploaderResponses {
				up.On(mar.FunctionName, mar.Arguments...).Return(mar.ReturnValues...)
			}
			for _, mar := range test.downloaderResponses {
				dw.On(mar.FunctionName, mar.Arguments...).Return(mar.ReturnValues...)
			}
			reader := bytes.NewReader(test.readerInput)
			processor := command.NewProcessor(communicator, dw, up, reader)
			processor.Process()

			up.AssertExpectations(t)
			dw.AssertExpectations(t)
			conn.AssertExpectations(t)
		})
	}
}
