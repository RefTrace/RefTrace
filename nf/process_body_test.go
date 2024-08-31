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
	inputs := processes[0].Inputs
	if len(inputs) != 1 {
		t.Fatalf("Expected 1 input, got %d", len(inputs))
	}
	if each, ok := inputs[0].(*inputs.Each); !ok {
		t.Fatalf("Expected each input, got %v", inputs[0])
	}
}
