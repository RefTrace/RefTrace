package parser

import (
	"strings"
)

// LambdaExpression represents a lambda expression such as:
// e -> e * 2
// (x, y) -> x + y
// (x, y) -> { x + y }
// (int x, int y) -> { x + y }
type LambdaExpression struct {
	*ClosureExpression
	serializable bool
}

// NewLambdaExpression creates a new LambdaExpression
func NewLambdaExpression(parameters []*Parameter, code Statement) *LambdaExpression {
	le := &LambdaExpression{
		ClosureExpression: NewClosureExpression(parameters, code),
	}
	return le
}

// Visit implements the GroovyCodeVisitor interface
func (l *LambdaExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitLambdaExpression(l)
}

// GetText returns a string representation of the lambda expression
func (l *LambdaExpression) GetText() string {
	paramText := getParametersText(l.GetParameters())
	if len(paramText) > 0 {
		return "(" + paramText + ") -> { ... }"
	}
	return "() -> { ... }"
}

// IsSerializable returns whether the lambda expression is serializable
func (l *LambdaExpression) IsSerializable() bool {
	return l.serializable
}

// SetSerializable sets the serializable flag for the lambda expression
func (l *LambdaExpression) SetSerializable(serializable bool) {
	l.serializable = serializable
}

// Helper function to get parameters text
func getParametersText(parameters []*Parameter) string {
	var paramTexts []string
	for _, param := range parameters {
		paramTexts = append(paramTexts, param.GetText())
	}
	return strings.Join(paramTexts, ", ")
}
