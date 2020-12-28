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
	Do() error
}

type Parser struct {
	conn net.Conn
	user memory.User
	um *memory.UserManager
	fm *memory.FileManager
	*regexSet

}

func NewParser(conn net.Conn, user memory.User, fileManager *memory.FileManager, userManager *memory.UserManager) *Parser {
	return &Parser{
		conn:     conn,
		user:     user,
		um: userManager,
		fm: fileManager,
		regexSet: newRegexSet(),
	}
}

func (p *Parser)Parse() (Doable, error) {
	r := bufio.NewReader(p.conn)
	commandString, err := r.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("while reading from %s: %w", p.conn.RemoteAddr().String() ,err)
	}

	if p.regexSet.upload.MatchString(commandString) {
		return p.uploadCommand(commandString)
	} else if p.regexSet.stopUpload.MatchString(commandString) {

	} else if p.regexSet.download.MatchString(commandString) {

	} else if p.regexSet.disconnect.MatchString(commandString) {

	}
	return NewInvalidCommand(p.conn), nil
}

func (p *Parser)uploadCommand(commandString string) (Doable, error) {
	captureGroups := p.regexSet.upload.FindStringSubmatch(commandString)
	if cgc := len(captureGroups); cgc != uploadCaptureGroupsCount {
		return nil, fmt.Errorf("request matched upload regex but got %d capture groups insted of %d", cgc, uploadCaptureGroupsCount)
	}
	return NewUploadCommand(p.conn, p.user,p.fm, captureGroups[1], captureGroups[2]), nil
}

func (p *Parser) stopUploadCommand(commandString string) (Doable, error) {
	captureGroups := p.regexSet.stopUpload.FindStringSubmatch(commandString)
	if cgc := len(captureGroups); cgc != stopUploadCaptureGroupsCount {
		return nil, fmt.Errorf("request matched upload regex but got %d capture groups insted of %d", cgc, stopUploadCaptureGroupsCount)
	}
	return NewUploadCommand(p.conn, p.user,p.fm, captureGroups[1], captureGroups[2]), nil
}