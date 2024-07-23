package parser

type NotExpression struct {
	*BooleanExpression
}

func NewNotExpression(expression Expression) *NotExpression {
	return &NotExpression{
		BooleanExpression: NewBooleanExpression(expression),
	}
}

func (n *NotExpression) GetText() string {
	return "!(" + n.BooleanExpression.GetText() + ")"
}

func (n *NotExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewNotExpression(transformer.Transform(n.GetExpression()))
	ret.SetSourcePosition(n)
	ret.CopyNodeMetaData(n)
	return ret
}

func (n *NotExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitNotExpression(n)
}
