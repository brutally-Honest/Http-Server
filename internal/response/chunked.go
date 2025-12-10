package response

import (
	"errors"
	"fmt"
	"net"
)

func (r *Response) WriteChunk(conn net.Conn, data []byte) error {
	if !r.chunked {
		return errors.New("chunked encoding is not enabled")
	}

	if err := r.checkCancel(); err != nil {
		return err
	}

	if !r.headerWritten {
		// write headers before first chunk
		if err := r.writeHeaders(conn); err != nil {
			return err
		}
		r.headerWritten = true
	}

	// chunk size in hex
	size := fmt.Sprintf("%x\r\n", len(data))

	if err := safeWriteString(conn, size); err != nil {
		return err
	}

	if _, err := safeWrite(conn, data); err != nil {
		return err
	}

	if err := safeWriteString(conn, "\r\n"); err != nil {
		return err
	}

	return nil
}

func (r *Response) EndChunked(conn net.Conn) error {
	if !r.chunked {
		return errors.New("not chunked response")
	}

	if err := r.checkCancel(); err != nil {
		return err
	}

	return safeWriteString(conn, "0\r\n\r\n")
}
