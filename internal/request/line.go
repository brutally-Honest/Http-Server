package request

import (
	"bytes"
	"errors"
	"log"
)

func parseRequestLine(headers []byte) (method, path, version string, err error) {
	idx := bytes.Index(headers, []byte("\r\n"))
	if idx == -1 {
		return "", "", "", errors.New("invalid request line: missing CRLF")
	}

	requestLine := bytes.TrimSpace(headers[:idx])

	// Split by any amount of whitespace
	parts := bytes.Fields(requestLine)
	if len(parts) != 3 {
		return "", "", "", errors.New("invalid request line format")
	}

	method = string(parts[0])
	path = string(parts[1])
	version = string(parts[2])

	if err := validateMethod(method); err != nil {
		return "", "", "", err
	}

	if err := validatePath(path); err != nil {
		return "", "", "", err
	}

	if err := validateVersion(version); err != nil {
		return "", "", "", err
	}

	log.Printf("Method:%s | Path :%s | Version:%s", method, path, version)
	return method, path, version, nil
}
