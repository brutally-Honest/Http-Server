package parser

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/brutally-Honest/http-server/internal/config"
	"github.com/brutally-Honest/http-server/internal/models"
)

var (
	ErrHeaderLimitExceeded = errors.New("header size limit exceeded")
	ErrBodyLimitExceeded   = errors.New("body size limit exceeded")
)

func ParseRequest(conn net.Conn, cfg *config.Config) (*models.Request, error) {

	conn.SetReadDeadline(time.Now().Add(cfg.ReadTimeout))
	buffer := make([]byte, cfg.BufferLimit)

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

	contentLength, err := parseContentLength(headersRaw)
	if err != nil {
		log.Printf("parseContentLength error: %v", err)
		return nil, err
	}

	log.Println("Content Length:", contentLength)

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
	log.Printf("Body: %d bytes", len(body))
	log.Printf("Total Request: %d bytes", len(headersRaw)+4+len(body))

	return &models.Request{
		Method:  method,
		Version: version,
		Path:    path,
		Headers: headerMap,
		Body:    body,
	}, nil
}

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

func parseRequestLine(buf []byte) (method, path, version string, err error) {
	idx := bytes.Index(buf, []byte("\r\n"))
	if idx == -1 {
		return "", "", "", errors.New("invalid request line: missing CRLF")
	}

	requestLine := bytes.TrimSpace(buf[:idx])

	// Split by any amount of whitespace
	parts := bytes.Fields(requestLine)
	if len(parts) != 3 {
		return "", "", "", errors.New("invalid request line format")
	}

	method = string(parts[0])
	path = string(parts[1])
	version = string(parts[2])

	if !strings.HasPrefix(version, "HTTP/1.") {
		return "", "", "", fmt.Errorf("unsupported HTTP version: %s", version)
	}

	log.Printf("Method:%s | Path :%s | Version:%s", method, path, version)
	return method, path, version, nil
}

func parseHeaders(headers []byte) (map[string]string, error) {
	headerMap := make(map[string]string)

	// Skip request line
	idx := bytes.Index(headers, []byte("\r\n"))
	if idx == -1 {
		return headerMap, nil
	}

	headerLines := headers[idx+2:] // Skip past first \r\n
	lines := bytes.Split(headerLines, []byte("\r\n"))

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		// Split on first colon
		colonIdx := bytes.IndexByte(line, ':')
		if colonIdx == -1 {
			continue // Skip malformed headers
		}

		key := string(bytes.TrimSpace(line[:colonIdx]))
		value := string(bytes.TrimSpace(line[colonIdx+1:]))
		headerMap[key] = value
	}

	return headerMap, nil
}

func parseContentLength(headers []byte) (int, error) {
	lines := bytes.Split(headers, []byte("\r\n"))
	found := false
	var length int

	for _, line := range lines {
		if len(line) < 15 {
			continue
		}

		// Use EqualFold - case-insensitive without allocation
		if !bytes.EqualFold(line[:15], []byte("content-length:")) {
			continue
		}

		// Multiple Content-Length headers = error
		if found {
			return 0, fmt.Errorf("multiple Content-Length headers")
		}

		value := bytes.TrimSpace(line[15:])

		n, err := strconv.Atoi(string(value))
		if err != nil {
			return 0, fmt.Errorf("invalid Content-Length value: %w", err)
		}

		if n < 0 {
			return 0, fmt.Errorf("negative Content-Length: %d", n)
		}

		length = n
		found = true
	}

	if !found {
		return 0, nil
	}
	return length, nil
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
