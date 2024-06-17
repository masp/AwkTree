package token

import "strconv"

type Type int

const (
	INVALID Type = iota

	COMMENT

	literal_begin
	IDENT   // e.g. $1, foo
	INT     // e.g. 123
	FLOAT   // e.g. 123.456
	STRING  // e.g. "foo"
	PATTERN // e.g. `foo()`
	literal_end

	LCURLY_BRACKET
	RCURLY_BRACKET
	LSQUARE_BRACKET
	RSQUARE_BRACKET
	LPAREN
	RPAREN

	// Operators
	COMMA
	PERIOD
	SEMICOLON
	COLON
	STAR
	PLUS
	MINUS
	SLASH

	BANG
	LESS
	GREATER
	LESS_EQUAL
	GREATER_EQUAL
	EQUAL
	BANG_EQUAL
	EQUAL_EQUAL

	EOF Type = 255 // must be at end
)

var types = [...]string{
	INVALID:         "INVALID",
	COMMENT:         "COMMENT",
	IDENT:           "IDENT",
	INT:             "INT",
	FLOAT:           "FLOAT",
	STRING:          "STRING",
	PATTERN:         "PATTERN",
	LCURLY_BRACKET:  "LCURLY_BRACKET",
	RCURLY_BRACKET:  "RCURLY_BRACKET",
	LSQUARE_BRACKET: "LSQUARE_BRACKET",
	RSQUARE_BRACKET: "RSQUARE_BRACKET",
	LPAREN:          "LPAREN",
	RPAREN:          "RPAREN",
	COMMA:           "COMMA",
	PERIOD:          "PERIOD",
	SEMICOLON:       "SEMICOLON",
	COLON:           "COLON",
	STAR:            "STAR",
	PLUS:            "PLUS",
	MINUS:           "MINUS",
	SLASH:           "SLASH",
	BANG:            "BANG",
	LESS:            "LESS",
	GREATER:         "GREATER",
	LESS_EQUAL:      "LESS_EQUAL",
	GREATER_EQUAL:   "GREATER_EQUAL",
	EQUAL:           "EQUAL",
	BANG_EQUAL:      "BANG_EQUAL",
	EQUAL_EQUAL:     "EQUAL_EQUAL",
	EOF:             "EOF",
}

func (tok Type) String() string {
	s := ""
	if 0 <= tok && tok < Type(len(types)) {
		s = types[tok]
	}
	if s == "" {
		s = "Token(" + strconv.Itoa(int(tok)) + ")"
	}
	return s
}

func (tok Type) IsLiteral() bool {
	return literal_begin < tok && tok < literal_end
}
