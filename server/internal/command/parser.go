//Package command is an implementation of the Command pattern. Implementations of
//Doable interface are structs which encapsulate all information needed to perform an action.
package command

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"

	"github.com/krasish/torrbalan/server/internal/memory"
)

const (
	uploadPattern     = `^[\s]*UPLOAD[\s]+([0-9A-Za-z.\-_\+$]+)[\s]+"(.{32,})"[\s]*$`
	stopUploadPattern = `^[\s]*STOP_UPLOAD[\s]+([0-9A-Za-z.\-_\+$]+)[\s]*$`
	getOwnersPattern  = `^[\s]*GET_OWNERS[\s]+([0-9A-Za-z.\-_\+$]+)[\s]*$`
	disconnectPattern = `^[\s]*DISCONNECT[\s]*$`

	uploadCaptureGroupsCount     = 3
	stopUploadCaptureGroupsCount = 2
	getOwnersCaptureGroupsCount  = 2
)

//Represents a set of *regexp.Regexp used by Parser to match client requests.
type regexSet struct {
	upload     *regexp.Regexp
	stopUpload *regexp.Regexp
	getOwners  *regexp.Regexp
	disconnect *regexp.Regexp
}

func newRegexSet() *regexSet {
	return &regexSet{
		upload:     regexp.MustCompile(uploadPattern),
		stopUpload: regexp.MustCompile(stopUploadPattern),
		getOwners:  regexp.MustCompile(getOwnersPattern),
		disconnect: regexp.MustCompile(disconnectPattern),
	}
}

type Doable interface {
	//Do method executes the command. It should return an error for the upload side
	//and its callers should handle those errors. It should also write errors messages
	//to clients if needed.
	Do() error
}

//Parser is a reads requests from a single client through Conn and creates the respective Doable
//representing the action needed to serve a request. Valid commands are matched by its regexSet.
//On errors which prevent further communication with client, ConnectionClosed is set to true.
type Parser struct {
	Conn             net.Conn
	user             memory.User
	um               *memory.UserManager
	fm               *memory.FileManager
	ConnectionClosed bool
	*regexSet
}

func NewParser(conn net.Conn, user memory.User, fileManager *memory.FileManager, userManager *memory.UserManager) *Parser {
	return &Parser{
		Conn:             conn,
		user:             user,
		um:               userManager,
		fm:               fileManager,
		regexSet:         newRegexSet(),
		ConnectionClosed: false,
	}
}

//Parse reads a client request (a string ending in '\n') from Conn and returns Doable
//with an action representing the request. An error is returned when either reading
//from client is impossible or a Doable cannot be constructed from the client request.
func (p *Parser) Parse() (Doable, error) {
	r := bufio.NewReader(p.Conn)
	commandString, err := r.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			p.ConnectionClosed = true
			return nil, nil
		}
		return nil, fmt.Errorf("while reading from %s: %w", p.Conn.RemoteAddr().String(), err)
	}

	if p.regexSet.upload.MatchString(commandString) {
		return p.uploadCommand(commandString)
	} else if p.regexSet.stopUpload.MatchString(commandString) {
		return p.stopUploadCommand(commandString)
	} else if p.regexSet.getOwners.MatchString(commandString) {
		return p.getOwnersCommand(commandString)
	} else if p.regexSet.disconnect.MatchString(commandString) {
		p.ConnectionClosed = true
		return nil, nil
	}
	return NewInvalidCommand(p.Conn), nil
}

func (p *Parser) uploadCommand(cmd string) (Doable, error) {
	captureGroups := p.regexSet.upload.FindStringSubmatch(cmd)
	if cgc := len(captureGroups); cgc != uploadCaptureGroupsCount {
		return nil, fmt.Errorf("request matched upload regex but got %d capture groups insted of %d", cgc, uploadCaptureGroupsCount)
	}
	return NewUploadCommand(p.Conn, p.user, p.fm, captureGroups[1], captureGroups[2]), nil
}

func (p *Parser) stopUploadCommand(cmd string) (Doable, error) {
	captureGroups := p.regexSet.stopUpload.FindStringSubmatch(cmd)
	if cgc := len(captureGroups); cgc != stopUploadCaptureGroupsCount {
		return nil, fmt.Errorf("request matched stop-upload regex but got %d capture groups insted of %d", cgc, stopUploadCaptureGroupsCount)
	}
	return NewStopUploadCommand(p.Conn, p.user, p.fm, captureGroups[1]), nil
}

func (p *Parser) getOwnersCommand(cmd string) (Doable, error) {
	captureGroups := p.regexSet.getOwners.FindStringSubmatch(cmd)
	if cgc := len(captureGroups); cgc != getOwnersCaptureGroupsCount {
		return nil, fmt.Errorf("request matched getOwners regex but got %d capture groups insted of %d", cgc, getOwnersCaptureGroupsCount)
	}
	return NewGetOwnersCommand(p.Conn, p.fm, captureGroups[1]), nil
}
