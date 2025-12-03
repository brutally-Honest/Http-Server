package request

import (
	"errors"
	"log"
	"net"
	"time"

	"github.com/brutally-Honest/http-server/internal/config"
)

func ParseRequest(conn net.Conn, cfg *config.Config) (*Request, error) {

	conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	buffer := make([]byte, cfg.BufferLimit)
	var contentLength int

	headersRaw, leftover, err := readHeaders(conn, cfg, buffer)
	if err != nil {
		return nil, err
	}

	method, path, version, err := parseRequestLine(headersRaw)
	if err != nil {
		log.Printf("Invalid request line: %v", err)
		return nil, err
	}

	headerMap, err := parseHeaders(headersRaw)
	if err != nil {
		log.Printf("Invalid headers: %v", err)
		return nil, err
	}

	// RFC Requirements
	if version == "HTTP/1.1" {
		if _, ok := headerMap["Host"]; !ok {
			return nil, errors.New("HTTP/1.1 requires Host header")
		}
	}

	_, hasCL := headerMap["content-length"]
	_, hasTE := headerMap["transfer-encoding"]

	if hasCL && hasTE {
		return nil, errors.New("both Content-Length and Transfer-Encoding present")
	}

	//TODO: implement Chunked encoding
	if hasTE {
		return nil, errors.New("Transfer-Encoding: chunked not implemented yet")
	}

	contentLength, err = parseContentLength(headersRaw)
	if err != nil {
		log.Printf("parseContentLength error: %v", err)
		return nil, err
	}

	if contentLength > cfg.BodyLimit {
		return nil, ErrBodyLimitExceeded
	}

	body := make([]byte, len(leftover))
	copy(body, leftover)

	if need := contentLength - len(body); need > 0 {
		more, err := readBody(conn, cfg, need, buffer)
		if err != nil {
			return nil, err
		}
		body = append(body, more...)
	}

	log.Printf("Headers: %d bytes", len(headersRaw)+4)
	log.Printf("Content Length: %d", contentLength)
	log.Printf("Body: %d bytes", len(body))
	log.Printf("Total Request: %d bytes", len(headersRaw)+4+len(body))

	return &Request{
		Method:  method,
		Version: version,
		Path:    path,
		Headers: headerMap,
		Body:    body,
	}, nil
}
