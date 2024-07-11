package parser

import (
	"fmt"
)

// ExpressionStatement represents a simple statement such as a method call
// where the return value is ignored
type ExpressionStatement struct {
	Statement
	expression Expression
}

// NewExpressionStatement creates a new ExpressionStatement
func NewExpressionStatement(expression Expression) (*ExpressionStatement, error) {
	if expression == nil {
		return nil, fmt.Errorf("expression cannot be nil")
	}
	return &ExpressionStatement{expression: expression}, nil
}

// Visit implements the Statement interface
func (e *ExpressionStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitExpressionStatement(e)
}

// GetExpression returns the expression
func (e *ExpressionStatement) GetExpression() Expression {
	return e.expression
}

// SetExpression sets the expression
func (e *ExpressionStatement) SetExpression(expression Expression) {
	e.expression = expression
}

// GetText implements the Statement interface
func (e *ExpressionStatement) GetText() string {
	return e.expression.GetText()
}

// String implements the Stringer interface
func (e *ExpressionStatement) String() string {
	return fmt.Sprintf("%s[expression:%v]", e.Statement.GetText(), e.expression)
}
