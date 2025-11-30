package server

import (
	"log"
	"net"
	"time"

	"github.com/brutally-Honest/http-server/internal/request"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	_, req_err := request.ParseRequest(conn, s.config)
	// default response for now
	if req_err != nil {
		resp := "HTTP/1.1 400 Bad Request\r\n" +
			"Content-Length: 11\r\n" +
			"Connection: close\r\n" +
			"\r\n" +
			"Bad Request"

		conn.Write([]byte(resp))
		conn.Close()
		return
	}
	// TODO: Handle the parse errors and write to response
	// TODO: Handle the request based on apt route

	//Send response
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
