package request

import (
	"bufio"
	"errors"
	"log"
	"net"
	"time"

	"github.com/brutally-Honest/http-server/internal/config"
)

func ParseRequest(conn net.Conn, reader *bufio.Reader, cfg *config.Config) (*Request, error) {

	conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	var contentLength int

	headersRaw, err := readHeaders(reader, cfg)
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
		if _, ok := headerMap["host"]; !ok {
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

	contentLength, err = getContentLength(headerMap)
	if err != nil {
		log.Printf("parseContentLength error: %v", err)
		return nil, err
	}

	if contentLength > cfg.BodyLimit {
		return nil, ErrBodyLimitExceeded
	}

	body, err := readBody(reader, cfg, contentLength)
	if err != nil {
		return nil, err
	}

	log.Printf("Headers: %d bytes", len(headersRaw)+4)
	log.Printf("Content Length: %d", contentLength)
	log.Printf("Total Request: %d bytes", len(headersRaw)+4+len(body))

	return &Request{
		Method:  method,
		Version: version,
		Path:    path,
		Headers: headerMap,
		Body:    body,
	}, nil
}
