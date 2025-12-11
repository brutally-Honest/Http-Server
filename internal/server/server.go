package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/brutally-Honest/http-server/internal/config"
	"github.com/brutally-Honest/http-server/internal/router"
)

type Server struct {
	Addr     string
	listener net.Listener
	running  bool
	mu       sync.Mutex
	config   *config.Config
	matcher  router.RouteMatcher
}

func NewServer(Addr string, config *config.Config, router router.RouteMatcher) *Server {
	return &Server{
		Addr:    Addr,
		config:  config,
		matcher: router,
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
