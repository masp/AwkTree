package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestCases = []struct {
	in   string
	want string
}{
	{
		"`foo` { print(@) }",
		`
PATTERN(` + "`foo`" + `)
LCURLY_BRACKET({)
	IDENT(print)
	LPAREN(()
		IDENT(@)
	RPAREN())
RCURLY_BRACKET(})
`,
	},
	{
		`(dotted_name
			(_) @prev-id
			.
			(identifier) @next-id)`,
		`
LPAREN(()
	IDENT(dotted_name)
	LPAREN(()
		IDENT(_)
	RPAREN())
	IDENT(@prev-id)
	PERIOD(.)
	LPAREN(()
		IDENT(identifier)
	RPAREN())
	IDENT(@next-id)
RPAREN())
		`,
	},
}

func TestLexT(t *testing.T) {
	for _, tt := range TestCases {
		t.Run(tt.in, func(t *testing.T) {
			tokens, err := Lex([]byte(tt.in))
			if err != nil {
				t.Fatal(err)
			}
			var got []string
			for _, tok := range tokens {
				got = append(got, encodeToken(tok))
			}
			want := strings.Fields(tt.want)
			assert.Equal(t, want, got, "got: %+v", got)
		})
	}
}

func encodeToken(t Token) string {
	return t.Type.String() + "(" + t.Lit + ")"
}
