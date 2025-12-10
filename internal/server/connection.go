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
			log.Println("connection: ", reqErr.Error())
			cancelConn()
			res := response.NewResponseWithContext(400, ctx, nil, conn, s.config)
			res.Write([]byte("Bad Request"))
			res.Flush(conn, nil, true)
			conn.Close()
			return
		}
		reqCtx, cancelReq := context.WithCancel(ctx)
		// TODO: Handle the request based on apt route with reqCtx passed

		handler, params, err := s.matcher.Match(req.Path)
		if err != nil {
			log.Println("router error: ", err.Error())
			res := response.NewResponseWithContext(404, ctx, reqCtx, conn, s.config)
			res.Write([]byte("Not Found"))
			conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
			res.Flush(conn, req, false)
			cancelReq()

			if strings.ToLower(req.Headers["Connection"]) == "close" {
				cancelConn()
				conn.Close()
				return
			}
			continue
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
				if err := res.WriteChunk(conn, []byte(chunk)); err != nil {
					log.Printf("WriteChunk failed: %v", err)
					streamErr = err
					break
				}
			}

			if streamErr == nil {
				conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
				if err := res.EndChunked(conn); err != nil {
					log.Printf("EndChunked failed: %v", err)
					streamErr = err
				}
			}

			cancelReq()

			if streamErr != nil {
				conn.Close()
				return
			}

			if strings.ToLower(req.Headers["Connection"]) == "close" {
				cancelConn()
				conn.Close()
				return
			}

			continue
		}
		// TODO: Wrap the response with config for write timeout
		// conn.SetWriteDeadline(time.Now().Add(s.config.WriteTimeout))
		// if err := res.Flush(conn, req, false); err != nil {
		// 	cancelReq()
		// 	conn.Close()
		// 	return
		// }
		cancelReq()
		if strings.ToLower(req.Headers["Connection"]) == "close" {
			cancelConn()
			conn.Close()
			return
		}
	}
}
