package parser

import (
	"strings"
)

// ClosureListExpression represents a list of expressions used to create closures.
// Example: def foo = (1;2;;)
// The right side is a ClosureListExpression consisting of two ConstantExpressions
// for the values 1 and 2, and two EmptyStatement entries. The ClosureListExpression
// defines a new variable scope. All created Closures share this scope.
type ClosureListExpression struct {
	ListExpression
	scope *VariableScope
}

func NewClosureListExpression(expressions []Expression) *ClosureListExpression {
	return &ClosureListExpression{
		ListExpression: ListExpression{expressions: expressions},
		scope:          NewVariableScope(),
	}
}

func NewEmptyClosureListExpression() *ClosureListExpression {
	return NewClosureListExpression(make([]Expression, 0, 3))
}

func (c *ClosureListExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitClosureListExpression(c)
}

func (c *ClosureListExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewClosureListExpression(TransformExpressions(c.GetExpressions(), transformer))
	ret.SetSourcePosition(c)
	ret.CopyNodeMetaData(c)
	return ret
}

func (c *ClosureListExpression) SetVariableScope(scope *VariableScope) {
	c.scope = scope
}

func (c *ClosureListExpression) GetVariableScope() *VariableScope {
	return c.scope
}

func (c *ClosureListExpression) GetText() string {
	var buffer strings.Builder
	buffer.WriteString("(")
	for i, expression := range c.GetExpressions() {
		if i > 0 {
			buffer.WriteString("; ")
		}
		buffer.WriteString(expression.GetText())
	}
	buffer.WriteString(")")
	return buffer.String()
}
