package nf

import (
	"path/filepath"
	"reft-go/parser"
	"testing"
)

func TestIncludes(t *testing.T) {
	filePath := filepath.Join(getTestDataDir(), "nf-testdata", "includes.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}
	includeVisitor := NewIncludeVisitor()
	includeVisitor.VisitBlockStatement(ast.StatementBlock)
	includes := includeVisitor.GetSortedIncludes()
	if len(includes) != 2 {
		t.Fatalf("Expected 2 includes, got %d", len(includes))
	}
}
