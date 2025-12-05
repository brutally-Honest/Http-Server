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
		fmt.Printf("%s:%s\n", key, value)
		headerMap[key] = value
	}

	return headerMap, nil
}

func getContentLength(headers map[string]string) (int, error) {
	val, ok := headers["content-length"]
	if !ok {
		return 0, nil
	}
	n, err := strconv.Atoi(val)
	if err != nil || n < 0 {
		return 0, fmt.Errorf("invalid Content-Length: %s", val)
	}
	return n, nil
}
