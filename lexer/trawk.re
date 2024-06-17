package lexer

import (
    "bytes"
    "github.com/masp/awktree/token"
	"fmt"
)

func (l *Lexer) lex() (pos token.Pos, tok token.Type, lit string, err error) {
    for {
		lit = ""
		pos = l.pos()
		l.token = l.cursor

/*!re2c
		re2c:yyfill:enable = 0;
		re2c:flags:nested-ifs = 1;
		re2c:define:YYCTYPE = byte;
		re2c:define:YYPEEK = "l.input[l.cursor]";
		re2c:define:YYSKIP = "l.cursor += 1";
		re2c:define:YYBACKUP = "l.marker = l.cursor";
		re2c:define:YYRESTORE = "l.cursor = l.marker";

		end = [\x00];
		end { tok = token.EOF; return }
		* { err = fmt.Errorf("%w: %c", ErrUnrecognizedToken, l.token); return }

		// Whitespace and new lines
		eol = ("\r\n" | "\n");
		eol {
			if l.insertSemi() {
				l.cursor = l.token // Has the effect of "inserting" the semicolon in the input
				tok = token.SEMICOLON
				lit = "\n"
				return
			} else {
				l.file.AddLine(l.token)
				continue
			}
		}

        // Skip whitespace
		[ \t]+ {
			continue
		}

		// Comments
		"//" [^\r\n\x00]* { tok = token.COMMENT; lit = l.literal(); return }
		"/*" ([^*\x00] | ("*" [^/]))* "*""/" { tok = token.COMMENT; lit = l.literal(); return }
		"/*" { return l.lexMultiComment() }

		// Keywords
		// "type" { tok = token.TypeKeyword; lit = "type"; return }

		// Operators and punctuation
		"(" { tok = token.LPAREN; lit = "("; return }
		")" { tok = token.RPAREN; lit = ")"; return }
		"{" { tok = token.LCURLY_BRACKET; lit = "{"; return }
		"}" { tok = token.RCURLY_BRACKET; lit = "}"; return }
		"[" { tok = token.LSQUARE_BRACKET; lit = "["; return }
		"]" { tok = token.RSQUARE_BRACKET; lit = "]"; return }
        "==" { tok = token.EQUAL_EQUAL; lit = "=="; return }
        "!=" { tok = token.EQUAL_EQUAL; lit = "!="; return }
        ">=" { tok = token.GREATER_EQUAL; lit = ">="; return }
        "<=" { tok = token.LESS_EQUAL; lit = "<="; return }
        ">" { tok = token.GREATER; lit = ">"; return }
        "<" { tok = token.LESS; lit = "<"; return }
        "+" { tok = token.PLUS; lit = "+"; return }
        "-" { tok = token.MINUS; lit = "-"; return }
        "*" { tok = token.STAR; lit = "*"; return }
        "/" { tok = token.SLASH; lit = "/"; return }

		"." { tok = token.PERIOD; lit = "."; return }
		"," { tok = token.COMMA; lit = ","; return }
		":" { tok = token.COLON; lit = ":"; return }
		";" { tok = token.SEMICOLON; lit = ";"; return }

		// Integer literals
		dec = [1-9][0-9]* | "0";
		dec { tok = token.INT; lit = l.literal(); return }

		// Floating point numbers
		// from excellent https://re2c.org/examples/c/real_world/example_cxx98.html
		frc = [0-9]* "." [0-9]+ | [0-9]+ ".";
		exp = 'e' [+-]? [0-9]+;
		flt = (frc exp? | [0-9]+ exp);
		flt { tok = token.FLOAT; lit = l.literal(); return }

		// Strings
		["] { return l.lexString('"') }
		[`] { return l.lexPattern('`') }

		// Identifiers
		id = [a-zA-Z_$@][a-zA-Z_0-9-]*;
		id { tok = token.IDENT; lit = l.literal(); return }
*/
    }
}

func (l *Lexer) lexString(quote byte) (pos token.Pos, tok token.Type, lit string, err error) {
	var buf bytes.Buffer
	buf.WriteByte(quote)
	for {
		var u byte
/*!re2c
		re2c:yyfill:enable = 0;
		re2c:flags:nested-ifs = 1;
		re2c:define:YYCTYPE = byte;
		re2c:define:YYPEEK = "l.input[l.cursor]";
		re2c:define:YYSKIP = "l.cursor += 1";

		* { err = ErrInvalidString; return }
		[\x00] {
			err = ErrUnterminatedString
			tok = token.EOF
            pos = l.file.Pos(l.token)
			return
		}
		[^\n\\]              {
			u = yych
			buf.WriteByte(u)
			if u == quote {
				tok = token.STRING
				pos = l.file.Pos(l.token)
				lit = string(buf.Bytes())
				return
			}
			continue
		}
		"\\a"                { buf.WriteByte('\a'); continue }
		"\\b"                { buf.WriteByte('\b'); continue }
		"\\f"                { buf.WriteByte('\f'); continue }
		"\\n"                { buf.WriteByte('\n'); continue }
		"\\r"                { buf.WriteByte('\r'); continue }
		"\\t"                { buf.WriteByte('\t'); continue }
		"\\v"                { buf.WriteByte('\v'); continue }
		"\\\\"               { buf.WriteByte('\\'); continue }
		"\\'"                { buf.WriteByte('\''); continue }
		"\\\""               { buf.WriteByte('"'); continue }
		"\\?"                { buf.WriteByte('?'); continue }
*/		
	}
}

func (l *Lexer) lexPattern(quote byte) (pos token.Pos, tok token.Type, lit string, err error) {
	for {
/*!re2c
		re2c:yyfill:enable = 0;
		re2c:flags:nested-ifs = 1;
		re2c:define:YYCTYPE = byte;
		re2c:define:YYPEEK = "l.input[l.cursor]";
		re2c:define:YYSKIP = "l.cursor += 1";

		[\x00] {
			err = ErrUnterminatedString
			tok = token.EOF
            pos = l.file.Pos(l.token)
			return
		}
		[^\x00] {
			if yych == quote {
				tok = token.PATTERN
				pos = l.file.Pos(l.token)
				lit = string(l.input[l.token:l.cursor])
				return
			}
			continue
		}
*/		
	}
}

func (l *Lexer) lexMultiComment() (pos token.Pos, tok token.Type, lit string, err error) {
	for {
/*!re2c
		re2c:yyfill:enable = 0;
		re2c:flags:nested-ifs = 1;
		re2c:define:YYCTYPE = byte;
		re2c:define:YYPEEK = "l.input[l.cursor]";
		re2c:define:YYSKIP = "l.cursor += 1";

		[\x00] {
			err = ErrUnterminatedComment
			tok = token.EOF
            pos = l.file.Pos(l.token)
			return
		}
		"*/" {
			tok = token.COMMENT
			pos = l.file.Pos(l.token)
			lit = string(l.input[l.token+2:l.cursor])
			return
		}
		[^\x00] { continue }
*/		
	}
}