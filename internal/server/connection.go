package server

import (
	"log"
	"net"
	"time"

	"github.com/brutally-Honest/http-server/internal/parser"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	parser.ParseRequest(conn, s.config)
	resp := "HTTP/1.1 200 OK\r\n" +
		"Content-Length: 5\r\n" +
		"Content-Type: text/plain\r\n" +
		"\r\n" +
		"Hello"
	respByte := []byte(resp)

	conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))

	_, err := conn.Write(respByte)
	if err != nil {
		log.Printf("Write Error :%v", err)
	}
}
