package parser

import (
	"strings"

	"github.com/masp/awktree/ast"
	"github.com/masp/awktree/token"
)

func (p *Parser) parsePatternAction() *ast.PatternAction {
	pattern := p.parsePattern()
	action := p.parseAction()
	return &ast.PatternAction{Pattern: pattern, Action: action}
}

func (p *Parser) parsePattern() *ast.QueryPattern {
	pattern := &ast.QueryPattern{Lparen: p.expect(token.LPAREN).Pos}
	pattern.Symbol = p.parseIdent()
loop:
	for {
		var arg ast.Node
		t := p.peek()
		switch t.Type {
		case token.IDENT:
			// Field names
			pf := &ast.PatternField{Name: p.parseIdent()}
			pf.Colon = p.expect(token.COLON).Pos
			arg = pf
		case token.STRING:
			// Anonymouse nodes
			arg = p.parseString()
		case token.LPAREN:
			// Subnode
			arg = p.parsePattern()
		case token.RPAREN:
			break loop
		default:
			p.errorf(t.Pos, "bad token in node pattern: %v", t.Type)
		}
		pattern.Args = append(pattern.Args, arg)
	}
	pattern.Rparen = p.expect(token.RPAREN).Pos
	alias := p.peek()
	if alias.Type == token.IDENT && strings.HasPrefix(alias.Lit, "@") {
		pattern.Capture = p.parseIdent()
	}
	return pattern
}
