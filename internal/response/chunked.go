package response

import (
	"errors"
	"fmt"
)

func (r *Response) WriteChunk(data []byte) error {
	if !r.chunked {
		return errors.New("chunked encoding is not enabled")
	}

	if err := r.checkCancel(); err != nil {
		return err
	}

	if !r.headerWritten {
		// write headers before first chunk
		if err := r.writeHeaders(); err != nil {
			return err
		}
		r.headerWritten = true
	}

	// chunk size in hex
	size := fmt.Sprintf("%x\r\n", len(data))

	if err := safeWriteString(r.Conn, size); err != nil {
		return err
	}

	if _, err := safeWrite(r.Conn, data); err != nil {
		return err
	}

	if err := safeWriteString(r.Conn, "\r\n"); err != nil {
		return err
	}

	return nil
}

func (r *Response) EndChunked() error {
	if !r.chunked {
		return errors.New("not chunked response")
	}

	if err := r.checkCancel(); err != nil {
		return err
	}

	return safeWriteString(r.Conn, "0\r\n\r\n")
}
