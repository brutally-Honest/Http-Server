package response

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"

	"github.com/brutally-Honest/http-server/internal/config"
	"github.com/brutally-Honest/http-server/internal/request"
)

type Response struct {
	StatusCode    int
	Headers       map[string]string
	Body          []byte
	headerWritten bool
	chunked       bool
	Conn          net.Conn
	Cfg           *config.Config

	hasContentLength bool
	contentLength    int

	connCtx context.Context
	reqCtx  context.Context
}

func NewResponseWithContext(code int, connCtx, reqCtx context.Context, conn net.Conn, cfg *config.Config) *Response {
	return &Response{
		StatusCode: code,
		Headers:    map[string]string{},
		connCtx:    connCtx,
		reqCtx:     reqCtx,
		Conn:       conn,
		Cfg:        cfg,
	}
}

func isConnectionError(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, io.EOF) {
		return true
	}
	// net.Error but not timeout → client closed/reset
	if ne, ok := err.(net.Error); ok && !ne.Timeout() {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "broken pipe") ||
		strings.Contains(msg, "connection reset") ||
		strings.Contains(msg, "use of closed network connection")
}

func determineConnectionHeader(req *request.Request, serverWantsClose bool) string {
	if serverWantsClose || req == nil {
		return "close"
	}

	if connHeader, exists := req.Headers["connection"]; exists {
		if strings.ToLower(connHeader) == "close" {
			return "close"
		}
	}
	// HTTP/1.0 default
	if req.Version == "HTTP/1.0" {
		return "close"
	}

	// HTTP/1.1 default
	return "keep-alive"
}

func getStatusText(code int) string {
	statusTexts := map[int]string{
		200: "OK",
		201: "Created",
		204: "No Content",

		400: "Bad Request",
		404: "Not Found",
		405: "Method Not Allowed",
		408: "Request Timeout",

		500: "Internal Server Error",
		503: "Service Unavailable",
	}

	if text, exists := statusTexts[code]; exists {
		return text
	}
	return "Unknown"
}

func (r *Response) checkCancel() error {
	if r.reqCtx != nil {
		select {
		case <-r.reqCtx.Done():
			return r.reqCtx.Err()
		default:
		}
	}

	select {
	case <-r.connCtx.Done():
		return context.Canceled
	default:
	}

	return nil
}
