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

	user, err := command.NewRegisterCommand(s.UserManager, conn).Do()
	defer s.closeConnection(conn, user.Name, err == nil)
	if err != nil {
		log.Printf("could not register %s: %v\n", remoteAddr, err)
		return
	}

	parser := command.NewParser(conn, user, s.FileManager, s.UserManager)
	for {
		var msg interface{}
		func() {
			defer func() {
				msg = recover()
			}()
			s.parseCommandAndExecute(parser, remoteAddr)
		}()
		if msg == command.DisconnectMessage {
			log.Printf("closed connection with %s: %v\n", remoteAddr, err)
			return
		}
	}
}

func (s *Server) parseCommandAndExecute(parser *command.Parser, remoteAddr string) {
	doable, err := parser.Parse()
	if err != nil {
		log.Printf("while parsing command of %s: %v\n",remoteAddr, err)
		if _, err := parser.Conn.Write([]byte("your command could not be parsed\n")); err != nil {
			log.Printf("could not send message to client for failed parsing: %v\n", err)
		}
	}

	if err = doable.Do(); err != nil {
		log.Printf("while executing command of %s: %v\n",remoteAddr, err)
		if _, err := parser.Conn.Write([]byte("your command could not be executed\n")); err != nil {
			log.Printf("could not send message to client for failed executing: %v\n", err)
		}
	}
}


func (s *Server) closeConnection(conn net.Conn, name string, registeredSuccessfully bool) {
	<- s.limiter

	if err := conn.Close(); err != nil {
		log.Printf("could not close connection with %s: %v", name, err)
	}
	if registeredSuccessfully {
		if err := s.DeleteUser(name); err != nil {
			log.Printf("could not delete user %q: %v", name, err)
			return
		}
		if err := s.RemoveUserFromOwners(name); err != nil {
			log.Printf("could not rmeove user %q from owners: %v", name, err)
		}
	}
}