package parser

// CodeVisitorSupport is a base implementation of GroovyCodeVisitor
type CodeVisitorSupport struct {
	VisitReturnStatementFunc    func(statement *ReturnStatement)
	VisitThrowStatementFunc     func(statement *ThrowStatement)
	VisitPropertyExpressionFunc func(expression *PropertyExpression) // Added this line
}

// Statements
func (v *CodeVisitorSupport) VisitBlockStatement(block *BlockStatement) {
	for _, statement := range block.GetStatements() {
		v.VisitStatement(statement)
	}
}

func (v *CodeVisitorSupport) VisitForLoop(statement *ForStatement) {
	v.VisitExpression(statement.GetCollectionExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *CodeVisitorSupport) VisitWhileLoop(statement *WhileStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetLoopBlock())
}

func (v *CodeVisitorSupport) VisitDoWhileLoop(statement *DoWhileStatement) {
	v.VisitStatement(statement.GetLoopBlock())
	v.VisitExpression(statement.GetBooleanExpression())
}

func (v *CodeVisitorSupport) VisitIfElse(statement *IfStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitStatement(statement.GetIfBlock())
	v.VisitStatement(statement.GetElseBlock())
}

func (v *CodeVisitorSupport) VisitExpressionStatement(statement *ExpressionStatement) {
	v.VisitExpression(statement.GetExpression())
}

func (v *CodeVisitorSupport) VisitReturnStatement(statement *ReturnStatement) {
	if v.VisitReturnStatementFunc != nil {
		v.VisitReturnStatementFunc(statement)
	} else {
		v.VisitExpression(statement.GetExpression())
	}
}

func (v *CodeVisitorSupport) VisitAssertStatement(statement *AssertStatement) {
	v.VisitExpression(statement.GetBooleanExpression())
	v.VisitExpression(statement.GetMessageExpression())
}

func (v *CodeVisitorSupport) VisitTryCatchFinally(statement *TryCatchStatement) {
	for _, resource := range statement.GetResourceStatements() {
		v.VisitStatement(resource)
	}
	v.VisitStatement(statement.GetTryStatement())
	for _, catchStatement := range statement.GetCatchStatements() {
		v.VisitStatement(catchStatement)
	}
	v.VisitStatement(statement.GetFinallyStatement())
}

func (v *CodeVisitorSupport) VisitSwitch(statement *SwitchStatement) {
	v.VisitExpression(statement.GetExpression())
	for _, caseStatement := range statement.GetCaseStatements() {
		v.VisitStatement(caseStatement)
	}
	v.VisitStatement(statement.GetDefaultStatement())
}

func (v *CodeVisitorSupport) VisitCaseStatement(statement *CaseStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *CodeVisitorSupport) VisitBreakStatement(statement *BreakStatement) {}

func (v *CodeVisitorSupport) VisitContinueStatement(statement *ContinueStatement) {}

func (v *CodeVisitorSupport) VisitThrowStatement(statement *ThrowStatement) {
	if v.VisitThrowStatementFunc != nil {
		v.VisitThrowStatementFunc(statement)
	} else {
		v.VisitExpression(statement.GetExpression())
	}
}

func (v *CodeVisitorSupport) VisitSynchronizedStatement(statement *SynchronizedStatement) {
	v.VisitExpression(statement.GetExpression())
	v.VisitStatement(statement.GetCode())
}

func (v *CodeVisitorSupport) VisitCatchStatement(statement *CatchStatement) {
	v.VisitStatement(statement.GetCode())
}

func (v *CodeVisitorSupport) VisitEmptyStatement(statement *EmptyStatement) {}

func (v *CodeVisitorSupport) VisitStatement(statement Statement) {
	statement.Visit(v)
}

// Expressions
func (v *CodeVisitorSupport) VisitMethodCallExpression(call *MethodCallExpression) {
	v.VisitExpression(call.GetObjectExpression())
	v.VisitExpression(call.GetMethod())
	v.VisitExpression(call.GetArguments())
}

func (v *CodeVisitorSupport) VisitStaticMethodCallExpression(call *StaticMethodCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *CodeVisitorSupport) VisitConstructorCallExpression(call *ConstructorCallExpression) {
	v.VisitExpression(call.GetArguments())
}

func (v *CodeVisitorSupport) VisitTernaryExpression(expression *TernaryExpression) {
	booleanExpr := expression.GetBooleanExpression()
	v.VisitExpression(booleanExpr)
	v.VisitExpression(expression.GetTrueExpression())
	v.VisitExpression(expression.GetFalseExpression())
}

func (v *CodeVisitorSupport) VisitShortTernaryExpression(expression *ElvisOperatorExpression) {
	v.VisitTernaryExpression(expression.TernaryExpression)
}

func (v *CodeVisitorSupport) VisitBinaryExpression(expression *BinaryExpression) {
	v.VisitExpression(expression.GetLeftExpression())
	v.VisitExpression(expression.GetRightExpression())
}

func (v *CodeVisitorSupport) VisitPrefixExpression(expression *PrefixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitPostfixExpression(expression *PostfixExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitBooleanExpression(expression *BooleanExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitClosureExpression(expression *ClosureExpression) {
	if expression.IsParameterSpecified() {
		for _, parameter := range expression.GetParameters() {
			if parameter.HasInitialExpression() {
				v.VisitExpression(parameter.GetInitialExpression())
			}
		}
	}
	v.VisitStatement(expression.GetCode())
}

func (v *CodeVisitorSupport) VisitLambdaExpression(expression *LambdaExpression) {
	v.VisitClosureExpression(expression.ClosureExpression)
}

func (v *CodeVisitorSupport) VisitTupleExpression(expression *TupleExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *CodeVisitorSupport) VisitMapExpression(expression *MapExpression) {
	entries := expression.GetMapEntryExpressions()
	exprs := make([]Expression, len(entries))
	for i, entry := range entries {
		exprs[i] = entry
	}
	v.VisitListOfExpressions(exprs)
}

func (v *CodeVisitorSupport) VisitMapEntryExpression(expression *MapEntryExpression) {
	v.VisitExpression(expression.GetKeyExpression())
	v.VisitExpression(expression.GetValueExpression())
}

func (v *CodeVisitorSupport) VisitListExpression(expression *ListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *CodeVisitorSupport) VisitRangeExpression(expression *RangeExpression) {
	v.VisitExpression(expression.GetFrom())
	v.VisitExpression(expression.GetTo())
}

func (v *CodeVisitorSupport) VisitPropertyExpression(expression *PropertyExpression) {
	if v.VisitPropertyExpressionFunc != nil {
		v.VisitPropertyExpressionFunc(expression)
	} else {
		v.VisitExpression(expression.GetObjectExpression())
		v.VisitExpression(expression.GetProperty())
	}
}

func (v *CodeVisitorSupport) VisitAttributeExpression(expression *AttributeExpression) {
	v.VisitExpression(expression.GetObjectExpression())
	v.VisitExpression(expression.GetProperty())
}

func (v *CodeVisitorSupport) VisitFieldExpression(expression *FieldExpression) {}

func (v *CodeVisitorSupport) VisitMethodPointerExpression(expression *MethodPointerExpression) {
	v.VisitExpression(expression.GetExpression())
	v.VisitExpression(expression.GetMethodName())
}

func (v *CodeVisitorSupport) VisitMethodReferenceExpression(expression *MethodReferenceExpression) {
	v.VisitMethodPointerExpression(expression.MethodPointerExpression)
}

func (v *CodeVisitorSupport) VisitConstantExpression(expression *ConstantExpression) {}

func (v *CodeVisitorSupport) VisitClassExpression(expression *ClassExpression) {}

func (v *CodeVisitorSupport) VisitVariableExpression(expression *VariableExpression) {}

func (v *CodeVisitorSupport) VisitDeclarationExpression(expression *DeclarationExpression) {
	v.VisitBinaryExpression(expression.BinaryExpression)
}

func convertToExpressionSlice(constExprs []*ConstantExpression) []Expression {
	exprs := make([]Expression, len(constExprs))
	for i, ce := range constExprs {
		exprs[i] = ce
	}
	return exprs
}

func (v *CodeVisitorSupport) VisitGStringExpression(expression *GStringExpression) {
	v.VisitListOfExpressions(convertToExpressionSlice(expression.GetStrings()))
	v.VisitListOfExpressions(expression.GetValues())
}

func (v *CodeVisitorSupport) VisitArrayExpression(expression *ArrayExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
	v.VisitListOfExpressions(expression.GetSizeExpression())
}

func (v *CodeVisitorSupport) VisitSpreadExpression(expression *SpreadExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitSpreadMapExpression(expression *SpreadMapExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitNotExpression(expression *NotExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitUnaryMinusExpression(expression *UnaryMinusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitUnaryPlusExpression(expression *UnaryPlusExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitBitwiseNegationExpression(expression *BitwiseNegationExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitCastExpression(expression *CastExpression) {
	v.VisitExpression(expression.GetExpression())
}

func (v *CodeVisitorSupport) VisitArgumentlistExpression(expression *ArgumentListExpression) {
	v.VisitTupleExpression(expression.TupleExpression)
}

func (v *CodeVisitorSupport) VisitClosureListExpression(expression *ClosureListExpression) {
	v.VisitListOfExpressions(expression.GetExpressions())
}

func (v *CodeVisitorSupport) VisitEmptyExpression(expression *EmptyExpression) {}

func (v *CodeVisitorSupport) VisitListOfExpressions(expressions []Expression) {
	for _, expr := range expressions {
		v.VisitExpression(expr)
	}
}

func (v *CodeVisitorSupport) VisitExpression(expression Expression) {
	expression.Visit(v)
}
