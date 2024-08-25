package parser

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"testing"
)

func TestGroovyParserGStringFile(t *testing.T) {
	filePath := filepath.Join("testdata", "gstring.groovy")
	_, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
}

func TestGroovyParserUtils(t *testing.T) {
	filePath := filepath.Join("testdata", "utils_nfcore_pipeline.nf")
	_, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
}

func TestGroovyParserExpression(t *testing.T) {
	filePath := filepath.Join("testdata", "expression", "01.groovy")
	_, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
}

func TestGroovyParserCommandExpr(t *testing.T) {
	filePath := filepath.Join("testdata", "cnvkit_batch_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "LL" {
		t.Fatalf("Expected parsing mode to be LL, but got %s", result.Mode)
	}

	if result.Tree == nil {
		t.Fatal("Parse tree is nil")
	}
}

func TestInclude(t *testing.T) {
	filePath := filepath.Join("testdata", "include.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
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
	args, ok := oe.Arguments.(*ArgumentListExpression)
	closureExpr, ok := args.GetExpression(0).(*ClosureExpression)
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
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
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
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
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
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
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
	args := truth.GetArguments().(*ArgumentListExpression)
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

func TestTernaryClosure(t *testing.T) {
	filePath := filepath.Join("testdata", "ternary_closure.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
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
	if left.variable != "fasta" {
		t.Fatalf("Expected 'fasta', but got '%s'", left.variable)
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
	if truthConst.GetText() != "collect" {
		t.Fatalf("Expected 'collect', but got '%s'", truthConst.GetText())
	}
	mce := truth.ObjectExpression.(*MethodCallExpression).ObjectExpression.(*MethodCallExpression)
	objExpr := mce.ObjectExpression.(*VariableExpression)
	if objExpr.variable != "Channel" {
		t.Fatalf("Expected 'Channel', but got '%s'", objExpr.variable)
	}
	constExpr := mce.Method.(*ConstantExpression)
	if constExpr.GetText() != "fromPath" {
		t.Fatalf("Expected 'fromPath', but got '%s'", constExpr.GetText())
	}
	// closure
	args := truth.GetObjectExpression().(*MethodCallExpression).GetArguments().(*ArgumentListExpression).GetExpressions()[0].(*ClosureExpression)
	params := args.GetParameters()
	param := params[0]
	if param.GetName() != "it" {
		t.Fatalf("Expected 'it', but got '%s'", param.GetName())
	}
	exprs := args.GetCode().(*BlockStatement).GetStatements()[0].(*ExpressionStatement).GetExpression().(*ListExpression).GetExpressions()
	mapExpr := exprs[0].(*MapExpression)
	varExpr := exprs[1].(*VariableExpression)
	if varExpr.variable != "it" {
		t.Fatalf("Expected 'it', but got '%s'", varExpr.variable)
	}
	mapEntry := mapExpr.GetMapEntryExpressions()[0]
	if mapEntry.GetKeyExpression().(*ConstantExpression).GetText() != "id" {
		t.Fatalf("Expected 'id', but got '%s'", mapEntry.GetKeyExpression().GetText())
	}
	propExpr := mapEntry.GetValueExpression().(*PropertyExpression)
	if propExpr.GetProperty().(*ConstantExpression).GetText() != "baseName" {
		t.Fatalf("Expected 'baseName', but got '%s'", propExpr.GetProperty().GetText())
	}
	if propExpr.GetObjectExpression().(*VariableExpression).variable != "it" {
		t.Fatalf("Expected 'it', but got '%s'", propExpr.GetObjectExpression().GetText())
	}
}

func TestTopLevelIf(t *testing.T) {
	filePath := filepath.Join("testdata", "top_level_if.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
	ifStmt := bs.GetStatements()[0].(*IfStatement)
	op := ifStmt.GetBooleanExpression().GetExpression().(*BinaryExpression).GetOperation()
	if op.GetText() != "&&" {
		t.Errorf("Expected '&&', but got '%s'", op.GetText())
	}
}

func TestSimpleWorkflow(t *testing.T) {
	debug.SetGCPercent(-1)
	filePath := filepath.Join("testdata", "simple_workflow.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
	stmt := bs.GetStatements()[0].(*ExpressionStatement)
	closure := stmt.GetExpression().(*MethodCallExpression).GetArguments().(*ArgumentListExpression).GetExpressions()[0].(*ClosureExpression)
	stmts := closure.GetCode().(*BlockStatement).GetStatements()
	if len(stmts) != 3 {
		t.Errorf("Expected exactly 3 statement in the block, but got %d", len(stmts))
	}
	mainStmt := stmts[0].(*ExpressionStatement)
	if mainStmt.GetStatementLabel() != "main" {
		t.Errorf("Expected 'main', but got '%s'", mainStmt.GetStatementLabel())
	}
}

func TestFunction(t *testing.T) {
	filePath := filepath.Join("testdata", "function.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	methods := ast.Methods
	if len(methods) != 1 {
		t.Errorf("Expected exactly 1 method in the block, but got %d", len(methods))
	}
}

func TestSarekMain(t *testing.T) {
	filePath := filepath.Join("testdata", "sarek_main_workflow.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
	workflow := bs.statements[0].(*ExpressionStatement).GetExpression().(*MethodCallExpression)
	workflowName := workflow.Method.(*ConstantExpression).GetText()
	if workflowName != "workflow" {
		t.Errorf("Expected 'workflow', but got '%s'", workflowName)
	}
	mce := workflow.GetArguments().(*ArgumentListExpression).GetExpressions()[0].(*MethodCallExpression)
	if mce.Method.(*ConstantExpression).GetText() != "NFCORE_SAREK" {
		t.Errorf("Expected 'NFCORE_SAREK', but got '%s'", mce.Method.(*ConstantExpression).GetText())
	}
	closure := mce.GetArguments().(*ArgumentListExpression).GetExpressions()[0].(*ClosureExpression)
	stmts := closure.GetCode().(*BlockStatement).GetStatements()
	if len(stmts) != 41 {
		t.Errorf("Expected exactly 41 statements in the block, but got %d", len(stmts))
	}
}

func TestSarekMain2(t *testing.T) {
	filePath := filepath.Join("testdata", "sarek_main_workflow2.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 37 {
		t.Errorf("Expected exactly 37 statements in the block, but got %d", len(bs.statements))
	}
}

func TestCnvKitBatchMain2(t *testing.T) {
	filePath := filepath.Join("testdata", "cnvkit_batch_main2.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "LL" {
		t.Fatalf("Expected parsing mode to be LL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
}

func TestVarcalMain(t *testing.T) {
	filePath := filepath.Join("testdata", "varcal_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "LL" {
		t.Fatalf("Expected parsing mode to be LL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
}

func TestPrepareIntervalsMain(t *testing.T) {
	filePath := filepath.Join("testdata", "prepare_intervals_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 6 {
		t.Errorf("Expected exactly 6 statements in the block, but got %d", len(bs.statements))
	}
}

func TestSamplesheetToChannelMain(t *testing.T) {
	filePath := filepath.Join("testdata", "samplesheet_to_channel_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
}

func TestUtilsNFcoreSarekPipelineMain(t *testing.T) {
	filePath := filepath.Join("testdata", "utils_nfcore_sarek_pipeline_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 15 {
		t.Errorf("Expected exactly 15 statements in the block, but got %d", len(bs.statements))
	}
}

func TestUtilsNFPipelineMain(t *testing.T) {
	filePath := filepath.Join("testdata", "utils_nextflow_pipeline_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
}

func TestUtilsNFCorePipelineMain(t *testing.T) {
	filePath := filepath.Join("testdata", "utils_nfcore_pipeline_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
}

func TestEagerMain(t *testing.T) {
	filePath := filepath.Join("testdata", "eager_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "LL" {
		t.Fatalf("Expected parsing mode to be LL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 183 {
		t.Errorf("Expected exactly 183 statements in the block, but got %d", len(bs.statements))
	}
}

func TestPathInProcess(t *testing.T) {
	filePath := filepath.Join("testdata", "path_in_process.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		t.Errorf("Expected exactly 1 statement in the block, but got %d", len(bs.statements))
	}
}

func TestSarekEntireMain(t *testing.T) {
	filePath := filepath.Join("testdata", "sarek_entire_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "SLL" {
		t.Fatalf("Expected parsing mode to be SLL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	bs := ast.StatementBlock
	if len(bs.statements) != 71 {
		t.Errorf("Expected exactly 71 statements in the block, but got %d", len(bs.statements))
	}
	_, ok := bs.statements[67].(*IfStatement)
	if !ok {
		t.Errorf("Expected statement to be an IfStatement, but got %T", bs.statements[67])
	}
	_, ok = bs.statements[68].(*IfStatement)
	if !ok {
		t.Errorf("Expected statement to be an IfStatement, but got %T", bs.statements[68])
	}
}

func TestDeepVariantMain(t *testing.T) {
	filePath := filepath.Join("testdata", "deepvariant_main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "LL" {
		t.Fatalf("Expected parsing mode to be LL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	_ = ast
}

func TestCreateIntervalsBed(t *testing.T) {
	filePath := filepath.Join("testdata", "sarek/modules/local/create_intervals_bed", "main.nf")
	result, err := BuildCST(filePath)
	if err != nil {
		t.Fatalf("Failed to build CST: %v", err)
	}
	if result.Mode != "LL" {
		t.Fatalf("Expected parsing mode to be LL, but got %s", result.Mode)
	}
	builder := NewASTBuilder(filePath)
	ast := builder.Visit(result.Tree).(*ModuleNode)
	_ = ast
}

/*
func getStatementsCount(ast *ModuleNode) int {
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		return len(bs.statements)
	}

	// Try to convert the single statement to a MethodCallExpression
	methodCall, ok := bs.statements[0].(*MethodCallExpression)
	if !ok {
		return 1
	}

	// Check if there are any arguments
	if len(methodCall.Arguments) == 0 {
		return 1
	}

	// Try to convert the first argument to a Closure
	closure, ok := methodCall.Arguments[0].(*Closure)
	if !ok {
		return 1
	}

	// Return the number of statements in the closure's code field
	return len(closure.Code.statements)
}
*/

func getStatementsCount(ast *ModuleNode) int {
	bs := ast.StatementBlock
	if len(bs.statements) != 1 {
		return len(bs.statements)
	}
	exprStmt, ok := bs.statements[0].(*ExpressionStatement)
	if !ok {
		return 1
	}
	expr, ok := exprStmt.GetExpression().(*MethodCallExpression)
	if !ok {
		return 1
	}
	args := expr.GetArguments().(*ArgumentListExpression)
	if len(args.GetExpressions()) != 1 {
		return 1
	}
	firstExpr := args.GetExpressions()[0]
	mce, ok := firstExpr.(*MethodCallExpression)
	if !ok {
		return 1
	}
	args = mce.GetArguments().(*ArgumentListExpression)
	if len(args.GetExpressions()) != 1 {
		return 1
	}
	firstExpr = args.GetExpressions()[0]
	closure, ok := firstExpr.(*ClosureExpression)
	if !ok {
		return 1
	}
	code := closure.GetCode().(*BlockStatement)
	stmts := code.GetStatements()
	return len(stmts)
}

/*
func TestParseAllSarekFiles(t *testing.T) {
	debug.SetGCPercent(-1)
	dir := filepath.Join("testdata", "sarek")

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".nf" {
			t.Run(path, func(t *testing.T) {
				input, err := antlr.NewFileStream(path)
				if err != nil {
					t.Fatalf("Failed to open file %s: %s", path, err)
				}

				lexer := NewGroovyLexer(input)
				stream := antlr.NewCommonTokenStream(lexer, 0)
				stream.Fill()
				parser := NewGroovyParser(stream)
				//parser.GetInterpreter().SetPredictionMode(antlr.PredictionModeSLL)

				defer func() {
					if r := recover(); r != nil {
						t.Fatalf("Parser panicked while processing file %s: %v", path, r)
					}
				}()

				tree := parser.CompilationUnit()
				builder := NewASTBuilder(path)
				ast := builder.Visit(tree).(*ModuleNode)
				numStatements := getStatementsCount(ast)
				t.Logf("File: %s, Number of statements: %d", path, numStatements)

				if ast == nil {
					t.Errorf("Failed to parse file: %s", path)
				}
			})
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Error walking the path %s: %v", dir, err)
	}
}
*/

func TestSarekGoldMapping(t *testing.T) {
	dir := filepath.Join("testdata", "sarek")

	// Load the gold mapping
	goldFile, err := os.Open(filepath.Join("testdata", "sarek_gold_mapping.json"))
	if err != nil {
		t.Fatalf("Failed to open gold mapping file: %v", err)
	}
	defer goldFile.Close()

	var goldMapping map[string]int
	decoder := json.NewDecoder(goldFile)
	if err := decoder.Decode(&goldMapping); err != nil {
		t.Fatalf("Failed to decode gold mapping: %v", err)
	}

	totalFiles := 0
	passedFiles := 0
	failedFiles := 0
	nLL := 0

	// Walk through the files and compare with gold mapping
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".nf" {
			totalFiles++
			relPath, _ := filepath.Rel(dir, path)

			result, err := BuildCST(path)
			if err != nil {
				t.Errorf("Failed to open file %s: %s", path, err)
				failedFiles++
				return nil
			}
			if result.Mode == "LL" {
				nLL++
			}
			builder := NewASTBuilder(path)
			ast := builder.Visit(result.Tree).(*ModuleNode)
			numStatements := getStatementsCount(ast)

			expectedCount, ok := goldMapping[relPath]
			if !ok {
				t.Errorf("File %s not found in gold mapping", relPath)
				failedFiles++
			} else if numStatements != expectedCount {
				t.Errorf("Mismatch for file %s: expected %d statements, got %d", relPath, expectedCount, numStatements)
				failedFiles++
			} else {
				passedFiles++
				t.Logf("File %s: Passed (expected %d, got %d)", relPath, expectedCount, numStatements)
			}

			// Print progress
			t.Logf("Progress: %d/%d files processed, %d passed, %d failed", totalFiles, len(goldMapping), passedFiles, failedFiles)
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Error walking the path %s: %v", dir, err)
	}

	// Print final summary
	t.Logf("Final Summary: %d/%d files processed, %d passed, %d failed, %d LL", totalFiles, len(goldMapping), passedFiles, failedFiles, nLL)
}

/*
func TestGenerateGoldMapping(t *testing.T) {
	debug.SetGCPercent(-1)
	dir := filepath.Join("testdata", "sarek")
	goldMapping := make(map[string]int)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && filepath.Ext(path) == ".nf" {
			input, err := antlr.NewFileStream(path)
			if err != nil {
				t.Fatalf("Failed to open file %s: %s", path, err)
			}

			lexer := NewGroovyLexer(input)
			stream := antlr.NewCommonTokenStream(lexer, 0)
			stream.Fill()
			parser := NewGroovyParser(stream)

			tree := parser.CompilationUnit()
			builder := NewASTBuilder(path)
			ast := builder.Visit(tree).(*ModuleNode)
			numStatements := getStatementsCount(ast)

			relPath, _ := filepath.Rel(dir, path)
			goldMapping[relPath] = numStatements
		}
		return nil
	})

	if err != nil {
		t.Fatalf("Error walking the path %s: %v", dir, err)
	}

	// Save the gold mapping to a JSON file
	goldFile, err := os.Create(filepath.Join("testdata", "sarek_gold_mapping.json"))
	if err != nil {
		t.Fatalf("Failed to create gold mapping file: %v", err)
	}
	defer goldFile.Close()

	encoder := json.NewEncoder(goldFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(goldMapping); err != nil {
		t.Fatalf("Failed to encode gold mapping: %v", err)
	}

	t.Logf("Gold mapping generated and saved to testdata/sarek_gold_mapping.json")
}
*/

func TestAmpliseqPipelineMain(t *testing.T) {
	filePath := filepath.Join("testdata", "ampliseq_pipeline_main.nf")
	ast := BuildASTTest(filePath)
	_ = ast
}

func TestCrisprseqScreening(t *testing.T) {
	filePath := filepath.Join("testdata", "crisprseq_screening.nf")
	_, err := BuildAST(filePath)
	if err == nil {
		t.Fatalf("Failed to report AST building error")
	}
}

func TestCrisprseqTargeted(t *testing.T) {
	filePath := filepath.Join("testdata", "crisprseq_targeted.nf")
	_, err := BuildAST(filePath)
	if err == nil {
		t.Fatalf("Failed to report AST building error")
	}
}

func TestChannelEdgeCase(t *testing.T) {
	filePath := filepath.Join("testdata", "channel_edge_case.nf")
	_, err := BuildAST(filePath)
	if err == nil {
		t.Fatalf("Failed to report AST building error")
	}
}

func TestRead2TreeMain(t *testing.T) {
	filePath := filepath.Join("testdata", "read2tree_main.nf")
	_, err := BuildAST(filePath)
	if err == nil {
		t.Fatalf("Failed to report AST building error")
	}
}

func TestOncoanalyzerMRD(t *testing.T) {
	filePath := filepath.Join("testdata", "mrd.nf")
	ast := BuildASTTest(filePath)
	_ = ast
}

func TestFastqFastaContigExtensionCobraMain(t *testing.T) {
	filePath := filepath.Join("testdata", "fastq_fasta_contig_extension_cobra_main.nf")
	ast := BuildASTTest(filePath)
	_ = ast
}

func TestQuantMSCreateInputChannel(t *testing.T) {
	filePath := filepath.Join("testdata", "create_input_channel.nf")
	ast := BuildASTTest(filePath)
	_ = ast
}

func TestVCFChrExtractMain(t *testing.T) {
	filePath := filepath.Join("testdata", "vcfchrextract_main.nf")
	ast, _ := BuildAST(filePath)
	_ = ast
}

func TestPreprocessRnaseqMain(t *testing.T) {
	filePath := filepath.Join("testdata", "preprocess_rnaseq_main.nf")
	ast := BuildASTTest(filePath)
	_ = ast
}

func TestDatasyncSync(t *testing.T) {
	filePath := filepath.Join("testdata", "datasync_sync.nf")
	_, err := BuildAST(filePath)
	if err == nil {
		t.Fatalf("Failed to report AST building error")
	}
}

func TestVCFGenotypeAnnotator(t *testing.T) {
	filePath := filepath.Join("testdata", "vcf_genotype_annotator.nf")
	_, err := BuildAST(filePath)
	if err == nil {
		t.Fatalf("Failed to report AST building error")
	}
}

func TestQiime2Intree(t *testing.T) {
	filePath := filepath.Join("testdata", "qiime2_intree.nf")

	// Capture stdout and stderr
	oldStdout := os.Stdout
	oldStderr := os.Stderr
	r, w, _ := os.Pipe()
	os.Stdout = w
	os.Stderr = w
	_, err := BuildAST(filePath)
	if err != nil {
		t.Fatalf("Error building AST")
	}
	// Restore original stdout and stderr
	err = w.Close()
	if err != nil {
		t.Errorf("Error closing pipe: %v", err)
	}
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	// Read the captured output
	var buf bytes.Buffer
	_, err = io.Copy(&buf, r)
	if err != nil {
		t.Errorf("Error copying buf: %v", err)
	}
	capturedOutput := buf.String()

	// Assert that no output was produced
	if capturedOutput != "" {
		t.Errorf("Expected no output, but got: %s", capturedOutput)
	}
}
