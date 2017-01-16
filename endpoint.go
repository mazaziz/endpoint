package endpoint

import (
	"net/http"
	"regexp"
	"fmt"
	"strings"
	"net"
	"context"
	"os"
)

type endpoint struct {
	root *node
	matchers map[string]*regexp.Regexp
	Handler404 http.HandlerFunc
	Handler405 http.HandlerFunc
}

func New() *endpoint {
	e := &endpoint{}
	e.root = e.newNode()
	e.Handler404 = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not Found", http.StatusNotFound)
	}
	e.Handler405 = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	return e
}

func (e *endpoint) GetPaths() []string {
	return e.root.getPaths("/")
}

func (e *endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p := newPath(r.URL.Path)
	n := e.root
	s := p.next()
	for s != nil {
		if n.staticLinks != nil {
			n = n.staticLinks[*s]
		} else if n.dynamicMatcher != nil && n.dynamicMatcher.MatchString(*s) {
			r = r.WithContext(context.WithValue(r.Context(), n.dynamicParam, *s))
			n = n.dynamicLink
		} else {
			n = nil
		}
		if n == nil {
			e.Handler404(w, r)
			return
		}
		s = p.next()
	}
	
	if n.methods == nil {
		e.Handler404(w, r)
	} else if r.Method == "OPTIONS" {
		allow := make([]string, 0)
		for m, _ := range n.methods {
			allow = append(allow, m)
		}
		w.Header().Set("Allow", strings.Join(allow, ","))
	} else if n.methods[r.Method] == nil {
		e.Handler405(w, r)
	} else {
		n.methods[r.Method](w, r)
	}
}

func (e *endpoint) Path(path string) *node {
	p := newPath(path)
	n := e.root
	s := p.next()
	for s != nil {
		if strings.HasPrefix(*s, ":") {
			n = n.makeDynamicLink((*s)[1:])
		} else {
			n = n.makeStaticLink(*s)
		}
		s = p.next()
	}
	return n
}

func (e *endpoint) Route(req string, h http.HandlerFunc) {
	s := strings.Fields(req)
	if 2 != len(s) {
		panic("invalid request definition, must be in format 'METHOD PATH'")
	}
	n := e.Path(s[1])
	n.Handle(s[0], h)
}

func (e *endpoint) newNode() *node {
	n := &node{}
	n.endpoint = e
	return n
}

func (e *endpoint) getMatcher(name string) *regexp.Regexp {
	if e.matchers == nil { 
		return nil
	}
	return e.matchers[name]
}

func (e *endpoint) Match(name string, expr string) {
	if name == "" {
		panic("empty matcher name")
	}
	if e.matchers == nil { 
		e.matchers = make(map[string]*regexp.Regexp)
	}
	if e.matchers[name] != nil {
		panic("already registered matcher (" + name + ")")
	}
	matcher, err := regexp.Compile(expr)
	if err != nil {
		panic("invalid match pattern, " + err.Error())
	}
	e.matchers[name] = matcher
}

func (e *endpoint) Serve(laddr string) {
	l, err := net.Listen("tcp4", laddr)
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(os.Stdout, "server started, listening on " + laddr)  
	http.Serve(l, e)
}
