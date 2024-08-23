package nf

import (
	"errors"
	"reft-go/parser"

	"reft-go/nf/directives"
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
	directives   []directives.Directive
}

// NewProcessBodyVisitor creates a new ProcessBodyVisitor
func NewProcessBodyVisitor() *ProcessBodyVisitor {
	return &ProcessBodyVisitor{mode: ScriptMode, hitDeclBlock: false}
}

func makeDirective(statement parser.Statement) (directives.Directive, error) {
	if exprStmt, ok := statement.(*parser.ExpressionStatement); ok {
		if mce, ok := exprStmt.GetExpression().(*parser.MethodCallExpression); ok {
			if mce.GetMethod().GetText() == "accelerator" {
				directive, err := directives.MakeAccelerator(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "afterScript" {
				directive, err := directives.MakeAfterScript(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "arch" {
				directive, err := directives.MakeArch(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "array" {
				directive, err := directives.MakeArrayDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "beforeScript" {
				directive, err := directives.MakeBeforeScript(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "cache" {
				directive, err := directives.MakeCacheDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "clusterOptions" {
				directive, err := directives.MakeClusterOptions(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "conda" {
				directive, err := directives.MakeConda(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "container" {
				directive, err := directives.MakeContainer(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "containerOptions" {
				directive, err := directives.MakeContainerOptions(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "cpus" {
				directive, err := directives.MakeCpusDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "debug" {
				directive, err := directives.MakeDebugDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "disk" {
				directive, err := directives.MakeDiskDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "echo" {
				directive, err := directives.MakeEchoDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "errorStrategy" {
				directive, err := directives.MakeErrorStrategyDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "executor" {
				directive, err := directives.MakeExecutorDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "ext" {
				directive, err := directives.MakeExtDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "fair" {
				directive, err := directives.MakeFairDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "label" {
				directive, err := directives.MakeLabelDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "machineType" {
				directive, err := directives.MakeMachineTypeDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "maxSubmitAwait" {
				directive, err := directives.MakeMaxSubmitAwaitDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "maxErrors" {
				directive, err := directives.MakeMaxErrorsDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "maxForks" {
				directive, err := directives.MakeMaxForksDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "maxRetries" {
				directive, err := directives.MakeMaxRetriesDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "memory" {
				directive, err := directives.MakeMemoryDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "module" {
				directive, err := directives.MakeModuleDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "penv" {
				directive, err := directives.MakePenvDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "pod" {
				directive, err := directives.MakePodDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "publishDir" {
				directive, err := directives.MakePublishDirDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "queue" {
				directive, err := directives.MakeQueueDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "resourceLabels" {
				directive, err := directives.MakeResourceLabelsDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
			if mce.GetMethod().GetText() == "resourceLimits" {
				directive, err := directives.MakeResourceLimitsDirective(mce)
				if err == nil {
					return directive, nil
				}
				return nil, err
			}
		}
	}
	return nil, errors.New("unknown directive")
}

func makeDirectives(statements []parser.Statement) []directives.Directive {
	var directives []directives.Directive
	for _, statement := range statements {
		directive, err := makeDirective(statement)
		if err == nil {
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
