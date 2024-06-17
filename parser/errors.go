package parser

import (
	"fmt"

	"github.com/masp/awktree/token"
)

var ErrBailout = &struct{}{}

func (p *Parser) error(pos token.Pos, err error) {
	p.errors.Add(p.file.Position(pos), err)
}

func (p *Parser) errorf(pos token.Pos, format string, args ...any) {
	p.error(pos, fmt.Errorf(format, args...))
}

func (p *Parser) Errors() token.ErrorList {
	p.errors.RemoveMultiples()
	return p.errors
}

func (p *Parser) HasErrors() bool {
	return p.errors.Len() > 0
}

func (p *Parser) catchErrors() token.ErrorList {
	if r := recover(); r != nil {
		if r == ErrBailout {
			return p.errors
		} else {
			panic(r)
		}
	}
	return p.errors
}
