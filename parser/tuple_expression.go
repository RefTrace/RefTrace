package parser

import (
	"strings"
)

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

func (t *TupleExpression) PrependExpression(expression Expression) *TupleExpression {
	t.expressions = append([]Expression{expression}, t.expressions...)
	return t
}

func (t *TupleExpression) AddExpression(expression Expression) *TupleExpression {
	t.expressions = append(t.expressions, expression)
	return t
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
