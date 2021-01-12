package config

import "errors"

type Server struct {
	Port             string
	ConcurrencyLimit uint
}

func NewServer(port string, limit uint) (*Server, error) {
	if limit < 1 {
		return nil, errors.New("attempted to start uploader with zero limit")
	}
	return &Server{
		Port:             port,
		ConcurrencyLimit: limit,
	}, nil
}
