package parser

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/antlr4-go/antlr/v4"
)

func TestGroovyParserGStringFile(t *testing.T) {
	filePath := filepath.Join("testdata", "gstring.groovy")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	parser := NewGroovyParser(stream)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("parser.CompilationUnit() panicked: %v", r)
		}
	}()

	// Parse the file
	parser.CompilationUnit()
}

func TestGroovyParserUtils(t *testing.T) {
	filePath := filepath.Join("testdata", "utils_nfcore_pipeline.nf")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	parser := NewGroovyParser(stream)
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("parser.CompilationUnit() panicked: %v", r)
		}
	}()

	// Parse the file
	parser.CompilationUnit()
}

func TestGroovyParserExpression(t *testing.T) {
	filePath := filepath.Join("testdata", "expression", "01.groovy")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	parser := NewGroovyParser(stream)
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("parser.CompilationUnit() panicked: %v", r)
		}
	}()

	// Parse the file
	parser.CompilationUnit()
}

func TestGroovyParserCommandExpr(t *testing.T) {
	filePath := filepath.Join("testdata", "cnvkit_batch_main.nf")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	parser := NewGroovyParser(stream)
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

	/*
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("parser.CompilationUnit() panicked: %v", r)
			}
		}()
	*/

	// Parse the file
	tree := parser.CompilationUnit()
	fmt.Println(tree)
}

func TestInclude(t *testing.T) {
	filePath := filepath.Join("testdata", "include.nf")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	stream.Fill()
	parser := NewGroovyParser(stream)
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

	/*
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("parser.CompilationUnit() panicked: %v", r)
			}
		}()
	*/

	// Parse the file
	tree := parser.CompilationUnit()
	fmt.Println(tree)
	builder := NewASTBuilder(filePath)
	foo := builder.Visit(tree)
	fmt.Println(foo)
}

func TestInclude2(t *testing.T) {
	filePath := filepath.Join("testdata", "include.nf")
	input, err := antlr.NewFileStream(filePath)
	if err != nil {
		t.Fatalf("Failed to open file %s: %s", filePath, err)
	}

	lexer := NewGroovyLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, 0)
	//tokens := lexer.GetAllTokens()
	//tokenStream := NewPreloadedTokenStream(tokens, lexer)
	stream.Fill()
	parser := NewGroovyParser(stream)
	parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeLLExactAmbigDetection)

	/*
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("parser.CompilationUnit() panicked: %v", r)
			}
		}()
	*/

	// Parse the file
	tree := parser.CompilationUnit()
	//fmt.Println(tree)
	//stmt := tree.ScriptStatements().ScriptStatement(0)
	//ctx := stmt.(*ScriptStatementContext).GetChild(0).(*ExpressionStmtAltContext).StatementContext.GetChild(0).(*CommandExprAltContext).CommandExpression()
	//cmdArgCtx := ctx.GetChild(1).(*CommandArgumentContext)
	//_ = cmdArgCtx
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(tree).(*ModuleNode)
	bs := ast.StatementBlock
	
	fmt.Println(ast.StatementBlock.GetText())
}
