package parser

// BreakStatement represents a break statement in a switch or loop statement
type BreakStatement struct {
	Statement
	label string
}

// NewBreakStatement creates a new BreakStatement with an optional label
func NewBreakStatement(label string) *BreakStatement {
	return &BreakStatement{
		label: label,
	}
}

// GetLabel returns the label of the BreakStatement
func (b *BreakStatement) GetLabel() string {
	return b.label
}

// Visit calls the VisitBreakStatement method of the GroovyCodeVisitor
func (b *BreakStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitBreakStatement(b)
}
