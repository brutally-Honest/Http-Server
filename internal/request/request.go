package request

import "context"

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
	Params  map[string]string
	Context context.Context
}
