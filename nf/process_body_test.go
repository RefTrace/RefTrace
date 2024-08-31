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
	if len(pinputs) != 20 {
		t.Fatalf("Expected 20 inputs, got %d", len(pinputs))
	}
	each, ok := pinputs[0].(*inputs.Each)
	if !ok {
		t.Fatalf("Expected each input, got %v", pinputs[0])
	}
	collection, ok := each.Collection.(*inputs.Val)
	if !ok {
		t.Fatalf("Expected val input, got %v", each.Collection)
	}
	if collection.Var != "mode" {
		t.Fatalf("Expected var to be mode, got %s", collection.Var)
	}
	each, ok = pinputs[1].(*inputs.Each)
	if !ok {
		t.Fatalf("Expected each input, got %v", pinputs[1])
	}
	collection2, ok := each.Collection.(*inputs.Path)
	if !ok {
		t.Fatalf("Expected path input, got %v", each.Collection)
	}
	if collection2.Path != "lib" {
		t.Fatalf("Expected path to be lib, got %s", collection2.Path)
	}
	path, ok := pinputs[2].(*inputs.Path)
	if !ok {
		t.Fatalf("Expected path input, got %v", pinputs[2])
	}
	// TODO: check gstring parsing
	if path.Path != "${x}fa\"" {
		t.Fatalf("Expected path to be ${x}.fa, got %s", path.Path)
	}
	path, ok = pinputs[3].(*inputs.Path)
	if !ok {
		t.Fatalf("Expected path input, got %v", pinputs[3])
	}
	if path.Path != "x" {
		t.Fatalf("Expected path to be x, got %s", path.Path)
	}
	if path.StageAs != "data.txt" {
		t.Fatalf("Expected stageAs to be data.txt, got %s", path.StageAs)
	}
	path, ok = pinputs[4].(*inputs.Path)
	if !ok {
		t.Fatalf("Expected path input, got %v", pinputs[4])
	}
	if path.Path != "query.fa" {
		t.Fatalf("Expected path to be query.fa, got %s", path.Path)
	}
	tuple, ok := pinputs[5].(*inputs.Tuple)
	if !ok {
		t.Fatalf("Expected tuple input, got %v", pinputs[5])
	}
	val, ok := tuple.Values[0].(*inputs.Val)
	if !ok {
		t.Fatalf("Expected val input, got %v", tuple.Values[0])
	}
	if val.Var != "x" {
		t.Fatalf("Expected var to be x, got %s", val.Var)
	}
	path, ok = tuple.Values[1].(*inputs.Path)
	if !ok {
		t.Fatalf("Expected path input, got %v", tuple.Values[1])
	}
	if path.Path != "input.txt" {
		t.Fatalf("Expected path to be input.txt, got %s", path.Path)
	}
	tuple, ok = pinputs[6].(*inputs.Tuple)
	if !ok {
		t.Fatalf("Expected tuple input, got %v", pinputs[6])
	}
	val, ok = tuple.Values[0].(*inputs.Val)
	if !ok {
		t.Fatalf("Expected val input, got %v", tuple.Values[0])
	}
	if val.Var != "x" {
		t.Fatalf("Expected var to be x, got %s", val.Var)
	}
	path, ok = tuple.Values[1].(*inputs.Path)
	if !ok {
		t.Fatalf("Expected path input, got %v", tuple.Values[1])
	}
	if path.Path != "many_*.txt" {
		t.Fatalf("Expected path to be many_*.txt, got %s", path.Path)
	}
	if path.Arity != "1..*" {
		t.Fatalf("Expected arity to be 1..*, got %s", path.Arity)
	}
	env, ok := pinputs[15].(*inputs.Env)
	if !ok {
		t.Fatalf("Expected env input, got %v", pinputs[15])
	}
	if env.Var != "HELLO" {
		t.Fatalf("Expected var to be HELLO, got %s", env.Var)
	}
	stdin, ok := pinputs[17].(*inputs.Stdin)
	if !ok {
		t.Fatalf("Expected stdin input, got %v", pinputs[17])
	}
	if stdin.Var != "str" {
		t.Fatalf("Expected var to be str, got %s", stdin.Var)
	}
}
