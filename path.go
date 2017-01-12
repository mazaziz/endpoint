package endpoint

import "strings"

type path struct {
	raw string
	pos int
	segments []string
}

func newPath(str string) *path {
	p := &path{raw: str, pos: 0}
	if str == "" || str == "/" {
		p.segments = []string{""}
	} else {
		p.segments = strings.Split(str, "/")	
	}
	return p
}

func (p *path) next() *string {
	if p.pos + 1 >= len(p.segments) {
		return nil
	}
	p.pos++
	return &p.segments[p.pos]
}
