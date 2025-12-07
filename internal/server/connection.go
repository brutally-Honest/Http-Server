package server

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/brutally-Honest/http-server/internal/request"
	"github.com/brutally-Honest/http-server/internal/response"
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
		req, reqErr := request.ParseRequest(conn, s.config)
		if reqErr != nil {
			log.Println(reqErr.Error())
			cancelConn()
			res := response.NewResponse(400)
			res.Write([]byte("Bad Request"))
			res.Flush(conn, nil, true)
			conn.Close()
			return
		}
		_, cancelReq := context.WithCancel(ctx)
		// TODO: Handle the request based on apt route with reqCtx passed

		// TODO: Wrap the response with config for write timeout
		res := response.NewResponse(200) // default for now
		conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
		res.Flush(conn, req, false)
		cancelReq()
		if strings.ToLower(req.Headers["Connection"]) == "close" {
			cancelConn()
			conn.Close()
			return
		}
	}
}
