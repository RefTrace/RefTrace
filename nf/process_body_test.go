package nf

import (
	"path/filepath"
	"reft-go/parser"
	"testing"
)

func TestProcessInputs(t *testing.T) {
	filePath := filepath.Join("./testdata", "process_inputs.nf")
	ast, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}

	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes
	if len(processes) != 1 {
		t.Fatalf("Expected 1 process, got %d", len(processes))
	}
}
