package parser

type UnaryPlusExpression struct {
	*BaseExpression
	expression Expression
}

func NewUnaryPlusExpression(expression Expression) *UnaryPlusExpression {
	return &UnaryPlusExpression{BaseExpression: NewBaseExpression(), expression: expression}
}

func (u *UnaryPlusExpression) GetExpression() Expression {
	return u.expression
}

func (u *UnaryPlusExpression) GetText() string {
	return "+(" + u.GetExpression().GetText() + ")"
}

func (u *UnaryPlusExpression) GetType() IClassNode {
	return u.GetExpression().GetType() // TODO: promote byte, char and short to int
}

// IsDynamic is deprecated
func (u *UnaryPlusExpression) IsDynamic() bool {
	return false
}

func (u *UnaryPlusExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewUnaryPlusExpression(transformer.Transform(u.GetExpression()))
	ret.SetSourcePosition(u)
	ret.CopyNodeMetaData(u)
	return ret
}

func (u *UnaryPlusExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitUnaryPlusExpression(u)
}
