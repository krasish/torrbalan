package command

import (
	"bufio"
	"fmt"
	"io"
	"path"
	"regexp"
	"strings"

	"github.com/krasish/torrbalan/client/pkg/eofutil"

	"github.com/krasish/torrbalan/client/internal/domain/download"

	"github.com/krasish/torrbalan/client/internal/domain/connection"
)

const (
	ExitKey       = "DC"
	DownloadKey   = "DW"
	GetOwnersKey  = "OWN"
	StopUploadKey = "SUP"
	UploadKey     = "UP"
)

type Downloader interface {
	Download(info download.Info)
}

type Uploader interface {
	AddFile(filePath string) (name string, hash string, err error)
	RemoveFile(fileName string)
}

//Processor reads commands from r, matches them using regexes and takes the
//respective action to serve the command.
type Processor struct {
	c        *connection.ServerCommunicator
	d        Downloader
	u        Uploader
	regexes  map[string]*regexp.Regexp
	stopChan chan struct{}
	r        io.Reader
}

func NewProcessor(c *connection.ServerCommunicator, d Downloader, u Uploader, r io.Reader) *Processor {
	return &Processor{
		c: c,
		d: d,
		u: u,
		regexes: map[string]*regexp.Regexp{
			ExitKey:       regexp.MustCompile(`^[\s]*exit[\s]*$`),
			DownloadKey:   regexp.MustCompile(`^[\s]*download[\s]+([0-9A-Za-z.\-_+$]+)[\s]+([^\s]+)[\s]+([^\s]+)[\s]*$`),
			GetOwnersKey:  regexp.MustCompile(`^[\s]*get-owners[\s]+([0-9A-Za-z.\-_+$]+)[\s]*$`),
			StopUploadKey: regexp.MustCompile(`^[\s]*stop-upload[\s]+([0-9A-Za-z.\-_+$]+)[\s]*$`),
			UploadKey:     regexp.MustCompile(`^[\s]*upload[\s]+([^\0\s]+)[\s]*$`),
		},
		stopChan: nil,
		r:        r,
	}
}

//Register is the first method of Processor that must be called. It asks
//user for username and attempts to register to the server using that username.
//Register loops until a successful registration occurs.
func (p Processor) Register(port uint) {
	for {
		username := p.getUsername()
		if err := p.c.Register(username, port); err != nil {
			fmt.Printf("Unsuccessful registration: %v", err)
			continue
		}
		break
	}
}

//Process loops and reads client commands. It matches commands using
//Processor regexes and takes respective aciton.
func (p Processor) Process() {
	reader := bufio.NewReader(p.r)
	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		if p.regexes[ExitKey].MatchString(cmd) {
			p.exit()
			break
		} else if p.regexes[DownloadKey].MatchString(cmd) {
			p.download(cmd)
		} else if p.regexes[GetOwnersKey].MatchString(cmd) {
			p.getOwners(cmd)
		} else if p.regexes[StopUploadKey].MatchString(cmd) {
			p.stopUpload(cmd)
		} else if p.regexes[UploadKey].MatchString(cmd) {
			p.upload(cmd)
		} else {
			fmt.Println("Invalid command!")
		}
		fmt.Println()
	}
}

func (p Processor) getUsername() (username string) {
	var (
		err    error
		reader = bufio.NewReader(p.r)
	)

	for {
		fmt.Println("Please enter a username: ")
		username, err = reader.ReadString('\n')
		if err == nil {
			username = strings.TrimSuffix(username, "\n")
			break
		}
		fmt.Println("That didn't work. Try again!")
	}

	return
}

func (p Processor) exit() {
	p.c.Disconnect()
	eofutil.TryWrite(p.stopChan)
}

func (p Processor) download(cmd string) {
	captureGroups := p.regexes[DownloadKey].FindStringSubmatch(cmd)
	info := download.Info{
		Filename:    captureGroups[1],
		PathToSave:  captureGroups[2],
		PeerAddress: captureGroups[3],
	}
	p.d.Download(info)
}

func (p Processor) getOwners(cmd string) {
	captureGroups := p.regexes[GetOwnersKey].FindStringSubmatch(cmd)
	p.c.GetOwners(captureGroups[1])
}

func (p Processor) upload(cmd string) {
	captureGroups := p.regexes[UploadKey].FindStringSubmatch(cmd)
	name, hash, err := p.u.AddFile(captureGroups[1])
	if err != nil {
		fmt.Printf("could not start uploading file: %v", err)
		return
	}
	p.c.StartUploading(name, hash)
}

func (p Processor) stopUpload(cmd string) {
	captureGroups := p.regexes[StopUploadKey].FindStringSubmatch(cmd)
	fileName := path.Base(captureGroups[1])
	p.c.StopUploading(fileName)
	p.u.RemoveFile(fileName)
}
