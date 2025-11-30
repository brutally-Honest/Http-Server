package request

import (
	"bytes"
	"log"
	"net"

	"github.com/brutally-Honest/http-server/internal/config"
)

// read until \r\n\r\n is found
func readHeaders(conn net.Conn, cfg *config.Config, buffer []byte) ([]byte, []byte, error) {
	headers := make([]byte, 0, cfg.HeaderLimit)
	for {
		streamLength, err := conn.Read(buffer)
		if err != nil {
			log.Printf("read error: %v", err)
			return nil, nil, err
		}

		if len(headers)+streamLength > cfg.HeaderLimit {
			log.Printf("header limit exceeded")
			return nil, nil, ErrHeaderLimitExceeded
		}

		headers = append(headers, buffer[:streamLength]...)

		if idx := bytes.Index(headers, []byte("\r\n\r\n")); idx != -1 {
			headerEnd := idx + 4
			return headers[:idx], headers[headerEnd:], nil
		}
	}
}

// read based on Content-Length
func readBody(conn net.Conn, cfg *config.Config, contentLength int, buffer []byte) ([]byte, error) {
	if contentLength == 0 {
		return nil, nil
	}

	if contentLength > cfg.BodyLimit {
		return nil, ErrBodyLimitExceeded
	}

	body := make([]byte, 0, contentLength)

	for len(body) < contentLength {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("body Read Error :%v", err)
			return nil, err
		}

		if len(body)+n > cfg.BodyLimit {
			return nil, ErrBodyLimitExceeded
		}

		remaining := contentLength - len(body)
		if n > remaining {
			n = remaining
		}

		body = append(body, buffer[:n]...)
	}

	return body, nil
}
