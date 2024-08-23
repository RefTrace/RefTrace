package nf

import (
	"reft-go/parser"
)

var _ parser.GroovyCodeVisitor = (*ProcessBodyVisitor)(nil)

type ProcessMode int

const (
	InputMode ProcessMode = iota
	OutputMode
	WhenMode
	ScriptMode
)

type ProcessBodyVisitor struct {
	mode         ProcessMode
	hitDeclBlock bool
	inputs       []parser.Statement
	outputs      []parser.Statement
	directives   []Directive
}

// NewProcessBodyVisitor creates a new ProcessBodyVisitor
func NewProcessBodyVisitor() *ProcessBodyVisitor {
	return &ProcessBodyVisitor{mode: ScriptMode, hitDeclBlock: false}
}

type DirectiveType int

const (
	AcceleratorType DirectiveType = iota
	AfterScriptType
	// Add new types here
)

type Directive interface {
	Type() DirectiveType
}

type Accelerator struct {
	NumGPUs int
	GPUType string
}

func (a Accelerator) Type() DirectiveType { return AcceleratorType }

func makeAccelerator(mce *parser.MethodCallExpression) *Accelerator {
	var numGPUs int = -1
	var gpuType string = ""
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		for _, expr := range exprs {
			if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
				value := constantExpr.GetValue()
				if intValue, ok := value.(int); ok {
					numGPUs = intValue
				}
			}
			if mapExpr, ok := expr.(*parser.MapExpression); ok {
				entries := mapExpr.GetMapEntryExpressions()
				for _, entry := range entries {
					if entry.GetKeyExpression().GetText() == "type" {
						if constantExpr, ok := entry.GetValueExpression().(*parser.ConstantExpression); ok {
							gpuTypeVal := constantExpr.GetValue()
							if gpuTypeStr, ok := gpuTypeVal.(string); ok {
								gpuType = gpuTypeStr
							}
						}
					}
				}
			}
		}
	}
	if numGPUs != -1 {
		return &Accelerator{NumGPUs: numGPUs, GPUType: gpuType}
	}
	return nil
}

type AfterScript struct{}

func (a AfterScript) Type() DirectiveType { return AfterScriptType }

func makeDirective(statement parser.Statement) Directive {
	if exprStmt, ok := statement.(*parser.ExpressionStatement); ok {
		if mce, ok := exprStmt.GetExpression().(*parser.MethodCallExpression); ok {
			if mce.GetMethod().GetText() == "accelerator" {
				return makeAccelerator(mce)
			}
		}
	}
	return nil
}

func makeDirectives(statements []parser.Statement) []Directive {
	var directives []Directive
	for _, statement := range statements {
		directive := makeDirective(statement)
		if directive != nil {
			directives = append(directives, directive)
		}
	}
	return directives
}

// Statements
func (v *ProcessBodyVisitor) VisitBlockStatement(block *parser.BlockStatement) {
	stmts := block.GetStatements()
	v.inputs = findInputs(stmts)
	v.outputs = findOutputs(stmts)
	possibleDirectives := findPossibleDirectives(stmts)
	v.directives = makeDirectives(possibleDirectives)
	for _, statement := range stmts {
		v.VisitStatement(statement)
	}
}

func findPossibleDirectives(statements []parser.Statement) []parser.Statement {
	var directives []parser.Statement

	for _, statement := range statements {
		// If we find an "input" labeled statement, stop collecting
		if statement.GetStatementLabel() == "input" {
			break
		}

		// Add the statement to directives, regardless of line numbers
		directives = append(directives, statement)
	}

	return directives
}

func findInputs(statements []parser.Statement) []parser.Statement {
	var inputStatements []parser.Statement
	foundInput := false
	var lastLineNumber int

	for _, statement := range statements {
		if !foundInput {
			// Check if this statement has the "input" label
			if statement.GetStatementLabel() == "input" {
				foundInput = true
				inputStatements = append(inputStatements, statement)
				lastLineNumber = statement.GetLineNumber()
			}
		} else {
			// Check if the line number is contiguous
			if statement.GetLineNumber() == lastLineNumber+1 {
				inputStatements = append(inputStatements, statement)
				lastLineNumber = statement.GetLineNumber()
			} else {
				// Break the loop if we find a non-contiguous line
				break
			}
		}
	}

	return inputStatements
}

func findOutputs(statements []parser.Statement) []parser.Statement {
	var outputStatements []parser.Statement
	foundOutput := false
	var lastLineNumber int

	for _, statement := range statements {
		if !foundOutput {
			// Check if this statement has the "output" label
			if statement.GetStatementLabel() == "output" {
				foundOutput = true
				outputStatements = append(outputStatements, statement)
				lastLineNumber = statement.GetLineNumber()
			}
		} else {
			// Check if the line number is contiguous
			if statement.GetLineNumber() == lastLineNumber+1 {
				outputStatements = append(outputStatements, statement)
				lastLineNumber = statement.GetLineNumber()
			} else {
				// Break the loop if we find a non-contiguous line
				break
			}
		}
	}

	return outputStatements
}

func (v *ProcessBodyVisitor) VisitForLoop(statement *parser.ForStatement) {
	v.VisitExpression(statement.GetCollectionExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *ProcessBodyVisitor) VisitWhileLoop(statement *parser.WhileStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *ProcessBodyVisitor) VisitDoWhileLoop(statement *parser.DoWhileStatement) {
	v.VisitStatement(statement.GetLoopBlock())
	v.VisitExpression(statement.GetBooleanExpression())
}

func (v *ProcessBodyVisitor) VisitIfElse(statement *parser.IfStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetIfBlock())
	v.VisitStatement(statement.GetElseBlock())
}

func (v *ProcessBodyVisitor) VisitExpressionStatement(statement *parser.ExpressionStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *ProcessBodyVisitor) VisitReturnStatement(statement *parser.ReturnStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *ProcessBodyVisitor) VisitAssertStatement(statement *parser.AssertStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitExpression(statement.GetMessageExpression())
}

func (v *ProcessBodyVisitor) VisitTryCatchFinally(statement *parser.TryCatchStatement) {
	for _, resource := range statement.GetResourceStatements() {
		v.VisitStatement(resource)
	}
	v.VisitStatement(statement.GetTryStatement())
	for _, catchStatement := range statement.GetCatchStatements() {
		v.VisitStatement(catchStatement)
	}
	v.VisitStatement(statement.GetFinallyStatement())
}

func (v *ProcessBodyVisitor) VisitSwitch(statement *parser.SwitchStatement) {
	v.VisitExpression(statement.GetExpression())
	for _, caseStatement := range statement.GetCaseStatements() {
		v.VisitStatement(caseStatement)
	}
	v.VisitStatement(statement.GetDefaultStatement())
}

func (v *ProcessBodyVisitor) VisitCaseStatement(statement *parser.CaseStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *ProcessBodyVisitor) VisitBreakStatement(statement *parser.BreakStatement) {}

func (v *ProcessBodyVisitor) VisitContinueStatement(statement *parser.ContinueStatement) {}

func (v *ProcessBodyVisitor) VisitThrowStatement(statement *parser.ThrowStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *ProcessBodyVisitor) VisitSynchronizedStatement(statement *parser.SynchronizedStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *ProcessBodyVisitor) VisitCatchStatement(statement *parser.CatchStatement) {
	v.VisitStatement(statement.GetCode())
}

func (v *ProcessBodyVisitor) VisitEmptyStatement(statement *parser.EmptyStatement) {}

func (v *ProcessBodyVisitor) VisitStatement(statement parser.Statement) {
	statement.Visit(v)
}

// Expressions
func (v *ProcessBodyVisitor) VisitMethodCallExpression(call *parser.MethodCallExpression) {
	v.VisitExpression(call.GetObjectExpression())
	v.VisitExpression(call.GetMethod())
	v.VisitExpression(call.GetArguments())
}

func (v *ProcessBodyVisitor) VisitStaticMethodCallExpression(call *parser.StaticMethodCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *ProcessBodyVisitor) VisitConstructorCallExpression(call *parser.ConstructorCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *ProcessBodyVisitor) VisitTernaryExpression(expression *parser.TernaryExpression) {
	booleanExpr := expression.GetBooleanExpression()
	v.VisitExpression(booleanExpr)
	v.VisitExpression(expression.GetTrueExpression())
	v.VisitExpression(expression.GetFalseExpression())
}

func (v *ProcessBodyVisitor) VisitShortTernaryExpression(expression *parser.ElvisOperatorExpression) {
	v.VisitTernaryExpression(expression.TernaryExpression)
}

func (v *ProcessBodyVisitor) VisitBinaryExpression(expression *parser.BinaryExpression) {
	v.VisitExpression(expression.GetLeftExpression())
	v.VisitExpression(expression.GetRightExpression())
}

func (v *ProcessBodyVisitor) VisitPrefixExpression(expression *parser.PrefixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitPostfixExpression(expression *parser.PostfixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitBooleanExpression(expression *parser.BooleanExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitClosureExpression(expression *parser.ClosureExpression) {
	if expression.IsParameterSpecified() {
		for _, parameter := range expression.GetParameters() {
			if parameter.HasInitialExpression() {
				v.VisitExpression(parameter.GetInitialExpression())
			}
		}
	}
	v.VisitStatement(expression.GetCode())
}

func (v *ProcessBodyVisitor) VisitLambdaExpression(expression *parser.LambdaExpression) {
	v.VisitClosureExpression(expression.ClosureExpression)
}

func (v *ProcessBodyVisitor) VisitTupleExpression(expression parser.ITupleExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *ProcessBodyVisitor) VisitMapExpression(expression *parser.MapExpression) {
	entries := expression.GetMapEntryExpressions()
	exprs := make([]parser.Expression, len(entries))
	for i, entry := range entries {
		exprs[i] = entry
	}
	v.VisitListOfExpressions(exprs)
}

func (v *ProcessBodyVisitor) VisitMapEntryExpression(expression *parser.MapEntryExpression) {
	v.VisitExpression(expression.GetKeyExpression())
	v.VisitExpression(expression.GetValueExpression())
}

func (v *ProcessBodyVisitor) VisitListExpression(expression *parser.ListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *ProcessBodyVisitor) VisitRangeExpression(expression *parser.RangeExpression) {
	v.VisitExpression(expression.GetFrom())
	v.VisitExpression(expression.GetTo())
}

func (v *ProcessBodyVisitor) VisitPropertyExpression(expression *parser.PropertyExpression) {
	v.VisitExpression(expression.GetObjectExpression())
	v.VisitExpression(expression.GetProperty())
}

func (v *ProcessBodyVisitor) VisitAttributeExpression(expression *parser.AttributeExpression) {
	v.VisitExpression(expression.GetObjectExpression())
	v.VisitExpression(expression.GetProperty())
}

func (v *ProcessBodyVisitor) VisitFieldExpression(expression *parser.FieldExpression) {}

func (v *ProcessBodyVisitor) VisitMethodPointerExpression(expression *parser.MethodPointerExpression) {
	v.VisitExpression(expression.GetExpression())
	v.VisitExpression(expression.GetMethodName())
}

func (v *ProcessBodyVisitor) VisitMethodReferenceExpression(expression *parser.MethodReferenceExpression) {
	v.VisitMethodPointerExpression(expression.MethodPointerExpression)
}

func (v *ProcessBodyVisitor) VisitConstantExpression(expression *parser.ConstantExpression) {}

func (v *ProcessBodyVisitor) VisitClassExpression(expression *parser.ClassExpression) {}

func (v *ProcessBodyVisitor) VisitVariableExpression(expression *parser.VariableExpression) {}

func (v *ProcessBodyVisitor) VisitDeclarationExpression(expression *parser.DeclarationExpression) {
	v.VisitBinaryExpression(expression.BinaryExpression)
}

func (v *ProcessBodyVisitor) VisitGStringExpression(expression *parser.GStringExpression) {
	v.VisitListOfExpressions(convertToExpressionSlice(expression.GetStrings()))
	v.VisitListOfExpressions(expression.GetValues())
}

func (v *ProcessBodyVisitor) VisitArrayExpression(expression *parser.ArrayExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
	v.VisitListOfExpressions(expression.GetSizeExpression())
}

func (v *ProcessBodyVisitor) VisitSpreadExpression(expression *parser.SpreadExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitSpreadMapExpression(expression *parser.SpreadMapExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitNotExpression(expression *parser.NotExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitUnaryMinusExpression(expression *parser.UnaryMinusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitUnaryPlusExpression(expression *parser.UnaryPlusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitBitwiseNegationExpression(expression *parser.BitwiseNegationExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitCastExpression(expression *parser.CastExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *ProcessBodyVisitor) VisitArgumentlistExpression(expression *parser.ArgumentListExpression) {
	v.VisitTupleExpression(expression)
}

func (v *ProcessBodyVisitor) VisitClosureListExpression(expression *parser.ClosureListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *ProcessBodyVisitor) VisitEmptyExpression(expression *parser.EmptyExpression) {}

func (v *ProcessBodyVisitor) VisitListOfExpressions(expressions []parser.Expression) {
	for _, expr := range expressions {
		v.VisitExpression(expr)
	}
}

func (v *ProcessBodyVisitor) VisitExpression(expression parser.Expression) {
	expression.Visit(v)
}
