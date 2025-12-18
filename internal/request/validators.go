package request

import (
	"errors"
	"fmt"
	"strings"
)

func validateMethod(method string) error {
	switch method {
	case "GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS", "PATCH":
		return nil
	default:
		return fmt.Errorf("invalid http method: %q", method)
	}
}

func validatePath(path string) error {
	if len(path) > 8192 {
		return errors.New("path too long")
	}
	if strings.Contains(path, "\x00") {
		return errors.New("null byte in path")
	}
	if !strings.HasPrefix(path, "/") && path != "*" {
		return errors.New("invalid path format")
	}
	return nil
}

func validateVersion(version string) error {
	switch version {
	case "HTTP/1.0", "HTTP/1.1":
		return nil
	default:
		return errors.New("unsupported HTTP version: only HTTP/1.0 and HTTP/1.1 supported")
	}
}
