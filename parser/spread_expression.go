package parser

// SpreadExpression represents a spread expression *x in the list expression [1, *x, 2].
type SpreadExpression struct {
	Expression
	expression Expression
}

// NewSpreadExpression creates a new SpreadExpression with the given expression.
func NewSpreadExpression(expression Expression) *SpreadExpression {
	return &SpreadExpression{
		expression: expression,
	}
}

// GetExpression returns the underlying expression of the SpreadExpression.
func (s *SpreadExpression) GetExpression() Expression {
	return s.expression
}

// GetText returns the string representation of the SpreadExpression.
func (s *SpreadExpression) GetText() string {
	return "*" + s.expression.GetText()
}

// GetType returns the type of the underlying expression.
func (s *SpreadExpression) GetType() *ClassNode {
	return s.expression.GetType()
}

// TransformExpression transforms the SpreadExpression using the given transformer.
func (s *SpreadExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewSpreadExpression(transformer.Transform(s.GetExpression()))
	ret.SetSourcePosition(s)
	ret.CopyNodeMetaData(s)
	return ret
}

// Visit calls the VisitSpreadExpression method of the given GroovyCodeVisitor.
func (s *SpreadExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitSpreadExpression(s)
}
