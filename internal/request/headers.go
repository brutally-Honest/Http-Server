package request

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

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

		key := strings.ToLower(string(bytes.TrimSpace(line[:colonIdx])))
		if _, exists := headerMap[key]; exists {
			return nil, errors.New("duplicate header: " + key)
		}

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
