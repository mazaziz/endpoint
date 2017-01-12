package endpoint

import (
	"strings"
//	"fmt"
//	"os"
)

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

//func (p *path) dump() {
//	util.Dumpl("path", p.raw)
//	util.Dumpl("pos", p.pos)
//	if p.segments == nil {
//		fmt.Println("path not parsed")
//	} else {
//		fmt.Println("path segments =")
//		for i := 0; i < len(p.segments); i++ {
//			fmt.Fprintln(os.Stdout, fmt.Sprintf(" [%d] => '%s'", i, p.segments[i]))
//		}
//	}
//}
