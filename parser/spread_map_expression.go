package parser

import (
	"fmt"
)

// SpreadMapExpression represents a spread map expression *:m
// in the map expression [1, *:m, 2, "c":100]
// or in the method invoke expression func(1, *:m, 2, "c":100).
type SpreadMapExpression struct {
	Expression Expression
}

func NewSpreadMapExpression(expression Expression) *SpreadMapExpression {
	return &SpreadMapExpression{Expression: expression}
}

func (s *SpreadMapExpression) GetExpression() Expression {
	return s.Expression
}

func (s *SpreadMapExpression) GetText() string {
	return fmt.Sprintf("*:%s", s.GetExpression().GetText())
}

func (s *SpreadMapExpression) GetType() *ClassNode {
	return s.GetExpression().GetType()
}

func (s *SpreadMapExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewSpreadMapExpression(transformer.Transform(s.GetExpression()))
	ret.SetSourcePosition(s)
	ret.CopyNodeMetaData(s)
	return ret
}

func (s *SpreadMapExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitSpreadMapExpression(s)
}
