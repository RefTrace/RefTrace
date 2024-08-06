package parser

import (
	"reflect"
)

// MethodPointerExpression represents a method pointer on an object such as
// foo.&bar which means find the method pointer for the bar method on the foo instance.
// This is equivalent to:
// foo.metaClass.getMethodPointer(foo, "bar")
type MethodPointerExpression struct {
	*BaseExpression
	expression Expression
	methodName Expression
}

func NewMethodPointerExpression(expression, methodName Expression) *MethodPointerExpression {
	mpe := &MethodPointerExpression{
		BaseExpression: NewBaseExpression(),
		expression:     expression,
		methodName:     methodName,
	}
	mpe.SetType(CLOSURE_TYPE.GetPlainNodeReference())
	return mpe
}

func (m *MethodPointerExpression) GetExpression() Expression {
	if m.expression == nil {
		return THIS_EXPRESSION
	}
	return m.expression
}

func (m *MethodPointerExpression) GetMethodName() Expression {
	return m.methodName
}

func (m *MethodPointerExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitMethodPointerExpression(m)
}

func (m *MethodPointerExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	mname := transformer.Transform(m.methodName)
	var ret Expression
	if m.expression == nil {
		ret = NewMethodPointerExpression(THIS_EXPRESSION, mname)
	} else {
		ret = NewMethodPointerExpression(transformer.Transform(m.expression), mname)
	}
	ret.SetSourcePosition(m)
	ret.CopyNodeMetaData(m)
	return ret
}

func (m *MethodPointerExpression) GetText() string {
	if m.expression == nil {
		return "&" + m.methodName.GetText()
	}
	return m.expression.GetText() + ".&" + m.methodName.GetText()
}

func (m *MethodPointerExpression) IsDynamic() bool {
	return false
}

func (m *MethodPointerExpression) GetTypeClass() reflect.Type {
	return reflect.TypeOf((*func())(nil)).Elem()
}
