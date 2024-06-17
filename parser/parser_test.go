package parser

import (
	"strings"
	"testing"

	"github.com/masp/awktree/ast"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []string{
		`(identifier){print(@)}`,
		`(binary_expression operator: "!=" right: (null)){}`,
		`(id){print({id:"test",id2:@,id3:{id4:"test"}})}`,
	}

	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			prog, err := ParseFile("<test>", []byte(tt), nil)
			if err != nil {
				t.Fatal(err)
			}
			got := ast.Format(prog)
			assert.Equal(t, tt, strings.TrimSpace(got))
		})
	}
}
