package response

import (
	"errors"
	"fmt"
	"time"
)

func (r *Response) WriteChunk(data []byte) (err error) {
	defer func() {
		if err != nil {
			r.writeErr = err
		}
	}()

	if !r.chunked {
		return errors.New("chunked encoding is not enabled")
	}

	if err = r.checkCancel(); err != nil {
		return err
	}

	if !r.headerWritten {
		// write headers before first chunk
		if err = r.writeHeaders(); err != nil {
			return err
		}
		r.headerWritten = true
	}

	// write deadline before writing
	r.Conn.SetWriteDeadline(time.Now().Add(r.Cfg.WriteTimeout))

	// chunk size in hex
	size := fmt.Sprintf("%x\r\n", len(data))

	if err = safeWriteString(r.Conn, size); err != nil {
		return err
	}

	if _, err = safeWrite(r.Conn, data); err != nil {
		return err
	}

	if err = safeWriteString(r.Conn, "\r\n"); err != nil {
		return err
	}

	return nil
}

func (r *Response) EndChunked() (err error) {
	defer func() {
		if err != nil {
			r.writeErr = err
		}
	}()

	if !r.chunked {
		return errors.New("not chunked response")
	}

	if err = r.checkCancel(); err != nil {
		return err
	}

	// write deadline before writing
	r.Conn.SetWriteDeadline(time.Now().Add(r.Cfg.WriteTimeout))

	return safeWriteString(r.Conn, "0\r\n\r\n")
}
