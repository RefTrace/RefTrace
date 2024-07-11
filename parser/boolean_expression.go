package parser

// BooleanExpression represents a boolean expression
type BooleanExpression struct {
	Expression
	expression Expression
}

// NewBooleanExpression creates a new BooleanExpression
func NewBooleanExpression(expression Expression) *BooleanExpression {
	be := &BooleanExpression{
		expression: expression,
	}
	be.SetType(BooleanType) // Assuming BooleanType is defined elsewhere
	return be
}

// GetExpression returns the underlying expression
func (be *BooleanExpression) GetExpression() Expression {
	return be.expression
}

// GetText returns the text representation of the expression
func (be *BooleanExpression) GetText() string {
	return be.expression.GetText()
}

// TransformExpression transforms the boolean expression
func (be *BooleanExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewBooleanExpression(transformer.Transform(be.GetExpression()))
	ret.SetSourcePosition(be)
	ret.CopyNodeMetaData(be)
	return ret
}

// Visit implements the Visitable interface
func (be *BooleanExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitBooleanExpression(be)
}
