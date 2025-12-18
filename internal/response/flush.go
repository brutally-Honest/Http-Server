package response

import (
	"errors"
	"time"

	"github.com/brutally-Honest/http-server/internal/request"
)

func (r *Response) Write(b []byte) (err error) {
	defer func() {
		if err != nil {
			r.writeErr = err
		}
	}()

	if r.chunked {
		return errors.New("cannot use Write() when chunked encoding is enabled; use WriteChunk()")
	}

	if r.hasContentLength && len(r.Body)+len(b) > r.contentLength {
		return errors.New("body exceeds Content-Length")
	}

	if !r.headerWritten {
		r.headerWritten = true
	}

	r.Body = append(r.Body, b...)
	return nil
}

func (r *Response) Flush(req *request.Request, serverWantsClose bool) (err error) {
	defer func() {
		if err != nil {
			r.writeErr = err
		}
	}()

	if r.chunked {
		return errors.New("use WriteChunk + EndChunked for streaming responses")
	}

	if err = r.checkCancel(); err != nil {
		return err
	}

	if !r.headerWritten {
		r.headerWritten = true
	}

	connHeader := determineConnectionHeader(req, serverWantsClose)
	r.Headers["Connection"] = connHeader

	if err = r.writeHeaders(); err != nil {
		return err
	}

	if r.hasContentLength && len(r.Body) != r.contentLength {
		return errors.New("actual body size does not match Content-Length")
	}

	if err = r.checkCancel(); err != nil {
		return err
	}

	// write deadline before writing
	r.Conn.SetWriteDeadline(time.Now().Add(r.Cfg.WriteTimeout))

	if _, err = safeWrite(r.Conn, r.Body); err != nil {
		return err
	}
	return nil
}
