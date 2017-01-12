package endpoint

import (
	"net/http"
	"regexp"
)

type node struct {
	endpoint *endpoint
	methods map[string]http.HandlerFunc
	staticLinks map[string]*node
	dynamicLink *node
	dynamicMatcher *regexp.Regexp
	dynamicParam string
}

func (n *node) makeDynamicLink(name string) *node {
	if n.staticLinks != nil {
		panic("cant make dynamic link when there are static links")
	}
	matcher := n.endpoint.matchers[name]
	if matcher == nil {
		panic("unknown matcher (" + name + ")")
	}
	if n.dynamicMatcher == nil {
		n.dynamicMatcher = matcher
		n.dynamicLink = n.endpoint.newNode()
		n.dynamicParam = name
	} else if n.dynamicParam != name {
		panic("another dynamic link registered on this path")
	}
	return n.dynamicLink
}

func (n *node) Handle(method string, h http.HandlerFunc) *node {
	if n.methods == nil {
		n.methods = make(map[string]http.HandlerFunc)
	}
	if n.methods[method] != nil {
		panic("method handler already registered")
	}
	n.methods[method] = h
	return n
}

func (n *node) makeStaticLink(segment string) *node {
	if n.dynamicMatcher != nil {
		panic("cant make static link when there is dynamic link")
	}
	if n.staticLinks == nil {
		n.staticLinks = make(map[string]*node)
	}
	if n.staticLinks[segment] == nil {
		n.staticLinks[segment] = n.endpoint.newNode()
	}
	return n.staticLinks[segment]
}
