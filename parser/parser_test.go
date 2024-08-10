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
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
	stmt := bs.statements[0]
	exprStmt, ok := stmt.(*ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement, but got %T", stmt)
	}
	mce, ok := exprStmt.GetExpression().(*MethodCallExpression)
	if !ok {
		t.Fatalf("Expected expression to be an MethodCallExpression, but got %T", exprStmt.GetExpression())
	}
	ce, ok := mce.GetMethod().(*ConstantExpression)
	if !ok {
		t.Fatalf("Expected expression to be an ConstantExpression, but got %T", mce.GetMethod())
	}
	if ce.GetText() != "from" {
		t.Errorf("Expected 'from', but got '%s'", ce.GetText())
	}
	oe, ok := mce.GetObjectExpression().(*MethodCallExpression)
	if !ok {
		t.Fatalf("Expected expression to be an MethodCallExpression, but got %T", mce.GetObjectExpression())
	}
	method, ok := oe.GetMethod().(*ConstantExpression)
	if !ok {
		t.Fatalf("Expected expression to be an ConstantExpression, but got %T", mce.GetMethod())
	}
	if method.GetText() != "include" {
		t.Errorf("Expected 'include', but got '%s'", method.GetText())
	}
	args, ok := oe.Arguments.(*TupleExpression)
	if !ok {
		t.Fatalf("Expected tuple expression to be an TupleExpression, but got %T", oe.Arguments)
	}
	if len(args.GetExpressions()) != 1 {
		t.Errorf("Expected exactly 1 argument in the tuple expression, but got %d", len(args.GetExpressions()))
	}
	argumentListExpr, ok := args.GetExpressions()[0].(*ArgumentListExpression)
	if !ok {
		t.Fatalf("Expected argument expression to be an ArgumentListExpression, but got %T", args.GetExpressions())
	}
	if len(argumentListExpr.GetExpressions()) != 1 {
		t.Errorf("Expected exactly 1 argument in the tuple expression, but got %d", len(args.GetExpressions()))
	}
	closureExpr, ok := argumentListExpr.GetExpression(0).(*ClosureExpression)
	if !ok {
		t.Fatalf("Expected closure expression to be an ClosureExpression, but got %T", args.GetExpressions())
	}
	codeStmt, ok := closureExpr.GetCode().(*BlockStatement)
	if !ok {
		t.Fatalf("Expected closure expression to be a BlockStatement, but got %T", closureExpr.GetCode())
	}
	if len(codeStmt.GetStatements()) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(codeStmt.GetStatements()))
	}
	exprStmt, ok = codeStmt.GetStatements()[0].(*ExpressionStatement)
	if !ok {
		t.Fatalf("Expected expression to be an ExpressionStatement, but got %T", codeStmt.GetStatements())
	}
	varExpr, ok := exprStmt.GetExpression().(*VariableExpression)
	if !ok {
		t.Fatalf("Expected variable expression to be an VariableExpression, but got %T", exprStmt.GetExpression())
	}
	if varExpr.GetText() != "SAREK" {
		t.Errorf("Expected 'SAREK', but got '%s'", varExpr.GetText())
	}
}

func TestParams(t *testing.T) {
	filePath := filepath.Join("testdata", "params.nf")
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
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
	stmt := bs.statements[0]
	exprStmt, ok := stmt.(*ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement, but got %T", stmt)
	}
	binaryExpr, ok := exprStmt.GetExpression().(*BinaryExpression)
	if !ok {
		t.Fatalf("Expected expression to be an BinaryExpression, but got %T", exprStmt.GetExpression())
	}
	left, ok := binaryExpr.leftExpression.(*PropertyExpression)
	if !ok {
		t.Fatalf("Expected expression to be an PropertyExpression, but got %T", binaryExpr.leftExpression)
	}
	object, ok := left.objectExpression.(*VariableExpression)
	if !ok {
		t.Fatalf("Expected expression to be an VariableExpression, but got %T", left.objectExpression)
	}
	if object.GetText() != "params" {
		t.Fatalf("Expected 'params', but got '%s'", object.GetText())
	}
	property, ok := left.property.(*ConstantExpression)
	if !ok {
		t.Fatalf("Expected expression to be an ConstantExpression, but got %T", left.property)
	}
	if property.GetText() != "ascat_alleles" {
		t.Fatalf("Expected 'ascat_alleles', but got '%s'", property.GetText())
	}
	right, ok := binaryExpr.rightExpression.(*MethodCallExpression)
	if !ok {
		t.Fatalf("Expected expression to be an MethodCallExpression, but got %T", binaryExpr.rightExpression)
	}
	object, ok = right.GetObjectExpression().(*VariableExpression)
	if !ok {
		t.Fatalf("Expected expression to be an VariableExpression, but got %T", left.objectExpression)
	}
	if object.GetText() != "this" {
		t.Fatalf("Expected 'this', but got '%s'", object.GetText())
	}
	method, ok := right.GetMethod().(*ConstantExpression)
	if !ok {
		t.Fatalf("Expected expression to be an ConstantExpression, but got %T", right.GetMethod())
	}
	if method.GetText() != "getGenomeAttribute" {
		t.Fatalf("Expected 'getGenomeAttribute', but got '%s'", method.GetText())
	}
}

func TestElvis(t *testing.T) {
	filePath := filepath.Join("testdata", "elvis.nf")
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

	// Parse the file
	tree := parser.CompilationUnit()
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
	stmt := bs.statements[0]
	exprStmt, ok := stmt.(*ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement, but got %T", stmt)
	}
	binaryExpr, ok := exprStmt.GetExpression().(*BinaryExpression)
	if !ok {
		t.Fatalf("Expected expression to be an BinaryExpression, but got %T", exprStmt.GetExpression())
	}
	left, ok := binaryExpr.leftExpression.(*VariableExpression)
	if !ok {
		t.Fatalf("Expected expression to be an PropertyExpression, but got %T", binaryExpr.leftExpression)
	}
	if left.variable != "ascat_genome" {
		t.Fatalf("Expected 'ascat_genome', but got '%s'", left.variable)
	}
	right, ok := binaryExpr.rightExpression.(*ElvisOperatorExpression)
	if !ok {
		t.Fatalf("Expected expression to be an ElvisOperatorExpression, but got %T", binaryExpr.rightExpression)
	}
	truth, ok := right.truthExpression.(*PropertyExpression)
	if !ok {
		t.Fatalf("Expected expression to be a PropertyExpression, but got %T", right.truthExpression)
	}
	truthVar, ok := truth.objectExpression.(*VariableExpression)
	if !ok {
		t.Fatalf("Expected expression to be a VariableExpression, but got %T", truth.objectExpression)
	}
	if truthVar.variable != "params" {
		t.Fatalf("Expected 'params', but got '%s'", truthVar.GetText())
	}
	constExpr := truth.property.(*ConstantExpression)
	if constExpr.GetText() != "ascat_genome" {
		t.Fatalf("Expected 'ascat_genome', but got '%s'", constExpr.GetText())
	}
}

func TestTernary(t *testing.T) {
	filePath := filepath.Join("testdata", "ternary.nf")
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

	// Parse the file
	tree := parser.CompilationUnit()
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
	stmt := bs.statements[0]
	exprStmt, ok := stmt.(*ExpressionStatement)
	if !ok {
		t.Fatalf("Expected statement to be an ExpressionStatement, but got %T", stmt)
	}
	binaryExpr, ok := exprStmt.GetExpression().(*BinaryExpression)
	if !ok {
		t.Fatalf("Expected expression to be an BinaryExpression, but got %T", exprStmt.GetExpression())
	}
	left, ok := binaryExpr.leftExpression.(*VariableExpression)
	if !ok {
		t.Fatalf("Expected expression to be an PropertyExpression, but got %T", binaryExpr.leftExpression)
	}
	if left.variable != "dbsnp_vqsr" {
		t.Fatalf("Expected 'dbsnp_vqsr', but got '%s'", left.variable)
	}
	right, ok := binaryExpr.rightExpression.(*TernaryExpression)
	if !ok {
		t.Fatalf("Expected expression to be an TernaryExpression, but got %T", binaryExpr.rightExpression)
	}
	truth, ok := right.truthExpression.(*MethodCallExpression)
	if !ok {
		t.Fatalf("Expected expression to be a PropertyExpression, but got %T", right.truthExpression)
	}
	truthConst, ok := truth.Method.(*ConstantExpression)
	if !ok {
		t.Fatalf("Expected expression to be a VariableExpression, but got %T", truth.Method)
	}
	if truthConst.GetText() != "value" {
		t.Fatalf("Expected 'value', but got '%s'", truthConst.GetText())
	}
	objExpr := truth.ObjectExpression.(*VariableExpression)
	if objExpr.variable != "Channel" {
		t.Fatalf("Expected 'Channel', but got '%s'", objExpr.variable)
	}
	args := truth.GetArguments().(*TupleExpression).expressions[0].(*ArgumentListExpression)
	if len(args.expressions) != 1 {
		t.Fatalf("Expected exactly 1 argument in the TupleExpression, but got %d", len(args.expressions))
	}
	property := args.expressions[0].(*PropertyExpression)
	obj := property.GetObjectExpression().(*VariableExpression)
	prop := property.GetProperty().(*ConstantExpression)
	if obj.variable != "params" {
		t.Fatalf("Expected 'params', but got '%s'", obj.variable)
	}
	if prop.GetText() != "dbsnp_vqsr" {
		t.Fatalf("Expected 'dbsnp_vqsr', but got '%s'", prop.GetText())
	}
}
