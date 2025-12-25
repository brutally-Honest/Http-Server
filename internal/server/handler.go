package server

import (
	"bufio"
	"context"
	"log"
	"net"
	"strings"

	"github.com/brutally-Honest/http-server/internal/request"
	"github.com/brutally-Honest/http-server/internal/response"
)

func handleRequest(conn net.Conn, reader *bufio.Reader, s *Server, ctx context.Context) bool {

	req, reqErr := request.ParseRequest(conn, reader, s.config)
	if reqErr != nil {
		log.Println("parse error: ", reqErr.Error())
		res := response.NewResponseWithContext(400, ctx, nil, conn, s.config)
		res.Write([]byte("Bad Request"))
		res.Flush(nil, true)
		return true
	}
	reqCtx, cancelReq := context.WithCancel(ctx)
	defer cancelReq()

	handler, params, err := s.matcher.Match(req.Method, req.Path)
	if err != nil {
		log.Println("router error: ", err.Error())
		res := response.NewResponseWithContext(404, ctx, reqCtx, conn, s.config)
		res.Write([]byte("Not Found"))
		res.Flush(req, false)
		return strings.ToLower(req.Headers["Connection"]) == "close"
	}

	req.Params = params
	req.Context = reqCtx

	res := response.NewResponseWithContext(200, ctx, reqCtx, conn, s.config)
	handler(req, res)

	if res.HasError() {
		return true // write errors
	}

	return strings.ToLower(req.Headers["Connection"]) == "close"
}
