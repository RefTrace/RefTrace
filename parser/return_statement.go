package parser

import (
	"fmt"
)

// ReturnStatement represents a return statement
type ReturnStatement struct {
	Statement
	expression Expression
}

// RETURN_NULL_OR_VOID is only used for synthetic return statements emitted by the compiler.
// For comparisons use IsReturningNullOrVoid() instead.
var RETURN_NULL_OR_VOID = &ReturnStatement{expression: nullX()}

// NewReturnStatementFromExpressionStatement creates a new ReturnStatement from an ExpressionStatement
func NewReturnStatementFromExpressionStatement(statement *ExpressionStatement) *ReturnStatement {
	rs := &ReturnStatement{expression: statement.GetExpression()}
	rs.CopyStatementLabels(&statement.Statement)
	return rs
}

// NewReturnStatement creates a new ReturnStatement with the given expression
func NewReturnStatement(expression Expression) *ReturnStatement {
	return &ReturnStatement{expression: expression}
}

// GetExpression returns the expression of the ReturnStatement
func (rs *ReturnStatement) GetExpression() Expression {
	return rs.expression
}

// SetExpression sets the expression of the ReturnStatement
func (rs *ReturnStatement) SetExpression(expression Expression) {
	rs.expression = expression
}

// GetText returns the text representation of the ReturnStatement
func (rs *ReturnStatement) GetText() string {
	return "return " + rs.expression.GetText()
}

// IsReturningNullOrVoid checks if the ReturnStatement is returning null or void
func (rs *ReturnStatement) IsReturningNullOrVoid() bool {
	if ce, ok := rs.expression.(*ConstantExpression); ok {
		return ce.IsNullExpression()
	}
	return false
}

// String returns a string representation of the ReturnStatement
func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("ReturnStatement[expression:%v]", rs.expression)
}

// Visit calls the VisitReturnStatement method of the GroovyCodeVisitor
func (rs *ReturnStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitReturnStatement(rs)
}
