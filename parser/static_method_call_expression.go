package parser

import (
	"fmt"
)

var _ MethodCall = (*StaticMethodCallExpression)(nil)

// StaticMethodCallExpression represents a static method call on a class
type StaticMethodCallExpression struct {
	*BaseExpression
	ownerType *ClassNode
	method    string
	arguments Expression
}

func NewStaticMethodCallExpression(expr Expression, ownerType *ClassNode, method string, arguments Expression) *StaticMethodCallExpression {
	return &StaticMethodCallExpression{
		BaseExpression: NewBaseExpression(),
		ownerType:      ownerType,
		method:         method,
		arguments:      arguments,
	}
}

func (s *StaticMethodCallExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitStaticMethodCallExpression(s)
}

/*
func (s *StaticMethodCallExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewStaticMethodCallExpression(s.GetOwnerType(), s.method, transformer.Transform(s.arguments))
	ret.SetSourcePosition(s)
	ret.CopyNodeMetaData(s)
	return ret
}
*/

func (s *StaticMethodCallExpression) GetReceiver() ASTNode {
	return s.ownerType
}

func (s *StaticMethodCallExpression) GetMethodAsString() string {
	return s.method
}

func (s *StaticMethodCallExpression) GetArguments() Expression {
	return s.arguments
}

func (s *StaticMethodCallExpression) GetMethod() string {
	return s.method
}

func (s *StaticMethodCallExpression) GetText() string {
	return fmt.Sprintf("%s.%s%s", s.GetOwnerType().GetName(), s.method, s.arguments.GetText())
}

func (s *StaticMethodCallExpression) String() string {
	return fmt.Sprintf("%s[%s#%s arguments: %v]", s.BaseExpression.GetText(), s.GetOwnerType().GetName(), s.method, s.arguments)
}

func (s *StaticMethodCallExpression) GetOwnerType() *ClassNode {
	return s.ownerType
}

func (s *StaticMethodCallExpression) SetOwnerType(ownerType *ClassNode) {
	s.ownerType = ownerType
}
