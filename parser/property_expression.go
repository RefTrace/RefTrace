package parser

import (
	"fmt"
	"strings"
)

// PropertyExpression represents a property access such as the expression "foo.bar".
type PropertyExpression struct {
	*BaseExpression
	objectExpression Expression
	property         Expression
	safe             bool
	spreadSafe       bool
	isStatic         bool
	implicitThis     bool
}

func NewPropertyExpression(objectExpression Expression, propertyName string) *PropertyExpression {
	return NewPropertyExpressionWithSafe(objectExpression, NewConstantExpression(propertyName), false)
}

func NewPropertyExpressionWithProperty(objectExpression, property Expression) *PropertyExpression {
	return NewPropertyExpressionWithSafe(objectExpression, property, false)
}

func NewPropertyExpressionWithSafe(objectExpression, property Expression, safe bool) *PropertyExpression {
	return &PropertyExpression{
		BaseExpression:   NewBaseExpression(),
		objectExpression: objectExpression,
		property:         property,
		safe:             safe,
	}
}

func (p *PropertyExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewPropertyExpressionWithSafe(
		transformer.Transform(p.GetObjectExpression()),
		transformer.Transform(p.GetProperty()),
		p.IsSafe(),
	)
	ret.SetImplicitThis(p.IsImplicitThis())
	ret.SetSpreadSafe(p.IsSpreadSafe())
	ret.SetStatic(p.IsStatic())
	ret.SetType(p.GetType())
	ret.SetSourcePosition(p)
	ret.CopyNodeMetaData(p)
	return ret
}

func (p *PropertyExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitPropertyExpression(p)
}

func (p *PropertyExpression) GetObjectExpression() Expression {
	return p.objectExpression
}

func (p *PropertyExpression) SetObjectExpression(objectExpression Expression) {
	p.objectExpression = objectExpression
}

func (p *PropertyExpression) GetProperty() Expression {
	return p.property
}

func (p *PropertyExpression) GetPropertyAsString() string {
	if constExpr, ok := p.GetProperty().(*ConstantExpression); ok {
		return constExpr.GetText()
	}
	return ""
}

func (p *PropertyExpression) GetText() string {
	var sb strings.Builder
	sb.WriteString(p.GetObjectExpression().GetText())
	if p.IsSpreadSafe() {
		sb.WriteRune('*')
	}
	if p.IsSafe() {
		sb.WriteRune('?')
	}
	sb.WriteRune('.')
	sb.WriteString(p.GetProperty().GetText())
	return sb.String()
}

func (p *PropertyExpression) IsDynamic() bool {
	return true
}

func (p *PropertyExpression) IsImplicitThis() bool {
	return p.implicitThis
}

func (p *PropertyExpression) SetImplicitThis(implicitThis bool) {
	p.implicitThis = implicitThis
}

func (p *PropertyExpression) IsSafe() bool {
	return p.safe
}

func (p *PropertyExpression) IsSpreadSafe() bool {
	return p.spreadSafe
}

func (p *PropertyExpression) SetSpreadSafe(spreadSafe bool) {
	p.spreadSafe = spreadSafe
}

func (p *PropertyExpression) IsStatic() bool {
	return p.isStatic
}

func (p *PropertyExpression) SetStatic(isStatic bool) {
	p.isStatic = isStatic
}

func (p *PropertyExpression) String() string {
	return fmt.Sprintf("%s[object: %v property: %v]", p.BaseExpression.GetText(), p.GetObjectExpression(), p.GetProperty())
}
