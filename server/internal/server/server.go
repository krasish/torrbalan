package server

import (
	"fmt"
	"github.com/krasish/torrbalan/server/internal/command"
	"github.com/krasish/torrbalan/server/internal/config"
	"github.com/krasish/torrbalan/server/internal/memory"
	"log"
	"net"
)

type Server struct {
	*memory.UserManager
	*memory.FileManager
	limiter chan struct{}
	port string
}

func NewServer(config *config.Server) *Server {
	return &Server{
		UserManager: memory.NewEmptyUserManager(),
		FileManager: memory.NewEmptyFileManager(),
		limiter:     make(chan struct{}, config.ConcurrencyLimit),
		port:        config.Port,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("while getting a listener: %w", err)
	}
	log.Printf("Started listeling on %s\n", listener.Addr().String())

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("An error ocurred while acceppting a connection: %v\n", err)
		}
		s.limiter <- struct{}{}
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	log.Printf("Started handling connection with %s\n", remoteAddr)
	defer s.closeConnection(conn, remoteAddr)

	user, err := command.NewRegisterCommand(s.UserManager, conn).Do()
	if err != nil {
		log.Printf("could not register %q: %v\n", remoteAddr, err)
		return
	}

	parser := command.NewParser(conn, user, s.FileManager, s.UserManager)
	for {
		err := s.parseCommandAndExecute(parser, remoteAddr)
		if err != nil {
			return
		}
	}
}

func (s *Server) parseCommandAndExecute(parser *command.Parser, remoteAddr string) error {
	doable, err := parser.Parse()
	if err != nil {
		//TODO: Tell client you cannot parse command
		log.Printf("could not parse command: %v", err)
		return nil
	}
	//if disconnect, ok := doable.(command.DisconnectCommand); ok {
	//	err := invalid.Do()
	//	return fmt.Errorf("")
	//}

	if err = doable.Do(); err != nil {
		//TODO: Write error to client
		log.Printf("while doing command for %s: %v\n", remoteAddr, err)
	}
	return nil
}


func (s *Server) closeConnection(conn net.Conn, remoteAddr string) {
	<- s.limiter

	if err := conn.Close(); err != nil {
		log.Printf("could not close connection with %s: %v", remoteAddr, err)
	}
	if err := s.DeleteUser(remoteAddr); err != nil {
		log.Printf("could not delete user with addr %s: %v", remoteAddr, err)
		return
	}
	if err := s.RemoveUserFromOwners(remoteAddr); err != nil {
		log.Printf("could not rmeove user with addr %s from owners: %v", remoteAddr, err)
	}
}