package server

import (
	"bufio"
	"context"
	"log"
	"net"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("panic in handleConnection: %v", r)
			conn.Close()
		}
	}()
	ctx, cancelConn := context.WithCancel(context.Background())
	defer cancelConn()

	// per connection buffer
	reader := bufio.NewReaderSize(conn, s.config.ReadBufferSize)
	for {
		if handleRequest(conn, reader, s, ctx) {
			conn.Close()
			return
		}
	}
}
