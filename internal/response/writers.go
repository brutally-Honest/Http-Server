package response

import (
	"errors"
	"io"
	"log"
	"net"
)

func safeWrite(conn net.Conn, buffer []byte) (int, error) {
	n, err := conn.Write(buffer)
	if err != nil {
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) || isConnectionError(err) {
			log.Printf("write error : connection closed")
			return n, ErrConnectionClosed
		}
		log.Printf("write error: %v", err)
		return n, err
	}
	return n, nil
}

func safeWriteString(conn net.Conn, s string) error {
	_, err := safeWrite(conn, []byte(s))
	return err
}
