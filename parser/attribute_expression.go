package parser

import (
	"strings"
)

// AttributeExpression represents an attribute access (accessing the field of a class)
// such as the expression "foo.@bar".
type AttributeExpression struct {
	PropertyExpression
}

func NewAttributeExpression(objectExpression, property Expression) *AttributeExpression {
	return &AttributeExpression{
		PropertyExpression{
			objectExpression: objectExpression,
			property:         property,
			safe:             false,
		},
	}
}

func NewAttributeExpressionWithSafe(objectExpression, property Expression, safe bool) *AttributeExpression {
	return &AttributeExpression{
		PropertyExpression{
			objectExpression: objectExpression,
			property:         property,
			safe:             safe,
		},
	}
}

func (a *AttributeExpression) GetText() string {
	var sb strings.Builder
	sb.WriteString(a.objectExpression.GetText())
	if a.IsSpreadSafe() {
		sb.WriteRune('*')
	}
	if a.IsSafe() {
		sb.WriteRune('?')
	}
	sb.WriteString(".@")
	sb.WriteString(a.property.GetText())
	return sb.String()
}

func (a *AttributeExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewAttributeExpressionWithSafe(
		transformer.Transform(a.objectExpression),
		transformer.Transform(a.property),
		a.IsSafe(),
	)
	ret.SetImplicitThis(a.IsImplicitThis())
	ret.SetSpreadSafe(a.IsSpreadSafe())
	ret.SetStatic(a.IsStatic())
	ret.SetType(a.GetType())
	ret.SetSourcePosition(a)
	ret.CopyNodeMetaData(a)
	return ret
}

func (a *AttributeExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitAttributeExpression(a)
}
