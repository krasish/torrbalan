package config

type Server struct {
	Port             string
	ConcurrencyLimit int
}

func NewServer(port string, limit int) *Server {
	return &Server{
		Port:             port,
		ConcurrencyLimit: limit,
	}
}
