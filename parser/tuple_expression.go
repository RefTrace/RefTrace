package parser

import (
	"strings"
)

type ITupleExpression interface {
	PrependExpression(expression Expression) ITupleExpression
	AddExpression(expression Expression)
	GetExpression(i int) Expression
	GetExpressions() []Expression
	Visit(visitor GroovyCodeVisitor)
	//TransformExpression(transformer ExpressionTransformer) Expression
	GetText() string
}

var _ ITupleExpression = (*TupleExpression)(nil)

type TupleExpression struct {
	*BaseExpression
	expressions []Expression
}

func NewTupleExpression() *TupleExpression {
	return &TupleExpression{BaseExpression: NewBaseExpression(), expressions: make([]Expression, 0)}
}

func NewTupleExpressionWithCapacity(capacity int) *TupleExpression {
	return &TupleExpression{BaseExpression: NewBaseExpression(), expressions: make([]Expression, 0, capacity)}
}

func NewTupleExpressionWithExpressions(expressions ...Expression) *TupleExpression {
	return &TupleExpression{BaseExpression: NewBaseExpression(), expressions: expressions}
}

func (t *TupleExpression) PrependExpression(expression Expression) ITupleExpression {
	t.expressions = append([]Expression{expression}, t.expressions...)
	return t
}

func (t *TupleExpression) AddExpression(expression Expression) {
	t.expressions = append(t.expressions, expression)
}

func (t *TupleExpression) GetExpression(i int) Expression {
	return t.expressions[i]
}

func (t *TupleExpression) GetExpressions() []Expression {
	return t.expressions
}

func (t *TupleExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitTupleExpression(t)
}

func (t *TupleExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewTupleExpressionWithExpressions(TransformExpressions(t.GetExpressions(), transformer)...)
	ret.SetSourcePosition(t)
	ret.CopyNodeMetaData(t)
	return ret
}

func (t *TupleExpression) GetText() string {
	var buffer strings.Builder
	buffer.WriteString("(")
	for i, expression := range t.GetExpressions() {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(expression.GetText())
	}
	buffer.WriteString(")")
	return buffer.String()
}

func (t *TupleExpression) String() string {
	return t.BaseExpression.GetText() + t.GetText()
}
