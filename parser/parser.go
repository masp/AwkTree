package parser

import (
	"fmt"

	"github.com/masp/awktree/ast"
	"github.com/masp/awktree/lexer"
	"github.com/masp/awktree/token"
)

type Parser struct {
	tokens []lexer.Token
	file   *token.File
	pos    int

	errors token.ErrorList
}

type Options struct{}

func ParseFile(filename string, src []byte, opts *Options) (prog *ast.Program, err error) {
	if opts == nil {
		opts = &Options{}
	}
	lex := lexer.NewLexer(filename, src)
	prog = &ast.Program{File: lex.File()}
	tokens := lex.All()
	if lex.HasErrors() {
		err = lex.Errors()
		return
	}

	p := &Parser{
		file:   lex.File(),
		tokens: tokens,
	}

	defer func() {
		errlist := p.catchErrors()
		errlist.Sort()
		if errlist.Len() > 0 {
			err = errlist.Err()
		}
	}()

	for {
		tok := p.peek()
		if tok.Type == token.EOF {
			break
		}

		switch tok.Type {
		case token.LPAREN:
			patternAction := p.parsePatternAction()
			prog.Patterns = append(prog.Patterns, patternAction)
		default:
			p.error(tok.Pos, fmt.Errorf("unexpected token %s, wanted pattern or action block", tok.String()))
		}
	}
	return
}

func (p *Parser) eat() (tok lexer.Token) {
	for ; p.pos < len(p.tokens); p.pos++ {
		tok = p.tokens[p.pos]
		if tok.Type == token.COMMENT {
			continue
		}
		p.pos++
		return tok
	}
	return lexer.Token{Type: token.EOF}
}

func (p *Parser) eatAll(tokenType token.Type) token.Type {
	for {
		if p.pos >= len(p.tokens) {
			return token.EOF
		}
		if p.tokens[p.pos].Type == tokenType {
			p.pos++
		} else {
			return tokenType
		}
	}
}

func (p *Parser) expect(tokenType token.Type) (tok lexer.Token) {
	tok = p.eat()
	if tok.Type != tokenType {
		p.error(tok.Pos, fmt.Errorf("expected %s, got %s", tokenType.String(), tok.String()))
	}
	return
}

func (p *Parser) advance(to map[token.Type]bool) (tok lexer.Token) {
	for p.peek().Type != token.EOF && !to[p.peek().Type] {
		tok = p.eat()
	}
	return
}

func (p *Parser) peek() (tok lexer.Token) {
	for i := p.pos; i < len(p.tokens); i++ {
		tok = p.tokens[i]
		if tok.Type == token.COMMENT {
			continue
		}
		return tok
	}
	return lexer.Token{Type: token.EOF}
}

func (p *Parser) peekN(n int) (toks []lexer.Token) {
	for i := p.pos; i < len(p.tokens); i++ {
		tok := p.tokens[i]
		if tok.Type == token.COMMENT {
			continue
		}
		toks = append(toks, tok)
		if len(toks) == n {
			return toks
		}
	}
	if len(toks) < n {
		for i := 0; i < n-len(toks); i++ {
			toks = append(toks, lexer.Token{Type: token.EOF})
		}
	}
	return toks
}

func (p *Parser) matches(types ...token.Type) bool {
	next := p.peek()
	for _, t := range types {
		if next.Type == t {
			return true
		}
	}
	return false
}
