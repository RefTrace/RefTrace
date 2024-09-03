package nf

import (
	"fmt"
	"reft-go/parser"
)

var _ parser.GroovyCodeVisitor = (*WorkflowBodyVisitor)(nil)

// WorkflowMode represents the different modes of a workflow
type WorkflowMode int

const (
	// InitialMode represents the initial mode before any labels
	InitialMode WorkflowMode = iota
	// TakeMode represents the 'take' mode of a workflow
	TakeMode
	// MainMode represents the 'main' mode of a workflow
	MainMode
	// EmitMode represents the 'emit' mode of a workflow
	EmitMode
)

type WorkflowBodyVisitor struct {
	mode    WorkflowMode
	Takes   []string
	Emits   []string
	hasMain bool
	hasTake bool
	errors  []string
}

// NewWorkflowBodyVisitor creates a new WorkflowBodyVisitor
func NewWorkflowBodyVisitor() *WorkflowBodyVisitor {
	return &WorkflowBodyVisitor{mode: InitialMode}
}

// Statements
func (v *WorkflowBodyVisitor) VisitBlockStatement(block *parser.BlockStatement) {
	for _, statement := range block.GetStatements() {
		label := statement.GetStatementLabel()
		switch label {
		case "take":
			if v.mode != InitialMode {
				v.errors = append(v.errors, "take: must be the first section in the workflow")
			}
			v.mode = TakeMode
			v.hasTake = true
		case "main":
			if v.mode == EmitMode {
				v.errors = append(v.errors, "main: cannot come after emit:")
			}
			v.mode = MainMode
			v.hasMain = true
		case "emit":
			v.mode = EmitMode
		case "":
			if v.mode == InitialMode {
				v.mode = MainMode
			}
		default:
			v.errors = append(v.errors, fmt.Sprintf("Unknown label: %s", label))
		}
		v.VisitStatement(statement)
	}

	// Validate workflow structure
	if v.hasTake && !v.hasMain {
		v.errors = append(v.errors, "When take: is used, main: must also be present")
	}
}

func (v *WorkflowBodyVisitor) VisitForLoop(statement *parser.ForStatement) {
	v.VisitExpression(statement.GetCollectionExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *WorkflowBodyVisitor) VisitWhileLoop(statement *parser.WhileStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *WorkflowBodyVisitor) VisitDoWhileLoop(statement *parser.DoWhileStatement) {
	v.VisitStatement(statement.GetLoopBlock())
	v.VisitExpression(statement.GetBooleanExpression())
}

func (v *WorkflowBodyVisitor) VisitIfElse(statement *parser.IfStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetIfBlock())
	v.VisitStatement(statement.GetElseBlock())
}

func (v *WorkflowBodyVisitor) VisitExpressionStatement(statement *parser.ExpressionStatement) {
	expr := statement.GetExpression()
	if variable, ok := expr.(*parser.VariableExpression); ok {
		if v.mode == TakeMode {
			v.Takes = append(v.Takes, variable.GetText())
			return
		} else if v.mode == EmitMode {
			v.Emits = append(v.Emits, variable.GetText())
			return
		}
	}
	// TODO: check property .out
	// TODO: check binary expr
	if propExpr, ok := expr.(*parser.PropertyExpression); ok {
		prop := propExpr.GetText()
		if v.mode == EmitMode {
			v.Emits = append(v.Emits, prop)
			return
		}
	}
	if binaryExpr, ok := expr.(*parser.BinaryExpression); ok {
		if binaryExpr.GetOperation().GetText() == "=" {
			if leftVar, ok := binaryExpr.GetLeftExpression().(*parser.VariableExpression); ok {
				if v.mode == EmitMode {
					v.Emits = append(v.Emits, leftVar.GetText())
				}
			}
			return
		}
	}
	v.VisitExpression(statement.GetExpression())
}

func (v *WorkflowBodyVisitor) GetErrors() []string {
	return v.errors
}

func (v *WorkflowBodyVisitor) VisitReturnStatement(statement *parser.ReturnStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitAssertStatement(statement *parser.AssertStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitExpression(statement.GetMessageExpression())
}

func (v *WorkflowBodyVisitor) VisitTryCatchFinally(statement *parser.TryCatchStatement) {
	for _, resource := range statement.GetResourceStatements() {
		v.VisitStatement(resource)
	}
	v.VisitStatement(statement.GetTryStatement())
	for _, catchStatement := range statement.GetCatchStatements() {
		v.VisitStatement(catchStatement)
	}
	v.VisitStatement(statement.GetFinallyStatement())
}

func (v *WorkflowBodyVisitor) VisitSwitch(statement *parser.SwitchStatement) {
	v.VisitExpression(statement.GetExpression())
	for _, caseStatement := range statement.GetCaseStatements() {
		v.VisitStatement(caseStatement)
	}
	v.VisitStatement(statement.GetDefaultStatement())
}

func (v *WorkflowBodyVisitor) VisitCaseStatement(statement *parser.CaseStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *WorkflowBodyVisitor) VisitBreakStatement(statement *parser.BreakStatement) {}

func (v *WorkflowBodyVisitor) VisitContinueStatement(statement *parser.ContinueStatement) {}

func (v *WorkflowBodyVisitor) VisitThrowStatement(statement *parser.ThrowStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitSynchronizedStatement(statement *parser.SynchronizedStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *WorkflowBodyVisitor) VisitCatchStatement(statement *parser.CatchStatement) {
	v.VisitStatement(statement.GetCode())
}

func (v *WorkflowBodyVisitor) VisitEmptyStatement(statement *parser.EmptyStatement) {}

func (v *WorkflowBodyVisitor) VisitStatement(statement parser.Statement) {
	statement.Visit(v)
}

// Expressions
func (v *WorkflowBodyVisitor) VisitMethodCallExpression(call *parser.MethodCallExpression) {
	v.VisitExpression(call.GetObjectExpression())
	v.VisitExpression(call.GetMethod())
	v.VisitExpression(call.GetArguments())
}

func (v *WorkflowBodyVisitor) VisitStaticMethodCallExpression(call *parser.StaticMethodCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *WorkflowBodyVisitor) VisitConstructorCallExpression(call *parser.ConstructorCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *WorkflowBodyVisitor) VisitTernaryExpression(expression *parser.TernaryExpression) {
	booleanExpr := expression.GetBooleanExpression()
	v.VisitExpression(booleanExpr)
	v.VisitExpression(expression.GetTrueExpression())
	v.VisitExpression(expression.GetFalseExpression())
}

func (v *WorkflowBodyVisitor) VisitShortTernaryExpression(expression *parser.ElvisOperatorExpression) {
	v.VisitTernaryExpression(expression.TernaryExpression)
}

func (v *WorkflowBodyVisitor) VisitBinaryExpression(expression *parser.BinaryExpression) {
	v.VisitExpression(expression.GetLeftExpression())
	v.VisitExpression(expression.GetRightExpression())
}

func (v *WorkflowBodyVisitor) VisitPrefixExpression(expression *parser.PrefixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitPostfixExpression(expression *parser.PostfixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitBooleanExpression(expression *parser.BooleanExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitClosureExpression(expression *parser.ClosureExpression) {
	if expression.IsParameterSpecified() {
		for _, parameter := range expression.GetParameters() {
			if parameter.HasInitialExpression() {
				v.VisitExpression(parameter.GetInitialExpression())
			}
		}
	}
	v.VisitStatement(expression.GetCode())
}

func (v *WorkflowBodyVisitor) VisitLambdaExpression(expression *parser.LambdaExpression) {
	v.VisitClosureExpression(expression.ClosureExpression)
}

func (v *WorkflowBodyVisitor) VisitTupleExpression(expression parser.ITupleExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *WorkflowBodyVisitor) VisitMapExpression(expression *parser.MapExpression) {
	entries := expression.GetMapEntryExpressions()
	exprs := make([]parser.Expression, len(entries))
	for i, entry := range entries {
		exprs[i] = entry
	}
	v.VisitListOfExpressions(exprs)
}

func (v *WorkflowBodyVisitor) VisitMapEntryExpression(expression *parser.MapEntryExpression) {
	v.VisitExpression(expression.GetKeyExpression())
	v.VisitExpression(expression.GetValueExpression())
}

func (v *WorkflowBodyVisitor) VisitListExpression(expression *parser.ListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *WorkflowBodyVisitor) VisitRangeExpression(expression *parser.RangeExpression) {
	v.VisitExpression(expression.GetFrom())
	v.VisitExpression(expression.GetTo())
}

func (v *WorkflowBodyVisitor) VisitPropertyExpression(expression *parser.PropertyExpression) {
	v.VisitExpression(expression.GetObjectExpression())
	v.VisitExpression(expression.GetProperty())
}

func (v *WorkflowBodyVisitor) VisitAttributeExpression(expression *parser.AttributeExpression) {
	v.VisitExpression(expression.GetObjectExpression())
	v.VisitExpression(expression.GetProperty())
}

func (v *WorkflowBodyVisitor) VisitFieldExpression(expression *parser.FieldExpression) {}

func (v *WorkflowBodyVisitor) VisitMethodPointerExpression(expression *parser.MethodPointerExpression) {
	v.VisitExpression(expression.GetExpression())
	v.VisitExpression(expression.GetMethodName())
}

func (v *WorkflowBodyVisitor) VisitMethodReferenceExpression(expression *parser.MethodReferenceExpression) {
	v.VisitMethodPointerExpression(expression.MethodPointerExpression)
}

func (v *WorkflowBodyVisitor) VisitConstantExpression(expression *parser.ConstantExpression) {}

func (v *WorkflowBodyVisitor) VisitClassExpression(expression *parser.ClassExpression) {}

func (v *WorkflowBodyVisitor) VisitVariableExpression(expression *parser.VariableExpression) {}

func (v *WorkflowBodyVisitor) VisitDeclarationExpression(expression *parser.DeclarationExpression) {
	v.VisitBinaryExpression(expression.BinaryExpression)
}

func (v *WorkflowBodyVisitor) VisitGStringExpression(expression *parser.GStringExpression) {
	v.VisitListOfExpressions(convertToExpressionSlice(expression.GetStrings()))
	v.VisitListOfExpressions(expression.GetValues())
}

func (v *WorkflowBodyVisitor) VisitArrayExpression(expression *parser.ArrayExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
	v.VisitListOfExpressions(expression.GetSizeExpression())
}

func (v *WorkflowBodyVisitor) VisitSpreadExpression(expression *parser.SpreadExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitSpreadMapExpression(expression *parser.SpreadMapExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitNotExpression(expression *parser.NotExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitUnaryMinusExpression(expression *parser.UnaryMinusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitUnaryPlusExpression(expression *parser.UnaryPlusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitBitwiseNegationExpression(expression *parser.BitwiseNegationExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitCastExpression(expression *parser.CastExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *WorkflowBodyVisitor) VisitArgumentlistExpression(expression *parser.ArgumentListExpression) {
	v.VisitTupleExpression(expression)
}

func (v *WorkflowBodyVisitor) VisitClosureListExpression(expression *parser.ClosureListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *WorkflowBodyVisitor) VisitEmptyExpression(expression *parser.EmptyExpression) {}

func (v *WorkflowBodyVisitor) VisitListOfExpressions(expressions []parser.Expression) {
	for _, expr := range expressions {
		v.VisitExpression(expr)
	}
}

func (v *WorkflowBodyVisitor) VisitExpression(expression parser.Expression) {
	expression.Visit(v)
}
