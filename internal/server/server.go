package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/brutally-Honest/http-server/internal/config"
)

type Server struct {
	Addr     string
	listener net.Listener
	running  bool
	mu       sync.Mutex
	config   *config.Config
}

func NewServer(Addr string, config *config.Config) *Server {
	return &Server{
		Addr:   Addr,
		config: config,
	}
}

func (s *Server) ListenAndServe() error {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return fmt.Errorf("listening Socket Error : %v", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("connection Error : %v", err)
		}
		go s.handleConnection(conn)
	}
}

// TODO: implement shutdown
