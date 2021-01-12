package config

import "errors"

type Downloader struct {
	ServerAddress       string
	ConcurrentDownloads uint
}

func NewDownloader(addr string, concurrency uint) (Downloader, error) {
	if concurrency < 1 {
		return Downloader{}, errors.New("attempted to start download with zero limit")
	}
	return Downloader{
		ServerAddress:       addr,
		ConcurrentDownloads: concurrency,
	}, nil
}
