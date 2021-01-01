package command

import (
	"bufio"
	"fmt"
	"github.com/krasish/torrbalan/server/internal/memory"
	"net"
	"regexp"
)

const (
	uploadPattern = `^[\s]*UPLOAD[\s]+([0-9A-Za-z.\-_\+$]+)[\s]+([A-Fa-f0-9]{64})[\s]*$`
	stopUploadPattern = `^[\s]*STOP_UPLOAD[\s]+([0-9A-Za-z.\-_\+$]+)[\s]*$`
	downloadPattern = `^[\s]*DOWNLOAD[\s]+([0-9A-Za-z.\-_\+$]+)[\s]*$`
	disconnectPattern = `^[\s]*DISCONNECT[\s]*$`

	uploadCaptureGroupsCount = 3
	stopUploadCaptureGroupsCount = 2
	downloadCaptureGroupsCount = 2
)

type regexSet struct {
	upload *regexp.Regexp
	stopUpload *regexp.Regexp
	download *regexp.Regexp
	disconnect *regexp.Regexp
}

func newRegexSet() *regexSet {
	return &regexSet{
		upload:     regexp.MustCompile(uploadPattern),
		stopUpload: regexp.MustCompile(stopUploadPattern),
		download:   regexp.MustCompile(downloadPattern),
		disconnect: regexp.MustCompile(disconnectPattern),
	}
}

type Doable interface {
	//Do method executes the command. It should return an error for the server side
	//and its callers should handle those errors. It should also write errors messages
	//to clients.
	Do() error
}

type Parser struct {
	Conn net.Conn
	user memory.User
	um   *memory.UserManager
	fm   *memory.FileManager
	*regexSet

}

func NewParser(conn net.Conn, user memory.User, fileManager *memory.FileManager, userManager *memory.UserManager) *Parser {
	return &Parser{
		Conn:     conn,
		user:     user,
		um:       userManager,
		fm:       fileManager,
		regexSet: newRegexSet(),
	}
}

func (p *Parser)Parse() (Doable, error) {
	r := bufio.NewReader(p.Conn)
	commandString, err := r.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("while reading from %s: %w", p.Conn.RemoteAddr().String() ,err)
	}

	if p.regexSet.upload.MatchString(commandString) {
		return p.uploadCommand(commandString)
	} else if p.regexSet.stopUpload.MatchString(commandString) {
		return p.stopUploadCommand(commandString)
	} else if p.regexSet.download.MatchString(commandString) {
		return p.downloadCommand(commandString)
	} else if p.regexSet.disconnect.MatchString(commandString) {
		return DisconnectCommand{}, nil
	}
	return NewInvalidCommand(p.Conn), nil
}

func (p *Parser)uploadCommand(commandString string) (Doable, error) {
	captureGroups := p.regexSet.upload.FindStringSubmatch(commandString)
	if cgc := len(captureGroups); cgc != uploadCaptureGroupsCount {
		return nil, fmt.Errorf("request matched upload regex but got %d capture groups insted of %d", cgc, uploadCaptureGroupsCount)
	}
	return NewUploadCommand(p.Conn, p.user,p.fm, captureGroups[1], captureGroups[2]), nil
}

func (p *Parser) stopUploadCommand(commandString string) (Doable, error) {
	captureGroups := p.regexSet.stopUpload.FindStringSubmatch(commandString)
	if cgc := len(captureGroups); cgc != stopUploadCaptureGroupsCount {
		return nil, fmt.Errorf("request matched stop-upload regex but got %d capture groups insted of %d", cgc, stopUploadCaptureGroupsCount)
	}
	return NewStopUploadCommand(p.Conn, p.user,p.fm, captureGroups[1]), nil
}

func (p *Parser) downloadCommand(commandString string) (Doable, error) {
	captureGroups := p.regexSet.download.FindStringSubmatch(commandString)
	if cgc := len(captureGroups); cgc != downloadCaptureGroupsCount {
		return nil, fmt.Errorf("request matched download regex but got %d capture groups insted of %d", cgc, downloadCaptureGroupsCount)
	}
	return NewDownloadCommand(p.Conn, p.fm, captureGroups[1]), nil
}