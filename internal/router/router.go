package router

import (
	"fmt"
	"strings"

	"github.com/brutally-Honest/http-server/internal/request"
	"github.com/brutally-Honest/http-server/internal/response"
)

type Handler func(req *request.Request, res *response.Response)

type RouteMatcher interface {
	Match(path string) (Handler, map[string]string, error)
}

type Node struct {
	segement      string
	children      map[string]*Node
	paramChild    *Node
	wildcardChild *Node
	handler       Handler
}

type Router struct {
	root *Node
}

func NewRouter() *Router {
	return &Router{
		root: &Node{
			children: make(map[string]*Node),
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
			return nil, fmt.Errorf("invalid segement in path")
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

func (r *Router) Insert(segments []string, handler Handler) error {
	curr := r.root

	for idx, segment := range segments {
		segmentType := getSegmentType(segment)

		switch segmentType {
		case static:
			if _, exists := curr.children[segment]; !exists {
				curr.children[segment] = &Node{
					children: make(map[string]*Node),
					segement: segment,
				}
			}
			curr = curr.children[segment]
		case param:
			if curr.paramChild == nil {
				curr.paramChild = &Node{
					children: make(map[string]*Node),
					segement: segment,
				}
			} else if curr.paramChild.segement != segment {
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
					segement: segment,
				}
			} else if curr.wildcardChild.segement != segment {
				return fmt.Errorf("conflicting wildcard route")
			}
			curr = curr.wildcardChild
		}
	}
	curr.handler = handler
	return nil
}

func (r *Router) GET(path string, handler Handler) {
	segments, err := splitPath(path)
	if err != nil {
		panic(err)
	}
	if err := r.Insert(segments, handler); err != nil {
		panic(err.Error())
	}
}

func (r *Router) Search(segments []string) (Handler, map[string]string) {
	curr := r.root
	params := make(map[string]string)

	for idx, segment := range segments {

		if _, exists := curr.children[segment]; exists {
			curr = curr.children[segment]
			continue
		}

		if curr.paramChild != nil {
			paramChild := curr.paramChild.segement[1:] //removing : from registered route
			params[paramChild] = segment
			curr = curr.paramChild
			continue
		}

		if curr.wildcardChild != nil {
			wildcardChild := curr.wildcardChild.segement[1:]
			params[wildcardChild] = strings.Join(segments[idx:], "/") //removing * from registered route
			curr = curr.wildcardChild
			break
		}

		return nil, nil

	}
	return curr.handler, params
}

func (r *Router) Match(path string) (Handler, map[string]string, error) {
	segments, err := splitPath(path)
	if err != nil {
		return nil, nil, err
	}

	handler, params := r.Search(segments)
	if handler == nil {
		return nil, nil, fmt.Errorf("route not found")
	}

	return handler, params, nil
}
