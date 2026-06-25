package evalop

import (
	"go/ast"
	"go/parser"
	"testing"

	"github.com/go-delve/delve/pkg/dwarf/godwarf"
)

type stubLookup struct{}

func (stubLookup) FindTypeExpr(_ ast.Expr) (godwarf.Type, error) { return nil, nil }
func (stubLookup) HasBuiltin(string) bool                        { return false }
func (stubLookup) PtrSize() int                                    { return 8 }

func TestCompileForkSyntax(t *testing.T) {
	lookup := stubLookup{}
	cases := []struct {
		expr    string
		wantErr bool
	}{
		{`x!`, false},
		{`x!.field`, false},
		{`a ?? b`, false},
		{`root?.child?.name`, false},
		{`if true { 1 } else { 2 }`, true},
		{`switch x { case 1: 2; default: 3 }`, true},
		{`(a, b) => a + b`, true},
	}
	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			_, err := Compile(lookup, tc.expr, 0)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error compiling %q", tc.expr)
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error compiling %q: %v", tc.expr, err)
			}
		})
	}
}

func TestParseForkSyntax(t *testing.T) {
	exprs := []string{
		`x!`,
		`x!.field`,
		`a ?? b`,
		`root?.child?.name`,
		`if true { 1 } else { 2 }`,
	}
	for _, expr := range exprs {
		if _, err := parser.ParseExpr(expr); err != nil {
			t.Errorf("parser.ParseExpr(%q): %v", expr, err)
		}
	}
}
