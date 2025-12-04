package request

import "errors"

var (
	ErrHeaderLimitExceeded = errors.New("header size limit exceeded")
	ErrBodyLimitExceeded   = errors.New("body size limit exceeded")
	ErrConnectionClosed    = errors.New("connection closed by client")
)
