package parser

type GroovyCodeVisitor interface {
	// Statements
	VisitBlockStatement(*BlockStatement)
	VisitForLoop(*ForStatement)
	VisitWhileLoop(*WhileStatement)
	VisitDoWhileLoop(*DoWhileStatement)
	VisitIfElse(*IfStatement)
	VisitExpressionStatement(*ExpressionStatement)
	VisitReturnStatement(*ReturnStatement)
	VisitAssertStatement(*AssertStatement)
	VisitTryCatchFinally(*TryCatchStatement)
	VisitSwitch(*SwitchStatement)
	VisitCaseStatement(*CaseStatement)
	VisitBreakStatement(*BreakStatement)
	VisitContinueStatement(*ContinueStatement)
	VisitThrowStatement(*ThrowStatement)
	VisitSynchronizedStatement(*SynchronizedStatement)
	VisitCatchStatement(*CatchStatement)
	VisitEmptyStatement(*EmptyStatement)
	VisitStatement(Statement)

	// Expressions
	VisitMethodCallExpression(*MethodCallExpression)
	VisitStaticMethodCallExpression(*StaticMethodCallExpression)
	VisitConstructorCallExpression(*ConstructorCallExpression)
	VisitTernaryExpression(*TernaryExpression)
	VisitShortTernaryExpression(*ElvisOperatorExpression)
	VisitBinaryExpression(*BinaryExpression)
	VisitPrefixExpression(*PrefixExpression)
	VisitPostfixExpression(*PostfixExpression)
	VisitBooleanExpression(*BooleanExpression)
	VisitClosureExpression(*ClosureExpression)
	VisitLambdaExpression(*LambdaExpression)
	VisitTupleExpression(*TupleExpression)
	VisitMapExpression(*MapExpression)
	VisitMapEntryExpression(*MapEntryExpression)
	VisitListExpression(*ListExpression)
	VisitRangeExpression(*RangeExpression)
	VisitPropertyExpression(*PropertyExpression)
	VisitAttributeExpression(*AttributeExpression)
	VisitFieldExpression(*FieldExpression)
	VisitMethodPointerExpression(*MethodPointerExpression)
	VisitMethodReferenceExpression(*MethodReferenceExpression)
	VisitConstantExpression(*ConstantExpression)
	VisitClassExpression(*ClassExpression)
	VisitVariableExpression(*VariableExpression)
	VisitDeclarationExpression(*DeclarationExpression)
	VisitGStringExpression(*GStringExpression)
	VisitArrayExpression(*ArrayExpression)
	VisitSpreadExpression(*SpreadExpression)
	VisitSpreadMapExpression(*SpreadMapExpression)
	VisitNotExpression(*NotExpression)
	VisitUnaryMinusExpression(*UnaryMinusExpression)
	VisitUnaryPlusExpression(*UnaryPlusExpression)
	VisitBitwiseNegationExpression(*BitwiseNegationExpression)
	VisitCastExpression(*CastExpression)
	VisitArgumentlistExpression(*ArgumentListExpression)
	VisitClosureListExpression(*ClosureListExpression)
	VisitBytecodeExpression(*BytecodeExpression)
	VisitEmptyExpression(*EmptyExpression)
	VisitListOfExpressions([]Expression)
	VisitExpression(Expression)
}

// DefaultGroovyCodeVisitor provides default implementations for GroovyCodeVisitor
type DefaultGroovyCodeVisitor struct{}

func (v *DefaultGroovyCodeVisitor) VisitEmptyStatement(*EmptyStatement) {}

func (v *DefaultGroovyCodeVisitor) VisitStatement(stmt Statement) {
	if stmt != nil {
		stmt.Visit(v)
	}
}

func (v *DefaultGroovyCodeVisitor) VisitEmptyExpression(*EmptyExpression) {}

func (v *DefaultGroovyCodeVisitor) VisitListOfExpressions(list []Expression) {
	for _, expr := range list {
		expr.Visit(v)
	}
}

func (v *DefaultGroovyCodeVisitor) VisitExpression(expr Expression) {
	if expr != nil {
		expr.Visit(v)
	}
}

// Implement other methods of DefaultGroovyCodeVisitor as needed
