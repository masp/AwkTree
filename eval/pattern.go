package eval

import (
	"fmt"
	"slices"
	"strings"

	"github.com/masp/awktree/ast"
	sitter "github.com/smacker/go-tree-sitter"
)

func (l *language) formatPattern(pa *ast.QueryPattern) (rootCapture string, tsPattern string, err error) {
	l.replaceSymbols(pa)
	tsPattern = ast.Format(pa)
	if pa.Capture != nil {
		rootCapture = pa.Capture.Name
	} else {
		tsPattern += " @__match" // implicit variable for the entire match unless there already is one
		rootCapture = "@__match"
	}
	return
}

// To support abbrevations of symbols, we define the sorted list of all symbols for a given language,
// and then do a longest-prefix search to find the appropriate symbol and replace it in the resulting
// query pattern.

type symbol string

type language struct {
	lang *sitter.Language

	symbols []symbol // sorted list of all symbols
}

func buildLanguage(lang *sitter.Language) *language {
	return &language{
		lang:    lang,
		symbols: allSymbols(lang),
	}
}

func allSymbols(lang *sitter.Language) (result []symbol) {
	for i := uint32(0); i < lang.SymbolCount(); i++ {
		result = append(result, symbol(lang.SymbolName(sitter.Symbol(i))))
	}
	slices.Sort(result)
	return slices.Compact(result)
}

func (l *language) replaceSymbols(x ast.Node) error {
	return ast.Walk(x, ast.VisitorFunc(func(x ast.Node) error {
		switch x := x.(type) {
		case *ast.QueryPattern:
			if isPunctuation(x.Symbol.Name) {
				return nil
			}
			opts := l.lookupAbbrev(x.Symbol.Name)
			if len(opts) > 1 {
				return fmt.Errorf("ambiguous symbol abbreviation %s (%+v)", x.Symbol.Name, opts)
			}
			x.Symbol.Name = string(opts[0])
		}
		return nil
	}))
}

// isPunctuation is true if the symbol is plain ASCII letters (e.g. identifier) and not
//
//	_, !=, ==, or some symbol similar.
func isPunctuation(sym string) bool {
	return sym[0] == '_'
}

// lookupAbbrev will look up all the potential matches abbrev symbol has in the language.
// If there are multiple matches, the abbrev is ambiguous.
func (l *language) lookupAbbrev(abbrev string) []symbol {
	// Search the symbols by the first part first. For all matches,
	// we then check the prefix of the rest.
	parts := strings.Split(abbrev, "_")
	i, found := slices.BinarySearch(l.symbols, symbol(parts[0]))
	if found {
		return []symbol{l.symbols[i]}
	}

	// Otherwise, search the candidates for the longest prefix match. If abbrev is a prefix
	// of a symbol, BinarySearch will pick the index of the first symbol that has the prefix.
	bestScore := 1
	var best []symbol
	for j := i; j < len(l.symbols); j++ {
		score, eol := matchAbbrev(l.symbols[j], symbol(abbrev))
		if eol {
			break
		}
		if score >= bestScore {
			bestScore = score
			best = append(best, l.symbols[j])
		}
	}
	return best
}

func matchAbbrev(sym symbol, abbrev symbol) (matching int, eol bool) {
	abbrevParts := strings.Split(string(abbrev), "_")
	symParts := strings.Split(string(sym), "_")
	if len(symParts) != len(abbrevParts) {
		return 0, false
	}

	for i := 0; i < len(abbrevParts); i++ {
		if strings.HasPrefix(symParts[i], abbrevParts[i]) {
			matching += len(abbrevParts[i])
		} else {
			return 0, true
		}
	}
	return matching, false
}
