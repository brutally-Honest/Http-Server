package parser

import (
	"errors"
	"strings"
)

func validateMethod(method string) (bool, error) {
	switch method {
	case "GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH":
		return true, nil
	default:
		return false, errors.New("invalid http method")
	}
}

func validatePath(path string) (bool, error) {
	if len(path) > 8192 {
		return false, errors.New("path too long")
	}
	if strings.Contains(path, "\x00") {
		return false, errors.New("null byte in path")
	}
	if !strings.HasPrefix(path, "/") && path != "*" {
		return false, errors.New("invalid path format")
	}
	return true, nil
}

func validateVersion(version string) (bool, error) {
	switch version {
	case "HTTP/1.0", "HTTP/1.1":
		return true, nil
	default:
		return false, errors.New("unsupported HTTP version: only HTTP/1.0 and HTTP/1.1 supported")
	}
}
