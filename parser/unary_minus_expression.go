package parser

type UnaryMinusExpression struct {
	*BaseExpression
	expression Expression
}

func NewUnaryMinusExpression(expression Expression) *UnaryMinusExpression {
	return &UnaryMinusExpression{BaseExpression: NewBaseExpression(), expression: expression}
}

func (u *UnaryMinusExpression) GetExpression() Expression {
	return u.expression
}

func (u *UnaryMinusExpression) GetText() string {
	return "-(" + u.GetExpression().GetText() + ")"
}

func (u *UnaryMinusExpression) GetType() *ClassNode {
	return u.GetExpression().GetType()
}

// IsDynamic is deprecated and always returns false
func (u *UnaryMinusExpression) IsDynamic() bool {
	return false
}

func (u *UnaryMinusExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewUnaryMinusExpression(transformer.Transform(u.GetExpression()))
	ret.SetSourcePosition(u)
	ret.CopyNodeMetaData(u)
	return ret
}

func (u *UnaryMinusExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitUnaryMinusExpression(u)
}
