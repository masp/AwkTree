package eval

import (
	"testing"

	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/stretchr/testify/assert"
)

func TestMatchAbbreviation(t *testing.T) {
	l := buildLanguage(javascript.GetLanguage())

	// Convert the above cases into a table based test
	tests := []struct {
		abbrev string
		want   []symbol
	}{
		{"array", []symbol{"array"}},
		{"arr", []symbol{"array"}},
		{"ar_pat", []symbol{"array_pattern"}},
		{"ar_pat_rep", []symbol{"array_pattern_repeat1"}},
		{"id", []symbol{"identifier"}},
		{"ca", []symbol{"case", "catch"}},
	}

	for _, tt := range tests {
		t.Run(tt.abbrev, func(t *testing.T) {
			got := l.lookupAbbrev(tt.abbrev)
			assert.Equal(t, tt.want, got)
		})
	}
}
