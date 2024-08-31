package nf

import (
	"path/filepath"
	"reft-go/nf/inputs"
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
	pinputs := processes[0].Inputs
	if len(pinputs) != 1 {
		t.Fatalf("Expected 1 input, got %d", len(pinputs))
	}
	each, ok := pinputs[0].(*inputs.Each)
	if !ok {
		t.Fatalf("Expected each input, got %v", pinputs[0])
	}
	_ = each
}
