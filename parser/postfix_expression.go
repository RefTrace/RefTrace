package parser

import (
	"fmt"
)

// PostfixExpression represents a postfix expression like foo++ or bar++
type PostfixExpression struct {
	Expression
	operation  *Token
	expression Expression
}

func NewPostfixExpression(expression Expression, operation *Token) *PostfixExpression {
	p := &PostfixExpression{
		operation: operation,
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

func (p *PostfixExpression) GetType() *ClassNode {
	return p.GetExpression().GetType()
}

func (p *PostfixExpression) String() string {
	return fmt.Sprintf("%s[%s%s]", p.Expression.GetText(), p.GetExpression(), p.GetOperation())
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