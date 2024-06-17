package ast

type Visitor interface {
	// Visit traverses the AST in depth-first order (every statement is completely traversed before the next).
	// If the visitor returns an error, then the traversal stops and the error is returned to the caller.
	Visit(n Node) error
}

type visitorFunc func(Node) error

func (v visitorFunc) Visit(n Node) error {
	return v(n)
}

func VisitorFunc(fn func(Node) error) Visitor {
	return visitorFunc(fn)
}

type bailout error

// mustVisit panics if the walk fails, expecting to be caught by the recover in Walk.
// This is because writing the err != nil is very repetitive and obscures the actual logic.
func mustVisit(v Visitor, n Node) {
	err := v.Visit(n)
	if err != nil {
		panic(bailout(err))
	}
}

func walk(n Node, v Visitor) {
	mustVisit(v, n) // visit root first

	switch n := n.(type) {
	case *Program:
		for _, pattern := range n.Patterns {
			walk(pattern, v)
		}
	case *PatternAction:
		walk(n.Pattern, v)
		walk(n.Action, v)
	case *QueryPattern:
		mustVisit(v, n.Symbol)
		for _, arg := range n.Args {
			walk(arg, v)
		}
		if n.Capture != nil {
			mustVisit(v, n.Capture)
		}
	case *Action:
		for _, stmt := range n.Stmts {
			walk(stmt, v)
		}
	default:
		// leaf node, no need to do anything
		return
	}
}

// Walk traverses each node calling v.Walk for each one, including the root first. If v.Walk returns an error,
// then the traversal stops and the error is returned to the caller.
func Walk(n Node, v Visitor) (result error) {
	defer func() {
		err := recover()
		if berr, ok := err.(bailout); ok {
			result = berr
		} else if err != nil {
			panic(err)
		}
	}()
	walk(n, v)
	return
}
