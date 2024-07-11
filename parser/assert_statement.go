package parser

// AssertStatement represents an assert statement.
// E.g.:
// assert i != 0 : "should never be zero"
type AssertStatement struct {
	Statement
	booleanExpression BooleanExpression
	messageExpression Expression
}

// NewAssertStatement creates a new AssertStatement with a boolean expression
func NewAssertStatement(booleanExpression BooleanExpression) *AssertStatement {
	return NewAssertStatementWithMessage(booleanExpression, NullX())
}

// NewAssertStatementWithMessage creates a new AssertStatement with a boolean expression and message
func NewAssertStatementWithMessage(booleanExpression BooleanExpression, messageExpression Expression) *AssertStatement {
	return &AssertStatement{
		booleanExpression: booleanExpression,
		messageExpression: messageExpression,
	}
}

// Visit implements the Statement interface
func (a *AssertStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitAssertStatement(a)
}

// GetMessageExpression returns the message expression
func (a *AssertStatement) GetMessageExpression() Expression {
	return a.messageExpression
}

// GetBooleanExpression returns the boolean expression
func (a *AssertStatement) GetBooleanExpression() BooleanExpression {
	return a.booleanExpression
}

// SetBooleanExpression sets the boolean expression
func (a *AssertStatement) SetBooleanExpression(booleanExpression BooleanExpression) {
	a.booleanExpression = booleanExpression
}

// SetMessageExpression sets the message expression
func (a *AssertStatement) SetMessageExpression(messageExpression Expression) {
	a.messageExpression = messageExpression
}
