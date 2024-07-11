package parser

import (
	"fmt"
)

// CastExpression represents a typecast expression.
type CastExpression struct {
	Expression
	expression       Expression
	ignoreAutoboxing bool
	coerce           bool
	strict           bool
}

// AsExpression creates a new CastExpression with coerce set to true.
func AsExpression(typ *ClassNode, expression Expression) *CastExpression {
	answer := NewCastExpression(typ, expression)
	answer.SetCoerce(true)
	return answer
}

// NewCastExpression creates a new CastExpression.
func NewCastExpression(typ *ClassNode, expression Expression) *CastExpression {
	return NewCastExpressionWithAutoboxing(typ, expression, false)
}

// NewCastExpressionWithAutoboxing creates a new CastExpression with autoboxing option.
func NewCastExpressionWithAutoboxing(typ *ClassNode, expression Expression, ignoreAutoboxing bool) *CastExpression {
	ce := &CastExpression{
		expression:       expression,
		ignoreAutoboxing: ignoreAutoboxing,
	}
	ce.SetType(typ)
	return ce
}

// GetExpression returns the expression being cast.
func (ce *CastExpression) GetExpression() Expression {
	return ce.expression
}

// SetExpression sets the expression being cast.
func (ce *CastExpression) SetExpression(expression Expression) {
	ce.expression = expression
}

// IsIgnoringAutoboxing returns whether autoboxing is ignored.
func (ce *CastExpression) IsIgnoringAutoboxing() bool {
	return ce.ignoreAutoboxing
}

// IsCoerce returns whether this is a coercion cast.
func (ce *CastExpression) IsCoerce() bool {
	return ce.coerce
}

// SetCoerce sets whether this is a coercion cast.
func (ce *CastExpression) SetCoerce(coerce bool) {
	ce.coerce = coerce
}

// IsStrict returns whether this is a strict cast.
func (ce *CastExpression) IsStrict() bool {
	return ce.strict
}

// SetStrict sets whether this is a strict cast.
func (ce *CastExpression) SetStrict(strict bool) {
	ce.strict = strict
}

// String returns a string representation of the CastExpression.
func (ce *CastExpression) String() string {
	return fmt.Sprintf("%s[%s]", ce.Expression.String(), ce.GetText())
}

// Visit calls the appropriate visit method on the GroovyCodeVisitor.
func (ce *CastExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitCastExpression(ce)
}

// TransformExpression transforms the CastExpression.
func (ce *CastExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewCastExpressionWithAutoboxing(ce.GetType(), transformer.Transform(ce.expression), ce.IsIgnoringAutoboxing())
	ret.SetCoerce(ce.IsCoerce())
	ret.SetStrict(ce.IsStrict())
	ret.SetSourcePosition(ce)
	ret.CopyNodeMetaData(ce)
	return ret
}

// GetText returns the text representation of the CastExpression.
func (ce *CastExpression) GetText() string {
	if ce.IsCoerce() {
		return fmt.Sprintf("%s as %s", ce.expression.GetText(), ce.GetType().ToString(false))
	}
	return fmt.Sprintf("(%s) %s", ce.GetType().ToString(false), ce.expression.GetText())
}

// SetType is not supported for CastExpression.
func (ce *CastExpression) SetType(typ *ClassNode) {
	panic("SetType is not supported for CastExpression")
}
