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
	if len(poutputs) != 2 {
		t.Fatalf("Expected 2 outputs, got %d", len(poutputs))
	}

	// Test the first file output
	foutput1, ok := poutputs[0].(*outputs.File)
	if !ok {
		t.Fatalf("Expected file output, got %T", poutputs[0])
	}
	if foutput1.Path != "proteins" {
		t.Errorf("Expected path to be 'proteins', got %s", foutput1.Path)
	}

	// Test the second file output
	foutput2, ok := poutputs[1].(*outputs.File)
	if !ok {
		t.Fatalf("Expected file output, got %T", poutputs[1])
	}
	if foutput2.Path != "foo.txt" {
		t.Errorf("Expected path to be 'foo.txt', got %s", foutput2.Path)
	}
	if foutput2.Emit != "hello" {
		t.Errorf("Expected emit to be 'hello', got %s", foutput2.Emit)
	}
	if !foutput2.Optional {
		t.Errorf("Expected optional to be true, got false")
	}
	if foutput2.Topic != "report" {
		t.Errorf("Expected topic to be 'report', got %s", foutput2.Topic)
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
	if len(poutputs) != 8 {
		t.Fatalf("Expected 8 outputs, got %d", len(poutputs))
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

	poutput = poutputs[7].(*outputs.Path)
	if poutput.Path != "foo.txt" {
		t.Errorf("Expected path to be foo.txt, got %s", poutput.Path)
	}
	if poutput.Emit != "hello" {
		t.Errorf("Expected emit to be hello, got %s", poutput.Emit)
	}
	if poutput.Topic != "report" {
		t.Errorf("Expected topic to be report, got %s", poutput.Topic)
	}
	if !poutput.Optional {
		t.Errorf("Expected optional to be true, got false")
	}
}

func TestProcessEnvOutputs(t *testing.T) {
	filePath := filepath.Join("./testdata", "process_env_output.nf")
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
	if len(poutputs) != 2 {
		t.Fatalf("Expected 2 outputs, got %d", len(poutputs))
	}

	// Test the first env output
	envOutput1, ok := poutputs[0].(*outputs.Env)
	if !ok {
		t.Fatalf("Expected env output, got %T", poutputs[0])
	}
	if envOutput1.Var != "VER" {
		t.Errorf("Expected var to be 'VER', got %s", envOutput1.Var)
	}
	if envOutput1.Emit != "tool_version" {
		t.Errorf("Expected emit to be 'tool_version', got %s", envOutput1.Emit)
	}

	// Test the second env output
	envOutput2, ok := poutputs[1].(*outputs.Env)
	if !ok {
		t.Fatalf("Expected env output, got %T", poutputs[1])
	}
	if envOutput2.Var != "DBVER" {
		t.Errorf("Expected var to be 'DBVER', got %s", envOutput2.Var)
	}
	if envOutput2.Emit != "db_version" {
		t.Errorf("Expected emit to be 'db_version', got %s", envOutput2.Emit)
	}
	if !envOutput2.Optional {
		t.Errorf("Expected optional to be true, got false")
	}
	if envOutput2.Topic != "report" {
		t.Errorf("Expected topic to be 'report', got %s", envOutput2.Topic)
	}
}

func TestProcessStdoutOutputs(t *testing.T) {
	filePath := filepath.Join("./testdata", "process_stdout_output.nf")
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
	if len(poutputs) != 1 {
		t.Fatalf("Expected 1 output, got %d", len(poutputs))
	}

	// Test the stdout output
	stdoutOutput, ok := poutputs[0].(*outputs.Stdout)
	if !ok {
		t.Fatalf("Expected stdout output, got %T", poutputs[0])
	}

	// Check Emit
	if stdoutOutput.Emit != "hello" {
		t.Errorf("Expected Emit to be 'hello', got %s", stdoutOutput.Emit)
	}

	// Check Optional
	if !stdoutOutput.Optional {
		t.Errorf("Expected Optional to be true, got false")
	}

	// Check Topic
	if stdoutOutput.Topic != "report" {
		t.Errorf("Expected Topic to be 'report', got %s", stdoutOutput.Topic)
	}
}

func TestProcessEvalOutputs(t *testing.T) {
	filePath := filepath.Join("./testdata", "process_eval_output.nf")
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
	if len(poutputs) != 2 {
		t.Fatalf("Expected 2 outputs, got %d", len(poutputs))
	}

	// Test the first eval output
	evalOutput1, ok := poutputs[0].(*outputs.Eval)
	if !ok {
		t.Fatalf("Expected eval output, got %T", poutputs[0])
	}
	if evalOutput1.Command != "bash --version" {
		t.Errorf("Expected command to be 'bash --version', got %s", evalOutput1.Command)
	}

	// Test the second eval output
	evalOutput2, ok := poutputs[1].(*outputs.Eval)
	if !ok {
		t.Fatalf("Expected eval output, got %T", poutputs[1])
	}
	if evalOutput2.Command != "echo \"foo\"" {
		t.Errorf("Expected command to be 'echo \"foo\"', got %s", evalOutput2.Command)
	}
	if evalOutput2.Emit != "hello" {
		t.Errorf("Expected emit to be 'hello', got %s", evalOutput2.Emit)
	}
	if !evalOutput2.Optional {
		t.Errorf("Expected optional to be true, got false")
	}
	if evalOutput2.Topic != "report" {
		t.Errorf("Expected topic to be 'report', got %s", evalOutput2.Topic)
	}
}

func TestProcessTupleOutputs(t *testing.T) {
	filePath := filepath.Join("./testdata", "process_tuple_output.nf")
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
	if len(poutputs) != 2 {
		t.Fatalf("Expected 2 outputs, got %d", len(poutputs))
	}

	// Test the first tuple output
	tupleOutput1, ok := poutputs[0].(*outputs.Tuple)
	if !ok {
		t.Fatalf("Expected tuple output, got %T", poutputs[0])
	}
	if len(tupleOutput1.Values) != 2 {
		t.Errorf("Expected 2 values in tuple, got %d", len(tupleOutput1.Values))
	}
	val1, ok := tupleOutput1.Values[0].(*outputs.Val)
	if !ok {
		t.Errorf("Expected first value to be Val, got %T", tupleOutput1.Values[0])
	}
	if val1.Var != "meta" {
		t.Errorf("Expected Val var to be 'meta', got %s", val1.Var)
	}
	path1, ok := tupleOutput1.Values[1].(*outputs.Path)
	if !ok {
		t.Errorf("Expected second value to be Path, got %T", tupleOutput1.Values[1])
	}
	if path1.Path != "$prefix.tsv" {
		t.Errorf("Expected Path path to be '$prefix.tsv', got %s", path1.Path)
	}
	if tupleOutput1.Emit != "report" {
		t.Errorf("Expected emit to be 'report', got %s", tupleOutput1.Emit)
	}
	if tupleOutput1.Topic != "report" {
		t.Errorf("Expected topic to be 'report', got %s", tupleOutput1.Topic)
	}

	// Test the second tuple output
	tupleOutput2, ok := poutputs[1].(*outputs.Tuple)
	if !ok {
		t.Fatalf("Expected tuple output, got %T", poutputs[1])
	}
	if len(tupleOutput2.Values) != 2 {
		t.Errorf("Expected 2 values in tuple, got %d", len(tupleOutput2.Values))
	}
	val2, ok := tupleOutput2.Values[0].(*outputs.Val)
	if !ok {
		t.Errorf("Expected first value to be Val, got %T", tupleOutput2.Values[0])
	}
	if val2.Var != "meta" {
		t.Errorf("Expected Val var to be 'meta', got %s", val2.Var)
	}
	path2, ok := tupleOutput2.Values[1].(*outputs.Path)
	if !ok {
		t.Errorf("Expected second value to be Path, got %T", tupleOutput2.Values[1])
	}
	if path2.Path != "$prefix-mutations.tsv" {
		t.Errorf("Expected Path path to be '$prefix-mutations.tsv', got %s", path2.Path)
	}
	if tupleOutput2.Emit != "mutation_report" {
		t.Errorf("Expected emit to be 'mutation_report', got %s", tupleOutput2.Emit)
	}
	if !tupleOutput2.Optional {
		t.Errorf("Expected optional to be true, got false")
	}
}
