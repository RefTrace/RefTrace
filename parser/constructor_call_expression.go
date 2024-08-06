package parser

import (
	"fmt"
)

// ConstructorCallExpression represents a constructor call.
type ConstructorCallExpression struct {
	*BaseExpression
	arguments               Expression
	usesAnonymousInnerClass bool
}

// NewConstructorCallExpression creates a new ConstructorCallExpression.
func NewConstructorCallExpression(typ *ClassNode, arguments Expression) *ConstructorCallExpression {
	cce := &ConstructorCallExpression{BaseExpression: NewBaseExpression()}
	cce.SetType(typ)

	if _, ok := arguments.(*TupleExpression); !ok {
		cce.arguments = NewTupleExpressionWithExpressions(arguments)
		cce.arguments.SetSourcePosition(arguments)
	} else {
		cce.arguments = arguments
	}

	return cce
}

// Visit implements the GroovyCodeVisitor interface.
func (cce *ConstructorCallExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitConstructorCallExpression(cce)
}

// TransformExpression transforms the expression.
func (cce *ConstructorCallExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	answer := NewConstructorCallExpression(cce.GetType(), transformer.Transform(cce.arguments))
	answer.SetUsingAnonymousInnerClass(cce.IsUsingAnonymousInnerClass())
	answer.SetSourcePosition(cce)
	answer.CopyNodeMetaData(cce)
	return answer
}

// GetReceiver implements the MethodCall interface.
func (cce *ConstructorCallExpression) GetReceiver() ASTNode {
	return nil
}

// GetMethodAsString implements the MethodCall interface.
func (cce *ConstructorCallExpression) GetMethodAsString() string {
	return "<init>"
}

// GetArguments implements the MethodCall interface.
func (cce *ConstructorCallExpression) GetArguments() Expression {
	return cce.arguments
}

// GetText returns the text representation of the expression.
func (cce *ConstructorCallExpression) GetText() string {
	var text string
	if cce.IsSuperCall() {
		text = "super "
	} else if cce.IsThisCall() {
		text = "this "
	} else {
		text = "new " + cce.GetType().GetText()
	}
	return text + cce.GetArguments().GetText()
}

// IsSpecialCall checks if it's a special call (this or super).
func (cce *ConstructorCallExpression) IsSpecialCall() bool {
	return cce.IsThisCall() || cce.IsSuperCall()
}

// IsSuperCall checks if it's a super call.
func (cce *ConstructorCallExpression) IsSuperCall() bool {
	return cce.GetType() == SUPER
}

// IsThisCall checks if it's a this call.
func (cce *ConstructorCallExpression) IsThisCall() bool {
	return cce.GetType() == THIS
}

// IsUsingAnonymousInnerClass checks if it's using an anonymous inner class.
func (cce *ConstructorCallExpression) IsUsingAnonymousInnerClass() bool {
	return cce.usesAnonymousInnerClass
}

// SetUsingAnonymousInnerClass sets the usage of anonymous inner class.
func (cce *ConstructorCallExpression) SetUsingAnonymousInnerClass(usage bool) {
	cce.usesAnonymousInnerClass = usage
}

// String returns a string representation of the expression.
func (cce *ConstructorCallExpression) String() string {
	return fmt.Sprintf("%s[type: %v arguments: %v]", cce.BaseExpression.GetText(), cce.GetType(), cce.arguments)
}
