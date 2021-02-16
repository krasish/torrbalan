package server

import (
	"bufio"
	"fmt"
	"log"
	"net"

	"github.com/krasish/torrbalan/client/pkg/eofutil"

	"github.com/krasish/torrbalan/server/internal/command"
	"github.com/krasish/torrbalan/server/internal/config"
	"github.com/krasish/torrbalan/server/internal/memory"
)

type Server struct {
	*memory.UserManager
	*memory.FileManager
	limiter chan struct{}
	port    string
}

func NewServer(config *config.Server) *Server {
	return &Server{
		UserManager: memory.NewEmptyUserManager(),
		FileManager: memory.NewEmptyFileManager(),
		limiter:     make(chan struct{}, config.ConcurrencyLimit),
		port:        config.Port,
	}
}

//Run starts listening on the configured port and loops accepting clients.
//A new goroutine which serves clients is started for each client.
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("while getting a listener: %w", err)
	}
	log.Printf("Started torrbalan server listeling on %s\n", listener.Addr().String())

	go s.FileManager.SyncFiles()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("An error ocurred while acceppting a connection: %v\n", err)
		}

		s.limiter <- struct{}{}
		go s.handleConnection(conn)
	}
}

//handleConnection registers client and creates a command.Parser associated
//with the client. The function then loops serving client requests
//using that command.Parser.
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
		s.parseCommandAndExecute(parser, remoteAddr)
		if parser.ConnectionClosed {
			log.Printf("closed connection with %s", remoteAddr)
			return
		}
	}
}

func (s *Server) parseCommandAndExecute(parser *command.Parser, remoteAddr string) {
	writer := bufio.NewWriter(parser.Conn)
	handler := eofutil.LoggingEOFHandler{DestName: parser.Conn.RemoteAddr().String()}

	doable, err := parser.Parse()
	if err != nil {
		log.Printf("while parsing command of %s: %v\n", remoteAddr, err)
		if err := eofutil.WriteCheckEOF(writer, "your command could not be parsed\n", handler); err != nil {
			log.Printf("could not send message to %s for failed parsing: %v\n", parser.Conn.RemoteAddr().String(), err)
		}
	}
	if parser.ConnectionClosed {
		return
	}
	if err = doable.Do(); err != nil {
		log.Printf("while executing command of %s: %v\n", remoteAddr, err)
	}
}

//closeConnection closes a connection with client and does cleanup on the Server components.
func (s *Server) closeConnection(conn net.Conn, name string, registeredSuccessfully bool) {
	if err := conn.Close(); err != nil {
		log.Printf("could not close connection with %s: %v", name, err)
	}
	defer func() { <-s.limiter }()
	if registeredSuccessfully {
		if err := s.DeleteUser(name); err != nil {
			log.Printf("could not delete user %q: %v", name, err)
			return
		}
		if err := s.RemoveUserFromOwners(name); err != nil {
			log.Printf("could not remove user %q from owners: %v", name, err)
		}
	}
}
