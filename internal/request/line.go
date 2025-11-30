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

	if validMethod, err := validateMethod(method); !validMethod {
		return "", "", "", err
	}

	if validPath, err := validatePath(path); !validPath {
		return "", "", "", err
	}

	if validVersion, err := validateVersion(version); !validVersion {
		return "", "", "", err
	}

	log.Printf("Method:%s | Path :%s | Version:%s", method, path, version)
	return method, path, version, nil
}
