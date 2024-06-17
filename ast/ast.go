package ast

import (
	"github.com/masp/awktree/token"
)

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Program struct {
	File     *token.File
	Patterns []*PatternAction
}

func (p *Program) Pos() token.Pos {
	if len(p.Patterns) == 0 {
		return token.NoPos
	}
	return p.Patterns[0].Pos()
}
func (p *Program) End() token.Pos {
	if len(p.Patterns) == 0 {
		return token.NoPos
	}
	return p.Patterns[len(p.Patterns)-1].End()
}

type PatternAction struct {
	Pattern *QueryPattern
	Action  *Action
}

func (pa *PatternAction) Pos() token.Pos {
	if pa.Pattern != nil {
		return pa.Pattern.Pos()
	}
	return pa.Action.Pos()
}
func (pa *PatternAction) End() token.Pos { return pa.Action.End() }

// QueryPattern is a tree-sitter query pattern like (identifier) @id. It's a basic
// LISP like tree structure.
type QueryPattern struct {
	Lparen  token.Pos
	Symbol  *Ident
	Args    []Node
	Rparen  token.Pos
	Capture *Ident
}

func (qp *QueryPattern) Pos() token.Pos {
	return qp.Lparen
}

func (qp *QueryPattern) End() token.Pos {
	if qp.Capture != nil {
		return qp.Capture.End()
	}
	return qp.Rparen
}

type PatternField struct {
	Name  *Ident
	Colon token.Pos
}

func (pf *PatternField) Pos() token.Pos { return pf.Name.Pos() }
func (pf *PatternField) End() token.Pos { return pf.Colon }

// Action is a block with a series of statements that will be executed when the pattern
// is matched.
type Action struct {
	OpenCurly  token.Pos
	Stmts      []Stmt
	CloseCurly token.Pos
}

func (a *Action) Pos() token.Pos { return a.OpenCurly }
func (a *Action) End() token.Pos { return a.CloseCurly }

type Stmt interface {
	Node
	stmtNode()
}

func (c *Call) stmtNode() {}

type Call struct {
	FuncName *Ident
	Lparen   token.Pos
	Args     []Expr
	Rparen   token.Pos
}

func (c *Call) Pos() token.Pos { return c.FuncName.Pos() }
func (c *Call) End() token.Pos { return c.Rparen }

type Ident struct {
	NamePos token.Pos
	Name    string
}

func (i *Ident) Pos() token.Pos { return i.NamePos }
func (i *Ident) End() token.Pos { return i.NamePos + token.Pos(len(i.Name)) }

type String struct {
	ValuePos token.Pos
	Value    string
}

func (s *String) Pos() token.Pos { return s.ValuePos }
func (s *String) End() token.Pos { return s.ValuePos + token.Pos(len(s.Value)) + 2 } // +2 for quotes

type Dict struct {
	LCurly  token.Pos
	Entries []*DictEntry
	RCurly  token.Pos
}

func (d *Dict) Pos() token.Pos { return d.LCurly }
func (d *Dict) End() token.Pos { return d.RCurly }

type DictEntry struct {
	Key, Val Expr
	Colon    token.Pos
	Comma    token.Pos // if there is a trailing comma
}

func (d *DictEntry) Pos() token.Pos { return d.Key.Pos() }
func (d *DictEntry) End() token.Pos {
	if d.Comma != token.NoPos {
		return d.Comma
	}
	return d.Val.End()
}

type Expr interface {
	Node
	exprNode()
}

type BadExpr struct {
	From, To token.Pos
}

func (b *BadExpr) exprNode() {}
func (b *BadExpr) Pos() token.Pos {
	return b.From
}
func (b *BadExpr) End() token.Pos {
	return b.To
}

func (c *Call) exprNode()   {}
func (i *Ident) exprNode()  {}
func (s *String) exprNode() {}
func (d *Dict) exprNode()   {}
