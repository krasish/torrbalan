package config

import (
	"errors"
)

type Client struct {
	ServerAddress       string
	ConcurrentDownloads uint
	ConcurrentUploads   uint
	Port                uint
}

func NewClient(addr string, concurrentDownloads, concurrentUploads, port uint) (Client, error) {
	if concurrentDownloads < 1 {
		return Client{}, errors.New("attempted to start client with zero download limit")
	} else if concurrentUploads < 1 {
		return Client{}, errors.New("attempted to start client with zero upload limit")
	} else if port < 1024 {
		return Client{}, errors.New("attempted to start client at well-known port")
	}
	return Client{
		ServerAddress:       addr,
		ConcurrentDownloads: concurrentDownloads,
		ConcurrentUploads:   concurrentUploads,
		Port:                port,
	}, nil
}
