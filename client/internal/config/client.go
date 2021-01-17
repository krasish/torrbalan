package config

import "errors"

type Client struct {
	ServerAddress       string
	ConcurrentDownloads uint
	ConcurrentUploads   uint
}

func NewClient(addr string, concurrentDownloads, concurrentUploads uint) (Client, error) {
	if concurrentDownloads < 1 {
		return Client{}, errors.New("attempted to start client with zero download limit")
	} else if concurrentUploads < 1 {
		return Client{}, errors.New("attempted to start client with zero upload limit")
	}
	return Client{
		ServerAddress:       addr,
		ConcurrentDownloads: concurrentDownloads,
		ConcurrentUploads:   concurrentUploads,
	}, nil
}
