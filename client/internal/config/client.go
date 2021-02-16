package config

import (
	"errors"
)

//Client is a struct wrapping all information needed
//initially for a client to be started.
type Client struct {
	ServerAddress       string
	ConcurrentDownloads uint
	ConcurrentUploads   uint
	Port                uint
}

//NewClient validates the configuration and returns a config.Client or
//an error if some of the arguments are unacceptable.
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
