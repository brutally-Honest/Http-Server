package server

import (
	"net"
	"time"

	"github.com/brutally-Honest/http-server/internal/request"
	"github.com/brutally-Honest/http-server/internal/response"
)

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	req, req_err := request.ParseRequest(conn, s.config)
	if req_err != nil {
		res := response.NewResponse(400)
		res.Write([]byte("Bad Request"))
		res.Flush(conn, nil, true)
		return
	}
	// TODO: Handle the request based on apt route

	// TODO: Wrap the response with config for write timeout
	res := response.NewResponse(200) // default for now
	conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
	res.Flush(conn, req, false)
}
