package eval

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"

	"github.com/go-enry/go-enry/v2"
	"github.com/masp/awktree/ast"
	"github.com/masp/awktree/parser"
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
)

type Program struct {
	Ast *ast.Program
}

func Compile(filename string, src []byte) (*Program, error) {
	prog, err := parser.ParseFile(filename, src, nil)
	if err != nil {
		return nil, err
	}
	return &Program{Ast: prog}, nil
}

type Options struct {
	Filename string           // optional, can be used to detect language
	Language *sitter.Language // optional, overrides from filename

	Stdout io.Writer // optional, defaults to /dev/null
}

func (p *Program) Eval(ctx context.Context, src []byte, opts *Options) error {
	tsLang, err := p.detectLanguage(opts)
	if err != nil {
		return err
	}
	lang := buildLanguage(tsLang)

	n, err := sitter.ParseCtx(ctx, src, tsLang)
	if err != nil {
		return err
	}

	for _, pa := range p.Ast.Patterns {
		var rootCapture string
		rootCapture, tsPattern, err := lang.formatPattern(pa.Pattern)
		if err != nil {
			return err
		}

		log.Printf("query pattern: %s", tsPattern)
		q, err := sitter.NewQuery([]byte(tsPattern), tsLang)
		if err != nil {
			return err
		}
		qc := sitter.NewQueryCursor()
		qc.Exec(q, n)
		state := &evalCtx{}
		if opts.Stdout == nil {
			state.Output = io.Discard
		} else {
			state.Output = opts.Stdout
		}
		state.Src = src
		state.Root = n
		state.SQuery = q
		for {
			state.Clear()
			m, ok := qc.NextMatch()
			if !ok {
				break
			}
			state.applyQuery(rootCapture, m)
			err := p.runAction(state, pa.Action)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Program) detectLanguage(opts *Options) (*sitter.Language, error) {
	if opts.Language != nil {
		return opts.Language, nil
	}

	if opts.Filename == "" {
		return nil, fmt.Errorf("opts.Filename is required")
	}

	lang, safe := enry.GetLanguageByExtension(filepath.Ext(opts.Filename))
	if !safe {
		return nil, fmt.Errorf("multiple languages possible from filename, manually specify opts.Language")
	}
	switch strings.ToLower(lang) {
	case "javascript":
		return javascript.GetLanguage(), nil
	case "go":
		return golang.GetLanguage(), nil
	case "python":
		return python.GetLanguage(), nil
	case "java":
		return java.GetLanguage(), nil
		// TODO: Add more languages
	default:
		return nil, fmt.Errorf("unsupported language %s, manually specify opts.Language", lang)
	}
}

type evalCtx struct {
	Src    []byte
	Root   *sitter.Node
	SQuery *sitter.Query

	// Variable name to value
	Vars map[string]Value

	Output io.Writer
}

func (c *evalCtx) applyQuery(rootVarName string, qm *sitter.QueryMatch) {
	for _, capture := range qm.Captures {
		name := c.SQuery.CaptureNameForId(capture.Index)
		if name == "" {
			continue
		}
		c.Vars["@"+name] = &NodeVal{N: capture.Node, Src: c.Src}
		if name == rootVarName[1:] {
			c.Vars["@"] = &NodeVal{N: capture.Node, Src: c.Src} // store in root variable
		}
	}
	if _, ok := c.Vars["@"]; !ok {
		panic(fmt.Sprintf("root var %s not found, must always be set", rootVarName))
	}
}

func (c *evalCtx) Clear() {
	c.Vars = make(map[string]Value)
}

func (p *Program) runAction(c *evalCtx, action *ast.Action) error {
	for _, stmt := range action.Stmts {
		switch stmt := stmt.(type) {
		case *ast.Call:
			return p.runFunc(c, stmt)
		default:
			return fmt.Errorf("unexpected statement type %T", stmt)
		}
	}
	return nil
}

func (p *Program) runFunc(c *evalCtx, f *ast.Call) error {
	switch f.FuncName.Name {
	case "print":
		if len(f.Args) != 1 {
			return fmt.Errorf("print expects 1 argument, got %d", len(f.Args))
		}
		arg := f.Args[0]
		val, err := p.eval(c, arg)
		if err != nil {
			return err
		}
		_, err = p.print(c, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Program) eval(c *evalCtx, expr ast.Expr) (Value, error) {
	switch expr := expr.(type) {
	case *ast.String:
		return &StringVal{S: expr.Value}, nil
	case *ast.Ident:
		if v, ok := c.Vars[expr.Name]; ok {
			return v, nil
		} else {
			return nil, fmt.Errorf("unknown variable %s", expr.Name)
		}
	case *ast.Dict:
		return p.evalDict(c, expr)
	default:
		return nil, fmt.Errorf("unexpected expression type %T", expr)
	}
}

func (p *Program) evalDict(c *evalCtx, d *ast.Dict) (Value, error) {
	var err error
	m := make(map[Value]Value)
	for _, kv := range d.Entries {
		var key Value
		if k, ok := kv.Key.(*ast.Ident); ok {
			if v, ok := c.Vars[k.Name]; ok {
				key = v
			} else {
				key = &StringVal{S: k.Name}
			}
		} else {
			key, err = p.eval(c, kv.Key)
			if err != nil {
				return nil, err
			}
		}
		val, err := p.eval(c, kv.Val)
		if err != nil {
			return nil, err
		}
		m[key] = val
	}
	return &DictVal{D: m}, nil
}

func (p *Program) print(c *evalCtx, v Value) (int, error) {
	switch v := v.(type) {
	case *StringVal:
		return fmt.Fprintf(c.Output, "%s\n", v.S)
	case *IntVal:
		return fmt.Fprintf(c.Output, "%d\n", v.I)
	case *NodeVal:
		if n, err := c.Output.Write(c.Src[v.N.StartByte():v.N.EndByte()]); err != nil {
			return n, err
		}
		return c.Output.Write([]byte{'\n'})
	case *DictVal:
		marshalled, err := json.Marshal(v)
		if err != nil {
			return 0, err
		}
		marshalled = append(marshalled, '\n')
		return c.Output.Write(marshalled)
	default:
		return fmt.Fprintf(c.Output, "<bad:%T>\n", v)
	}
}
