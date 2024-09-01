package nf

import (
	"path/filepath"
	"reft-go/nf/inputs"
	"reft-go/nf/outputs"
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
	if len(pinputs) != 21 {
		t.Fatalf("Expected 21 inputs, got %d", len(pinputs))
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
	if path.Path != "$x.fa" {
		t.Fatalf("Expected path to be $x.fa, got %s", path.Path)
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
	env, ok := pinputs[16].(*inputs.Env)
	if !ok {
		t.Fatalf("Expected env input, got %v", pinputs[15])
	}
	if env.Var != "HELLO" {
		t.Fatalf("Expected var to be HELLO, got %s", env.Var)
	}
	stdin, ok := pinputs[18].(*inputs.Stdin)
	if !ok {
		t.Fatalf("Expected stdin input, got %v", pinputs[17])
	}
	if stdin.Var != "str" {
		t.Fatalf("Expected var to be str, got %s", stdin.Var)
	}
}

func TestProcessFileInputs(t *testing.T) {
	filePath := filepath.Join("./testdata", "process_file_inputs.nf")
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
	finput, ok := pinputs[0].(*inputs.File)
	if !ok {
		t.Fatalf("Expected file input, got %v", pinputs[0])
	}
	if finput.Path != "proteins" {
		t.Fatalf("Expected path to be proteins, got %s", finput.Path)
	}
}

func TestProcessValOutputs(t *testing.T) {
	filePath := filepath.Join("./testdata", "process_val_output.nf")
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
	poutputs := processes[0].Outputs
	if len(poutputs) != 3 {
		t.Fatalf("Expected 3 outputs, got %d", len(poutputs))
	}
	poutput := poutputs[0].(*outputs.Val)
	if poutput.Var != "x" {
		t.Fatalf("Expected var to be x, got %s", poutput.Var)
	}
	if poutput.Emit != "hello" {
		t.Fatalf("Expected emit to be hello, got %s", poutput.Emit)
	}
	if poutput.Topic != "report" {
		t.Fatalf("Expected topic to be report, got %s", poutput.Topic)
	}
	poutput = poutputs[1].(*outputs.Val)
	if poutput.Var != "BB11" {
		t.Fatalf("Expected var to be BB11, got %s", poutput.Var)
	}
	poutput = poutputs[2].(*outputs.Val)
	if poutput.Var != "${infile.baseName}.out" {
		t.Fatalf("Expected var to be ${infile.baseName}.out, got %s", poutput.Var)
	}
}

func TestProcessFileOutputs(t *testing.T) {
	filePath := filepath.Join("./testdata", "process_file_output.nf")
	_, err := parser.BuildAST(filePath)
	if err != nil {
		t.Fatalf("Failed to build AST: %v", err)
	}
}

func TestProcessPathOutputs(t *testing.T) {
	filePath := filepath.Join("./testdata", "process_path_output.nf")
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
	poutputs := processes[0].Outputs
	if len(poutputs) != 7 {
		t.Fatalf("Expected 7 outputs, got %d", len(poutputs))
	}
	poutput := poutputs[0].(*outputs.Path)

	if poutput.Path != "foo.txt" {
		t.Fatalf("Expected path to be foo.txt, got %s", poutput.Path)
	}

	poutput = poutputs[1].(*outputs.Path)
	if poutput.Path != "one.txt" {
		t.Fatalf("Expected path to be one.txt, got %s", poutput.Path)
	}
	if poutput.Arity != "1" {
		t.Fatalf("Expected arity to be 1, got %s", poutput.Arity)
	}

	// Test case for 'pair_*.txt'
	poutput = poutputs[2].(*outputs.Path)
	if poutput.Path != "pair_*.txt" {
		t.Errorf("Expected path to be pair_*.txt, got %s", poutput.Path)
	}
	if poutput.Arity != "2" {
		t.Errorf("Expected arity to be 2, got %s", poutput.Arity)
	}

	// Test case for 'many_*.txt'
	poutput = poutputs[3].(*outputs.Path)
	if poutput.Path != "many_*.txt" {
		t.Errorf("Expected path to be many_*.txt, got %s", poutput.Path)
	}
	if poutput.Arity != "1..*" {
		t.Errorf("Expected arity to be 1..*, got %s", poutput.Arity)
	}

	// Test case for "${species}.aln"
	poutput = poutputs[4].(*outputs.Path)
	if poutput.Path != "$species.aln" {
		t.Errorf("Expected path to be $species.aln, got %s", poutput.Path)
	}

	// Test case for 'blah.txt'
	poutput = poutputs[5].(*outputs.Path)
	if poutput.Path != "blah.txt" {
		t.Errorf("Expected path to be blah.txt, got %s", poutput.Path)
	}
	if poutput.FollowLinks {
		t.Errorf("Expected FollowLinks to be false, got true")
	}
	if poutput.Glob {
		t.Errorf("Expected Glob to be false, got true")
	}
	if !poutput.Hidden {
		t.Errorf("Expected Hidden to be true, got false")
	}

	// Test case for '*.txt'
	poutput = poutputs[6].(*outputs.Path)
	if poutput.Path != "*.txt" {
		t.Errorf("Expected path to be *.txt, got %s", poutput.Path)
	}
	if !poutput.IncludeInputs {
		t.Errorf("Expected IncludeInputs to be true, got false")
	}
	if poutput.MaxDepth != 1 {
		t.Errorf("Expected MaxDepth to be 1, got %d", poutput.MaxDepth)
	}
	if poutput.PathType != "dir" {
		t.Errorf("Expected PathType to be dir, got %s", poutput.PathType)
	}

}
