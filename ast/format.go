package ast

import (
	"bytes"
	"fmt"
)

// Format creates reverses parse and creates a valid program from a given AST.
func Format(x Node) string {
	var buf bytes.Buffer
	format(x, &buf)
	return buf.String()
}

// TODO: need to actually format the code properly and not just dump each action pattern to a single line
func format(x Node, buf *bytes.Buffer) {
	switch x := x.(type) {
	case *Program:
		for _, pa := range x.Patterns {
			format(pa, buf)
			fmt.Fprintf(buf, "\n")
		}
	case *PatternAction:
		format(x.Pattern, buf)
		format(x.Action, buf)
	case *QueryPattern:
		buf.WriteString("(")
		format(x.Symbol, buf)
		for _, arg := range x.Args {
			buf.WriteString(" ")
			format(arg, buf)
		}
		buf.WriteString(")")
		if x.Capture != nil {
			buf.WriteString(" ")
			format(x.Capture, buf)
			buf.WriteString(" ")
		}
	case *PatternField:
		format(x.Name, buf)
		buf.WriteString(":")
	case *Action:
		buf.WriteString("{")
		for i, stmt := range x.Stmts {
			if i > 0 {
				buf.WriteString(";")
			}
			format(stmt, buf)
		}
		buf.WriteString("}\n")
	case *Call:
		format(x.FuncName, buf)
		buf.WriteString("(")
		for i, arg := range x.Args {
			if i > 0 {
				buf.WriteString(",")
			}
			format(arg, buf)
		}
		buf.WriteString(")")
	case *String:
		buf.WriteString(`"` + x.Value + `"`)
	case *Ident:
		buf.WriteString(x.Name)
	case *Dict:
		buf.WriteString("{")
		for _, kv := range x.Entries {
			format(kv.Key, buf)
			if kv.Colon.IsValid() {
				buf.WriteString(":")
			}
			format(kv.Val, buf)
			if kv.Comma.IsValid() {
				buf.WriteString(",")
			}
		}
		buf.WriteString("}")
	}
}
