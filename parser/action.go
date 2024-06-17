package parser

import (
	"github.com/masp/awktree/ast"
	"github.com/masp/awktree/token"
)

var (
	exprStart = map[token.Type]bool{
		token.IDENT:          true,
		token.INT:            true,
		token.STRING:         true,
		token.BANG:           true,
		token.LCURLY_BRACKET: true,
		token.FLOAT:          true,
	}

	exprEnd = map[token.Type]bool{
		token.EOF:            true,
		token.SEMICOLON:      true,
		token.RPAREN:         true,
		token.RCURLY_BRACKET: true,
		token.COMMA:          true,
	}
)

func (p *Parser) parseAction() *ast.Action {
	action := &ast.Action{
		OpenCurly: p.expect(token.LCURLY_BRACKET).Pos,
	}
	for {
		t := p.peek()
		if t.Type == token.RCURLY_BRACKET || t.Type == token.EOF {
			break
		}
		stmt := p.parseStmt()
		if stmt == nil {
			break
		}
		action.Stmts = append(action.Stmts, stmt)
	}
	action.CloseCurly = p.expect(token.RCURLY_BRACKET).Pos
	return action
}

func (p *Parser) parseStmt() ast.Stmt {
	t := p.peek()
	switch t.Type {
	case token.IDENT:
		return p.parseCall()
	default:
		p.errorf(t.Pos, "unexpected token %s, wanted statement", t.String())
	}
	return nil
}

func (p *Parser) parseCall() *ast.Call {
	call := &ast.Call{FuncName: p.parseIdent()}
	call.Lparen = p.expect(token.LPAREN).Pos
	call.Args = p.parseArgs()
	call.Rparen = p.expect(token.RPAREN).Pos
	return call
}

func (p *Parser) parseArgs() []ast.Expr {
	var args []ast.Expr
	for {
		t := p.peek()
		if t.Type == token.RPAREN {
			break
		}
		args = append(args, p.parseExpr())
	}
	return args
}

func (p *Parser) parseExpr() ast.Expr {
	t := p.peek()
	switch t.Type {
	case token.IDENT:
		return p.parseIdent()
	case token.STRING:
		return p.parseString()
	case token.LCURLY_BRACKET:
		return p.parseDict()
	default:
		p.errorf(t.Pos, "unexpected token %s, wanted expression", t.String())
	}
	return nil
}

func (p *Parser) parseString() *ast.String {
	tok := p.expect(token.STRING)
	return &ast.String{ValuePos: tok.Pos, Value: tok.Lit[1 : len(tok.Lit)-1]}
}

func (p *Parser) parseIdent() *ast.Ident {
	tok := p.expect(token.IDENT)
	return &ast.Ident{NamePos: tok.Pos, Name: tok.Lit}
}

func (p *Parser) parseDict() ast.Expr {
	dict := &ast.Dict{}
	dict.LCurly = p.expect(token.LCURLY_BRACKET).Pos
entries:
	for {
		next := p.peek()
		switch next.Type {
		case token.RCURLY_BRACKET:
			break entries
		case token.EOF:
			here := p.peek().Pos
			p.errorf(here, "expected '}', unterminated dict starting at %s", p.file.Position(dict.LCurly))
			return &ast.BadExpr{From: dict.LCurly, To: here}
		default:
			if _, ok := exprStart[next.Type]; !ok {
				p.errorf(p.peek().Pos, "expected entry, got %s", next.Type)
				to := p.advance(exprEnd)
				return &ast.BadExpr{From: dict.LCurly, To: to.Pos}
			}
		}

		var entry ast.DictEntry
		entry.Key = p.parseExpr()
		entry.Colon = p.expect(token.COLON).Pos
		entry.Val = p.parseExpr()
		if p.peek().Type == token.COMMA {
			entry.Comma = p.eat().Pos
		}
		dict.Entries = append(dict.Entries, &entry)
	}
	dict.RCurly = p.expect(token.RCURLY_BRACKET).Pos
	return dict
}
