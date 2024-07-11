package parser

type CaseStatement struct {
	Statement
	code       Statement
	expression Expression
}

func NewCaseStatement(expression Expression, code Statement) *CaseStatement {
	return &CaseStatement{
		expression: expression,
		code:       code,
	}
}

func (c *CaseStatement) GetCode() Statement {
	return c.code
}

func (c *CaseStatement) SetCode(code Statement) {
	c.code = code
}

func (c *CaseStatement) GetExpression() Expression {
	return c.expression
}

func (c *CaseStatement) SetExpression(e Expression) {
	c.expression = e
}

func (c *CaseStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitCaseStatement(c)
}

func (c *CaseStatement) String() string {
	return c.Statement.GetText() + "[expression: " + c.expression.GetText() + "; code: " + c.code.GetText() + "]"
}
