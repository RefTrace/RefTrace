package parser

// SynchronizedStatement represents a synchronized statement
type SynchronizedStatement struct {
	Statement
	code       Statement
	expression Expression
}

// NewSynchronizedStatement creates a new SynchronizedStatement
func NewSynchronizedStatement(expression Expression, code Statement) *SynchronizedStatement {
	return &SynchronizedStatement{
		code:       code,
		expression: expression,
	}
}

// GetCode returns the code block of the synchronized statement
func (s *SynchronizedStatement) GetCode() Statement {
	return s.code
}

// SetCode sets the code block of the synchronized statement
func (s *SynchronizedStatement) SetCode(statement Statement) {
	s.code = statement
}

// GetExpression returns the expression of the synchronized statement
func (s *SynchronizedStatement) GetExpression() Expression {
	return s.expression
}

// Visit implements the Statement interface
func (s *SynchronizedStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitSynchronizedStatement(s)
}

// SetExpression sets the expression of the synchronized statement
func (s *SynchronizedStatement) SetExpression(expression Expression) {
	s.expression = expression
}
