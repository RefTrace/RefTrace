package parser

import (
	"fmt"
	"strings"
)

// ListExpression represents a list expression [1, 2, 3] which creates a mutable List
type ListExpression struct {
	Expression
	expressions []Expression
	wrapped     bool
}

// NewListExpression creates a new ListExpression
func NewListExpression() *ListExpression {
	return &ListExpression{
		expressions: make([]Expression, 0),
		wrapped:     false,
	}
}

// NewListExpressionWithExpressions creates a new ListExpression with initial expressions
func NewListExpressionWithExpressions(expressions []Expression) *ListExpression {
	le := &ListExpression{
		expressions: expressions,
		wrapped:     false,
	}
	// TODO: get the types of the expressions to specify the
	// list type to List<X> if possible.
	le.SetType(LIST_TYPE)
	return le
}

// AddExpression adds an expression to the list
func (le *ListExpression) AddExpression(expression Expression) {
	le.expressions = append(le.expressions, expression)
}

// GetExpressions returns the list of expressions
func (le *ListExpression) GetExpressions() []Expression {
	return le.expressions
}

// SetWrapped sets the wrapped flag
func (le *ListExpression) SetWrapped(value bool) {
	le.wrapped = value
}

// IsWrapped returns the wrapped flag
func (le *ListExpression) IsWrapped() bool {
	return le.wrapped
}

// Visit implements the GroovyCodeVisitor interface
/*
func (le *ListExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitListExpression(le)
}
*/

// TransformExpression transforms the ListExpression
// TODO: implement
/*
func (le *ListExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewListExpressionWithExpressions(TransformExpressions(le.GetExpressions(), transformer))
	ret.SetSourcePosition(le)
	ret.CopyNodeMetaData(le)
	return ret
}
*/

// GetExpression returns the expression at the given index
func (le *ListExpression) GetExpression(i int) Expression {
	return le.expressions[i]
}

// GetText returns the string representation of the ListExpression
func (le *ListExpression) GetText() string {
	var buffer strings.Builder
	buffer.WriteString("[")
	for i, expression := range le.expressions {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(expression.GetText())
	}
	buffer.WriteString("]")
	return buffer.String()
}

// String returns a string representation of the ListExpression
func (le *ListExpression) String() string {
	return fmt.Sprintf("ListExpression%v", le.expressions)
}
