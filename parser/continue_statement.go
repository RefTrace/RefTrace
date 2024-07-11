package parser

// ContinueStatement represents a continue statement in a loop statement
type ContinueStatement struct {
	Statement
	label string
}

// NewContinueStatement creates a new ContinueStatement
func NewContinueStatement(label string) *ContinueStatement {
	return &ContinueStatement{
		label: label,
	}
}

// GetLabel returns the label of the ContinueStatement
func (c *ContinueStatement) GetLabel() string {
	return c.label
}

// Visit implements the Statement interface
func (c *ContinueStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitContinueStatement(c)
}
