package request

import (
	"bufio"
	"bytes"
	"errors"
	"io"

	"github.com/brutally-Honest/http-server/internal/config"
)

func readHeaders(reader *bufio.Reader, cfg *config.Config) ([]byte, error) {
	delimiter := []byte("\r\n\r\n")
	headers := make([]byte, 0, cfg.HeaderLimit)
	for {
		line, err := reader.ReadSlice('\n')
		if err != nil {
			if err == bufio.ErrBufferFull {
				// Line too long
				return nil, errors.New("header line too long")
			}
			if err == io.EOF {
				return nil, ErrConnectionClosed
			}
			return nil, err
		}
		if len(headers)+len(line) > cfg.HeaderLimit {
			return nil, ErrHeaderLimitExceeded
		}

		headers = append(headers, line...)
		if bytes.HasSuffix(headers, delimiter) {
			// Remove the \r\n\r\n delimiter from the end
			return headers[:len(headers)-4], nil
		}
	}
}

func readBody(reader *bufio.Reader, cfg *config.Config, contentLength int) ([]byte, error) {
	if contentLength == 0 {
		return nil, nil
	}

	if contentLength > cfg.BodyLimit {
		return nil, ErrBodyLimitExceeded
	}

	body := make([]byte, contentLength)

	_, err := io.ReadFull(reader, body)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, ErrConnectionClosed
		}
		return nil, err
	}

	return body, nil
}
