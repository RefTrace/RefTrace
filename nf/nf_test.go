package nf

import (
	"path/filepath"
	"reft-go/parser"
	"runtime/debug"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

func TestSarekEntireMain(t *testing.T) {
	debug.SetGCPercent(-1)
	filePath := filepath.Join("../parser/testdata", "sarek_entire_main.nf")
	result, err := parser.BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	builder := parser.NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*parser.ModuleNode)
	//cls := ast.GetClasses()[0]
	paramVisitor := NewParamVisitor()
	paramVisitor.VisitBlockStatement(ast.StatementBlock)
	params := paramVisitor.GetSortedParams()
	if len(params) != 65 {
		t.Fatalf("Expected 65 params, got %d", len(params))
	}
	includeVisitor := NewIncludeVisitor()
	includeVisitor.VisitBlockStatement(ast.StatementBlock)
	includes := includeVisitor.GetSortedIncludes()
	if len(includes) != 8 {
		t.Fatalf("Expected 8 includes, got %d", len(includes))
	}
	workflowVisitor := NewWorkflowVisitor()
	workflowVisitor.VisitBlockStatement(ast.StatementBlock)
	workflows := workflowVisitor.workflows
	if len(workflows) != 1 {
		t.Fatalf("Expected 1 workflow, got %d", len(workflows))
	}
	//stcVisitor := NewStcVisitor(cls)
	//stcVisitor.VisitBlockStatement(ast.StatementBlock)
}

func TestSimpleWorkflow(t *testing.T) {
	debug.SetGCPercent(-1)
	filePath := filepath.Join("./testdata", "simple_workflow.nf")
	result, err := parser.BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	builder := parser.NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*parser.ModuleNode)

	workflowVisitor := NewWorkflowVisitor()
	workflowVisitor.VisitBlockStatement(ast.StatementBlock)
	workflows := workflowVisitor.workflows
	if len(workflows) != 2 {
		t.Fatalf("Expected 2 workflows, got %d", len(workflows))
	}
	workflow := workflows[0]
	if len(workflow.Takes) != 2 {
		t.Fatalf("Expected 2 takes, got %d", len(workflow.Takes))
	}
	if len(workflow.Emits) != 2 {
		t.Fatalf("Expected 2 emits, got %d", len(workflow.Emits))
	}
	//stcVisitor := NewStcVisitor(cls)
	//stcVisitor.VisitBlockStatement(ast.StatementBlock)
}

func TestSimpleProcess(t *testing.T) {
	debug.SetGCPercent(-1)
	filePath := filepath.Join("./testdata", "simple_process.nf")
	result, err := parser.BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	builder := parser.NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*parser.ModuleNode)

	processVisitor := NewProcessVisitor()
	processVisitor.VisitBlockStatement(ast.StatementBlock)
	processes := processVisitor.processes
	if len(processes) != 1 {
		t.Fatalf("Expected 1 process, got %d", len(processes))
	}
}

func TestChannelFromPath(t *testing.T) {
	debug.SetGCPercent(-1)
	filePath := filepath.Join("./testdata", "channel_frompath.nf")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := parser.NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	stream.Fill()
	groovyParser := parser.NewGroovyParser(stream)
	//parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

	// Parse the file
	tree := groovyParser.CompilationUnit()
	builder := parser.NewASTBuilder(filePath)
	ast := builder.Visit(tree).(*parser.ModuleNode)
	cls := ast.GetClasses()[0]
	stcVisitor := NewStcVisitor(cls)
	stcVisitor.VisitBlockStatement(ast.StatementBlock)
}
