package response

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/brutally-Honest/http-server/internal/request"
)

type Response struct {
	StatusCode    int
	Headers       map[string]string
	Body          []byte
	headerWritten bool
	chunked       bool

	hasContentLength bool
	contentLength    int
}

func NewResponse(code int) *Response {
	return &Response{
		StatusCode: code,
		Headers:    map[string]string{},
	}
}

func (r *Response) SetHeader(k, v string) error {
	if r.headerWritten {
		return errors.New("headers already written")
	}

	if strings.ToLower(k) == "content-length" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("invalid content-length")
		}
		r.hasContentLength = true
		r.contentLength = n
	}

	if strings.ToLower(k) == "transfer-encoding" && strings.ToLower(v) == "chunked" {
		r.chunked = true
	}

	r.Headers[k] = v
	return nil
}

func (r *Response) WriteHeader(code int) error {
	if r.headerWritten {
		return errors.New("WriteHeader called twice")
	}
	r.StatusCode = code
	r.headerWritten = true
	return nil
}

func (r *Response) writeHeaders(conn net.Conn) error {
	statusText := getStatusText(r.StatusCode)
	if statusText == "Unknown" {
		return errors.New("invalid status code")
	}
	if _, err := fmt.Fprintf(conn, "HTTP/1.1 %d %s\r\n", r.StatusCode, statusText); err != nil {
		return err
	}

	if !r.chunked && !r.hasContentLength {
		r.Headers["Content-Length"] = strconv.Itoa(len(r.Body))
	}

	if r.chunked {
		r.Headers["Transfer-Encoding"] = "chunked"
	}

	for k, v := range r.Headers {
		if _, err := fmt.Fprintf(conn, "%s: %s\r\n", k, v); err != nil {
			return err
		}
	}

	if _, err := conn.Write([]byte("\r\n")); err != nil {
		return err
	}

	return nil
}

func (r *Response) Write(b []byte) error {
	if r.chunked {
		return errors.New("cannot use Write() when chunked encoding is enabled; use WriteChunk()")
	}

	if r.hasContentLength && len(r.Body)+len(b) > r.contentLength {
		return errors.New("body exceeds Content-Length")
	}

	if !r.headerWritten {
		r.headerWritten = true
	}

	r.Body = append(r.Body, b...)
	return nil
}

func (r *Response) WriteChunk(conn net.Conn, data []byte) error {
	if !r.chunked {
		return errors.New("chunked encoding is not enabled")
	}

	if !r.headerWritten {
		// write headers before first chunk
		if err := r.writeHeaders(conn); err != nil {
			return err
		}
		r.headerWritten = true
	}

	// chunk size in hex
	size := fmt.Sprintf("%x\r\n", len(data))

	if _, err := conn.Write([]byte(size)); err != nil {
		return err
	}
	if _, err := conn.Write(data); err != nil {
		return err
	}
	if _, err := conn.Write([]byte("\r\n")); err != nil {
		return err
	}

	return nil
}

func (r *Response) EndChunked(conn net.Conn) error {
	if !r.chunked {
		return errors.New("not chunked response")
	}

	_, err := conn.Write([]byte("0\r\n\r\n"))
	return err
}

func (r *Response) Flush(conn net.Conn, req *request.Request, serverWantsClose bool) error {
	if r.chunked {
		return errors.New("use WriteChunk + EndChunked for streaming responses")
	}

	if !r.headerWritten {
		r.headerWritten = true
	}

	connHeader := determineConnectionHeader(req, serverWantsClose)
	r.Headers["Connection"] = connHeader

	if err := r.writeHeaders(conn); err != nil {
		return err
	}

	if r.hasContentLength && len(r.Body) != r.contentLength {
		return errors.New("actual body size does not match Content-Length")
	}

	if _, err := conn.Write(r.Body); err != nil {
		return err
	}
	return nil
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
