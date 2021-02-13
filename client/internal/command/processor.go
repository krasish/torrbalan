package command

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/krasish/torrbalan/client/pkg/eofutil"

	"github.com/krasish/torrbalan/client/internal/domain/upload"

	"github.com/krasish/torrbalan/client/internal/domain/download"

	"github.com/krasish/torrbalan/client/internal/domain/connection"
)

const (
	DisconnectKey = "DC"
	DownloadKey   = "DW"
	GetOwnersKey  = "OWN"
	StopUploadKey = "SUP"
	UploadKey     = "UP"
)

type Processor struct {
	c        *connection.ServerCommunicator
	d        download.Downloader
	u        upload.Uploader
	regexes  map[string]*regexp.Regexp
	stopChan chan struct{}
}

func NewProcessor(c *connection.ServerCommunicator) *Processor {
	return &Processor{c: c, regexes: map[string]*regexp.Regexp{
		DisconnectKey: regexp.MustCompile(`^[\s]*exit[\s]*$`),
		DownloadKey:   regexp.MustCompile(`^[\s]*download[\s]+([0-9A-Za-z.\-_+$]+)[\s]+([^\s]+)[\s]*$`),
		GetOwnersKey:  regexp.MustCompile(`^[\s]*get-owners[\s]+([0-9A-Za-z.\-_+$]+)[\s]*$`),
		StopUploadKey: regexp.MustCompile(`^[\s]*stop-upload[\s]+([0-9A-Za-z.\-_+$]+)[\s]*$`),
		UploadKey:     regexp.MustCompile(`^[\s]*upload[\s]+([^\0\s]+)[\s]*$`),
	}}
}

func (p Processor) Register() {
	for {
		username := p.getUsername()
		if err := p.c.Register(username); err != nil {
			fmt.Printf("Unsuccessful registration: %v", err)
			continue
		}
		break
	}
}

func (p Processor) Process() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		if p.regexes[DisconnectKey].MatchString(cmd) {
			p.disconnect()
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
	}
}

func (p Processor) getUsername() (username string) {
	var (
		err    error
		reader = bufio.NewReader(os.Stdin)
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

func (p Processor) disconnect() {
	p.c.Disconnect()
	eofutil.TryWrite(p.stopChan)
}

func (p Processor) download(cmd string) {
	captureGroups := p.regexes[DownloadKey].FindStringSubmatch(cmd)
	info := download.Info{
		Filename:    captureGroups[1],
		PeerAddress: captureGroups[2],
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
