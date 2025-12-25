package response

import (
	"errors"
	"fmt"
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
		r.setWriteDeadlineOnce()
		// write headers before first chunk
		if err = r.writeHeaders(); err != nil {
			return err
		}
		r.headerWritten = true
	}

	if _, err = fmt.Fprintf(r.writer, "%x\r\n", len(data)); err != nil {
		return err
	}

	if _, err = r.writer.Write(data); err != nil {
		return err
	}

	if _, err = r.writer.WriteString("\r\n"); err != nil {
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
	r.setWriteDeadlineOnce()
	if _, err = r.writer.WriteString("0\r\n\r\n"); err != nil {
		return err
	}
	return r.writer.Flush()
}

// gives users control over when to send data (like net/http's Flusher)
func (r *Response) FlushChunk() (err error) {
	if !r.chunked {
		return errors.New("not a chunked response")
	}

	if err := r.checkCancel(); err != nil {
		return err
	}

	r.setWriteDeadlineOnce()

	return r.writer.Flush()
}
