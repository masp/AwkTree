package ast

import (
	"strings"
	"testing"

	"github.com/masp/awktree/token"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	tree := &Program{
		Patterns: []*PatternAction{
			{
				Pattern: &QueryPattern{
					Lparen: token.Pos(1),
					Symbol: &Ident{
						NamePos: token.Pos(2),
						Name:    "identifier",
					},
					Args: []Node{
						&Ident{
							NamePos: token.Pos(3),
							Name:    "_",
						},
					},
					Rparen: token.Pos(4),
					Capture: &Ident{
						NamePos: token.Pos(5),
						Name:    "@id",
					},
				},
				Action: &Action{
					OpenCurly: token.Pos(6),
					Stmts: []Stmt{
						&Call{
							FuncName: &Ident{
								NamePos: token.Pos(7),
								Name:    "print",
							},
							Args: []Expr{
								&String{
									ValuePos: token.Pos(8),
									Value:    "Hello, World!",
								},
							},
						},
					},
				},
			},
		},
	}

	got := Format(tree)
	assert.Equal(t, `(identifier _) @id {print("Hello, World!")}`, strings.TrimSpace(got))
}
