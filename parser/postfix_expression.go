package parser

import (
	"fmt"
)

// PostfixExpression represents a postfix expression like foo++ or bar++
type PostfixExpression struct {
	*BaseExpression
	operation  *Token
	expression Expression
}

func NewPostfixExpression(expression Expression, operation *Token) *PostfixExpression {
	p := &PostfixExpression{
		BaseExpression: NewBaseExpression(),
		operation:      operation,
	}
	p.SetExpression(expression)
	return p
}

func (p *PostfixExpression) SetExpression(expression Expression) {
	p.expression = expression
}

func (p *PostfixExpression) GetExpression() Expression {
	return p.expression
}

func (p *PostfixExpression) GetOperation() *Token {
	return p.operation
}

func (p *PostfixExpression) GetText() string {
	return fmt.Sprintf("(%s%s)", p.GetExpression().GetText(), p.GetOperation().GetText())
}

func (p *PostfixExpression) GetType() IClassNode {
	return p.GetExpression().GetType()
}

func (p *PostfixExpression) String() string {
	return fmt.Sprintf("%s[%s%s]", p.BaseExpression.GetText(), p.GetExpression(), p.GetOperation().GetText())
}

func (p *PostfixExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewPostfixExpression(transformer.Transform(p.GetExpression()), p.GetOperation())
	ret.SetSourcePosition(p)
	ret.CopyNodeMetaData(p)
	return ret
}

func (p *PostfixExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitPostfixExpression(p)
}
