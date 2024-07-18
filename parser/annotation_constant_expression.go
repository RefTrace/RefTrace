package parser

import (
	"fmt"
)

// AnnotationConstantExpression represents an annotation "constant" that may appear in annotation attributes
// (mainly used as a marker).
type AnnotationConstantExpression struct {
	*ConstantExpression
}

// NewAnnotationConstantExpression creates a new AnnotationConstantExpression
func NewAnnotationConstantExpression(node *AnnotationNode) *AnnotationConstantExpression {
	ace := &AnnotationConstantExpression{
		ConstantExpression: NewConstantExpression(node),
	}
	ace.SetType(node.GetClassNode())
	return ace
}

// Visit implements the visitor pattern
func (a *AnnotationConstantExpression) Visit(visitor GroovyCodeVisitor) {
	a.ConstantExpression.Visit(visitor) // GROOVY-9980

	node := a.GetValue().(*AnnotationNode)
	attrs := node.GetMembers()
	for _, expr := range attrs {
		expr.Visit(visitor)
	}
}

// String returns a string representation of the AnnotationConstantExpression
func (a *AnnotationConstantExpression) String() string {
	return fmt.Sprintf("%s[%v]", a.ConstantExpression.String(), a.GetValue())
}
