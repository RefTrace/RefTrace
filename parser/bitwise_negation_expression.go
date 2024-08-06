package parser

type BitwiseNegationExpression struct {
	*BaseExpression
	expression Expression
}

func NewBitwiseNegationExpression(expression Expression) *BitwiseNegationExpression {
	return &BitwiseNegationExpression{BaseExpression: NewBaseExpression(), expression: expression}
}

func (b *BitwiseNegationExpression) GetExpression() Expression {
	return b.expression
}

func (b *BitwiseNegationExpression) GetText() string {
	return "~(" + b.GetExpression().GetText() + ")"
}

func (b *BitwiseNegationExpression) GetType() *ClassNode {
	exprType := b.GetExpression().GetType()
	if IsStringType(exprType) || IsGStringType(exprType) {
		return PATTERN_TYPE // Assuming PATTERN_TYPE is defined elsewhere
	}
	return exprType
}

func (b *BitwiseNegationExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewBitwiseNegationExpression(transformer.Transform(b.GetExpression()))
	ret.SetSourcePosition(b)
	ret.CopyNodeMetaData(b)
	return ret
}

func (b *BitwiseNegationExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitBitwiseNegationExpression(b)
}
