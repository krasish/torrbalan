package download

import (
	"net"
)

type DownloadInfo struct {
	Filename    string
	PeerAddress string
}

type Downloader struct {
	q chan DownloadInfo
}

func NewDownloader(concurrentDownloads uint, conn net.Conn) Downloader {
	return Downloader{
		q: make(chan DownloadInfo, concurrentDownloads),
	}
}

func (d Downloader) Start() {
	for {
		info := <-d.q
		conn := d.connectToPeer(info)
		d.processDownloading(conn)
	}
}

func (d Downloader) Download(info DownloadInfo) {
	d.q <- info
}

func (d Downloader) connectToPeer(info DownloadInfo) net.Conn {
	panic("not implemented")
}

func (d Downloader) processDownloading(conn net.Conn) {
	//TODO: Consider what happens when two simultaneous
	// downloads for the same file are started.
}
