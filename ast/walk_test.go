package ast

import (
	"fmt"
	"strings"
	"testing"

	"github.com/masp/awktree/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalk(t *testing.T) {
	// Given an input AST, make sure that ast.Walk with a print func produces the whole syntax tree
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
				Action: &Action{},
			},
		},
	}

	var output strings.Builder
	err := Walk(tree, VisitorFunc(func(n Node) error {
		fmt.Fprintf(&output, "%T\n", n)
		return nil
	}))
	require.NoError(t, err)
	assert.Equal(t, `*ast.Program
*ast.PatternAction
*ast.QueryPattern
*ast.Ident
*ast.Ident
*ast.Ident
*ast.Action
`, output.String())
}
