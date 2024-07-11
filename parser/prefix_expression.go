package parser

import (
	"fmt"
)

// PrefixExpression represents a prefix expression like ++foo or --bar
type PrefixExpression struct {
	Expression
	operation  *Token
	expression Expression
}

func NewPrefixExpression(operation *Token, expression Expression) *PrefixExpression {
	pe := &PrefixExpression{
		operation: operation,
	}
	pe.SetExpression(expression)
	return pe
}

func (pe *PrefixExpression) SetExpression(expression Expression) {
	pe.expression = expression
}

func (pe *PrefixExpression) GetExpression() Expression {
	return pe.expression
}

func (pe *PrefixExpression) GetOperation() *Token {
	return pe.operation
}

func (pe *PrefixExpression) GetText() string {
	return fmt.Sprintf("(%s%s)", pe.GetOperation().GetText(), pe.GetExpression().GetText())
}

func (pe *PrefixExpression) GetType() *ClassNode {
	return pe.GetExpression().GetType()
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("%s[%s%s]", pe.Expression.GetText(), pe.GetOperation(), pe.GetExpression())
}

func (pe *PrefixExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewPrefixExpression(pe.GetOperation(), transformer.Transform(pe.GetExpression()))
	ret.SetSourcePosition(pe)
	ret.CopyNodeMetaData(pe)
	return ret
}

func (pe *PrefixExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitPrefixExpression(pe)
}
