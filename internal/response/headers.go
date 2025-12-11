package response

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (r *Response) SetHeader(k, v string) error {
	if r.headerWritten {
		return errors.New("headers already written")
	}

	if strings.ToLower(k) == "content-length" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return errors.New("invalid content-length")
		}
		r.hasContentLength = true
		r.contentLength = n
	}

	if strings.ToLower(k) == "transfer-encoding" && strings.ToLower(v) == "chunked" {
		r.chunked = true
	}

	r.Headers[k] = v
	return nil
}

func (r *Response) WriteHeader(code int) error {
	if r.headerWritten {
		return errors.New("WriteHeader called twice")
	}
	r.StatusCode = code
	r.headerWritten = true
	return nil
}

func (r *Response) writeHeaders() error {
	statusText := getStatusText(r.StatusCode)
	if statusText == "Unknown" {
		return errors.New("invalid status code")
	}
	if _, err := fmt.Fprintf(r.Conn, "HTTP/1.1 %d %s\r\n", r.StatusCode, statusText); err != nil {
		return err
	}

	if !r.chunked && !r.hasContentLength {
		r.Headers["Content-Length"] = strconv.Itoa(len(r.Body))
	}

	if r.chunked {
		r.Headers["Transfer-Encoding"] = "chunked"
	}

	for k, v := range r.Headers {
		header := fmt.Sprintf("%s: %s\r\n", k, v)
		if err := safeWriteString(r.Conn, header); err != nil {
			return err
		}
	}

	if err := safeWriteString(r.Conn, "\r\n"); err != nil {
		return err
	}

	return nil
}
