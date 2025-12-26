package router

import (
	"fmt"
	"strings"

	"github.com/brutally-Honest/http-server/internal/request"
	"github.com/brutally-Honest/http-server/internal/response"
)

type Handler func(req *request.Request, res *response.Response)

type RouteMatcher interface {
	Match(method, path string) (Handler, map[string]string, error)
	Register(method, path string, handler Handler)
}

type Node struct {
	segment       string
	children      map[string]*Node
	paramChild    *Node
	wildcardChild *Node
	handlers      map[string]Handler
}

type Router struct {
	root *Node
}

func NewRouter() *Router {
	return &Router{
		root: &Node{
			children: make(map[string]*Node),
			handlers: make(map[string]Handler),
		},
	}
}

type nodeType uint8

const (
	static nodeType = iota
	param
	wildcard
)

func splitPath(path string) ([]string, error) {
	if path == "" || path[0] != '/' {
		return nil, fmt.Errorf("invalid path")
	}
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}, nil
	}

	segments := strings.Split(path, "/")

	for _, segment := range segments {
		if segment == "" {
			return nil, fmt.Errorf("invalid segment in path")
		}
	}
	return segments, nil
}

func getSegmentType(segment string) nodeType {
	switch segment[0] {
	case ':':
		return param
	case '*':
		return wildcard
	default:
		return static
	}
}

func (r *Router) insert(method string, segments []string, handler Handler) error {
	curr := r.root

	for idx, segment := range segments {
		segmentType := getSegmentType(segment)

		switch segmentType {
		case static:
			if _, exists := curr.children[segment]; !exists {
				curr.children[segment] = &Node{
					children: make(map[string]*Node),
					handlers: make(map[string]Handler),
					segment:  segment,
				}
			}
			curr = curr.children[segment]
		case param:
			if curr.paramChild == nil {
				curr.paramChild = &Node{
					children: make(map[string]*Node),
					handlers: make(map[string]Handler),
					segment:  segment,
				}
			} else if curr.paramChild.segment != segment {
				return fmt.Errorf("conflicting param route")
			}
			curr = curr.paramChild
		case wildcard:
			if idx != len(segments)-1 {
				return fmt.Errorf("invalid route: wildcard %s must be final segment", segment)
			}
			if curr.wildcardChild == nil {
				curr.wildcardChild = &Node{
					children: make(map[string]*Node),
					handlers: make(map[string]Handler),
					segment:  segment,
				}
			} else if curr.wildcardChild.segment != segment {
				return fmt.Errorf("conflicting wildcard route")
			}
			curr = curr.wildcardChild
		}
	}
	if curr.handlers[method] != nil {
		return fmt.Errorf("handler already registered for %s %s", method, segments)
	}

	curr.handlers[method] = handler
	return nil
}

func (r *Router) Register(method, path string, handler Handler) {
	segments, err := splitPath(path)
	if err != nil {
		panic(err)
	}
	if err := r.insert(method, segments, handler); err != nil {
		panic(err.Error())
	}
}

func (r *Router) GET(path string, handler Handler) {
	r.Register("GET", path, handler)
}

func (r *Router) POST(path string, handler Handler) {
	r.Register("POST", path, handler)
}

func (r *Router) PUT(path string, handler Handler) {
	r.Register("PUT", path, handler)
}

func (r *Router) DELETE(path string, handler Handler) {
	r.Register("DELETE", path, handler)
}

func (r *Router) search(method string, segments []string) (Handler, map[string]string) {
	curr := r.root
	params := make(map[string]string)

	for idx, segment := range segments {

		if _, exists := curr.children[segment]; exists {
			curr = curr.children[segment]
			continue
		}

		if curr.paramChild != nil {
			paramChild := curr.paramChild.segment[1:] //removing : from registered route
			params[paramChild] = segment
			curr = curr.paramChild
			continue
		}

		if curr.wildcardChild != nil {
			wildcardChild := curr.wildcardChild.segment[1:]
			params[wildcardChild] = strings.Join(segments[idx:], "/") //removing * from registered route
			curr = curr.wildcardChild
			break
		}

		return nil, nil

	}
	handler := curr.handlers[method]
	if handler == nil {
		return nil, nil
	}

	return handler, params
}

func (r *Router) Match(method, path string) (Handler, map[string]string, error) {
	segments, err := splitPath(path)
	if err != nil {
		return nil, nil, err
	}

	handler, params := r.search(method, segments)
	if handler == nil {
		return nil, nil, fmt.Errorf("route not found")
	}

	return handler, params, nil
}
