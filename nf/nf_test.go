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
	stcVisitor := NewStcVisitor(cls)
	stcVisitor.VisitBlockStatement(ast.StatementBlock)
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
