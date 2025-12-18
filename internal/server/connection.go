package server

import (
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

	for {
		if handleRequest(conn, s, ctx) {
			conn.Close()
			return
		}
	}
}
