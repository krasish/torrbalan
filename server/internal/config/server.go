package config

import "errors"

//Server is a struct wrapping all information needed
//initially for a server to be started.
type Server struct {
	Port             string
	ConcurrencyLimit uint
}

//NewServer validates the configuration and returns a *config.Server or
//an error if some of the arguments are unacceptable.
func NewServer(port string, limit uint) (*Server, error) {
	if limit < 1 {
		return nil, errors.New("attempted to start upload with zero limit")
	}
	return &Server{
		Port:             port,
		ConcurrencyLimit: limit,
	}, nil
}
