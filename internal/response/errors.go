package response

import "errors"

var (
	ErrConnectionClosed = errors.New("connection closed by client")
)
