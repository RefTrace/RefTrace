package parser

// ThrowStatement represents a throw statement
type ThrowStatement struct {
	*BaseStatement
	expression Expression
}

// NewThrowStatement creates a new ThrowStatement
func NewThrowStatement(expression Expression) *ThrowStatement {
	return &ThrowStatement{
		BaseStatement: NewBaseStatement(),
		expression:    expression,
	}
}

// GetExpression returns the expression of the throw statement
func (t *ThrowStatement) GetExpression() Expression {
	return t.expression
}

// Accept implements the Visitable interface
func (t *ThrowStatement) Accept(visitor GroovyCodeVisitor) {
	visitor.VisitThrowStatement(t)
}

// SetExpression sets the expression of the throw statement
func (t *ThrowStatement) SetExpression(expression Expression) {
	t.expression = expression
}

// GetText returns the text representation of the throw statement
func (t *ThrowStatement) GetText() string {
	return "throw " + t.expression.GetText()
}
