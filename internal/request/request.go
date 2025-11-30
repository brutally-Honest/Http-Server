package request

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    []byte
}
