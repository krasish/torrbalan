package command

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

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
	c        connection.ServerCommunicator
	d        download.Downloader
	u        upload.Uploader
	regexes  map[string]*regexp.Regexp
	stopChan chan struct{}
}

func NewProcessor(c connection.ServerCommunicator) *Processor {
	return &Processor{c: c, regexes: map[string]*regexp.Regexp{
		DisconnectKey: regexp.MustCompile(`^[\s]*disconnect[\s]*$`),
		DownloadKey:   regexp.MustCompile(`^[\s]*download[\s]+([0-9A-Za-z.\-_+$]+)[\s]+([^\s]+)[\s]*$`),
		GetOwnersKey:  regexp.MustCompile(`^[\s]*get-owners[\s]+([0-9A-Za-z.\-_+$]+)[\s]*$`),
		StopUploadKey: regexp.MustCompile(`^[\s]*stop-upload[\s]+([0-9A-Za-z.\-_+$]+)[\s]*$`),
		UploadKey:     regexp.MustCompile(`^[\s]*upload[\s]+([^\0\s]+)[\s]*$`),
	}}
}

func (p Processor) Process() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		cmd, _ := reader.ReadString('\n')
		if p.regexes[DisconnectKey].MatchString(cmd) {
			close(p.stopChan)
		} else if p.regexes[DownloadKey].MatchString(cmd) {
			p.download(cmd)
		} else if p.regexes[GetOwnersKey].MatchString(cmd) {
			p.getOwners(cmd)
		} else if p.regexes[StopUploadKey].MatchString(cmd) {

		} else if p.regexes[UploadKey].MatchString(cmd) {

		} else {
			fmt.Println("Invalid command!")
		}
	}
}

func (p Processor) download(cmd string) {
	captureGroups := p.regexes[DownloadKey].FindStringSubmatch(cmd)
	info := download.Info{
		Filename:    captureGroups[1],
		PeerAddress: captureGroups[2],
	}
	p.d.Download(info)
}

func (p Processor) upload(cmd string) {
	captureGroups := p.regexes[UploadKey].FindStringSubmatch(cmd)
	hash := p.u.AddFile(captureGroups[1])
	p.c.StartUploading()

}

func (p Processor) getOwners(cmd string) {
	captureGroups := p.regexes[GetOwnersKey].FindStringSubmatch(cmd)
	p.c.GetOwners(captureGroups[1])
}
