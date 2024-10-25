package parser

// CatchStatement represents a catch (Exception var) { } statement
type CatchStatement struct {
	*BaseStatement
	variable *Parameter
	code     Statement
}

// NewCatchStatement creates a new CatchStatement
func NewCatchStatement(variable *Parameter, code Statement) *CatchStatement {
	return &CatchStatement{
		BaseStatement: NewBaseStatement(),
		variable:      variable,
		code:          code,
	}
}

// Visit implements the Statement interface
func (c *CatchStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitCatchStatement(c)
}

// GetCode returns the code block of the catch statement
func (c *CatchStatement) GetCode() Statement {
	return c.code
}

// GetExceptionType returns the exception type of the catch statement
func (c *CatchStatement) GetExceptionType() IClassNode {
	return c.variable.GetType()
}

// GetVariable returns the variable of the catch statement
func (c *CatchStatement) GetVariable() *Parameter {
	return c.variable
}

// SetCode sets the code block of the catch statement
func (c *CatchStatement) SetCode(code Statement) {
	c.code = code
}
