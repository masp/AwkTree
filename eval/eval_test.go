package eval

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvalExample(t *testing.T) {
	prog, err := Compile("<test>", []byte(`(identifier) @id {print(@)}`))
	if err != nil {
		t.Fatalf("compile error: %v", err)
	}

	var stdout bytes.Buffer
	err = prog.Eval(context.Background(), []byte(`let a = 10;`), &Options{
		Language: javascript.GetLanguage(),
		Stdout:   &stdout,
	})
	if err != nil {
		t.Errorf("eval error: %v", err)
	}
	t.Logf("\n=== OUTPUT ===\n%s\n==============", stdout.String())
	assert.Equal(t, "a\n", stdout.String(), "unexpected output")
}

func TestUseCases(t *testing.T) {
	tests := []struct {
		filename string
	}{
		{"testdata/JavaBasic/Example.java"},
		{"testdata/Abbrev/test.py"},
		{"testdata/Unknown/test.go"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			getTestFile := func(ext string) string {
				return strings.TrimSuffix(tt.filename, filepath.Ext(tt.filename)) + ext
			}

			src, err := os.ReadFile(tt.filename)
			if err != nil {
				t.Fatalf("reading test source: %v", err)
			}

			// Change suffix of filename to .tra
			traSrc, err := os.ReadFile(getTestFile(".tra"))
			if err != nil {
				t.Fatalf("reading pattern: %v", err)
			}

			prog, err := Compile(getTestFile(".tra"), traSrc)
			if err != nil {
				t.Fatalf("compile error: %v", err)
			}

			var stdout bytes.Buffer
			err = prog.Eval(context.Background(), src, &Options{
				Filename: tt.filename,
				Stdout:   &stdout,
			})
			if stdout.Len() == 0 {
				stdout.WriteString("<NO OUTPUT>\n")
			}
			if err != nil {
				t.Fatalf("eval error: %v", err)
			}

			wantFile := getTestFile(".out")
			gotFile := getTestFile(".got")
			if _, err := os.Stat(wantFile); os.IsNotExist(err) {
				t.Logf("want file does not exist, updating test")
				require.NoError(t, os.WriteFile(wantFile, stdout.Bytes(), 0644))
			} else {
				require.NoError(t, os.WriteFile(gotFile, stdout.Bytes(), 0644))
			}
			want, err := os.ReadFile(wantFile)
			if err != nil {
				t.Fatalf("reading want: %v", err)
			}
			require.Equal(t, string(want), stdout.String(), "unexpected output, want %s, got %s", wantFile, gotFile)
		})
	}
}
