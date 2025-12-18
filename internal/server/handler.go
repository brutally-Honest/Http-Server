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

func handleRequest(conn net.Conn, s *Server, ctx context.Context, cancelConn context.CancelFunc) bool {
	req, reqErr := request.ParseRequest(conn, s.config)
	if reqErr != nil {
		log.Println("connection: ", reqErr.Error())
		cancelConn()
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
		conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
		res.Flush(req, false)

		return strings.ToLower(req.Headers["Connection"]) == "close"
	}
	req.Params = params
	req.Context = reqCtx

	res := response.NewResponseWithContext(200, ctx, reqCtx, conn, s.config)
	handler(req, res)
	// temp route for chunked transfer encoding
	if req.Path == "/stream" && req.Method == "GET" {
		res.SetHeader("Content-Type", "text/plain")
		res.SetHeader("Transfer-Encoding", "chunked")

		chunks := []string{"Testing\n", "Transfer\n", "Encoding\n", "With\n", "HTTP\n", "1.1\n"}
		var streamErr error
		for _, chunk := range chunks {
			// Reset write deadline for each chunk
			conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
			if err := res.WriteChunk([]byte(chunk)); err != nil {
				log.Printf("WriteChunk failed: %v", err)
				streamErr = err
				break
			}
		}

		if streamErr == nil {
			conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
			if err := res.EndChunked(); err != nil {
				log.Printf("EndChunked failed: %v", err)
				streamErr = err
			}
		}

		if streamErr != nil {
			return true
		}

		if strings.ToLower(req.Headers["Connection"]) == "close" {
			return true
		}

		return false
	}
	// TODO: Wrap the response with config for write timeout

	return strings.ToLower(req.Headers["Connection"]) == "close"
}
