package parser

import (
	"fmt"
)

// MethodReferenceExpression represents a method reference or a constructor reference,
// e.g. System.out::println OR Objects::requireNonNull OR Integer::new OR int[]::new
type MethodReferenceExpression struct {
	MethodPointerExpression
}

func NewMethodReferenceExpression(expression, methodName Expression) *MethodReferenceExpression {
	return &MethodReferenceExpression{
		MethodPointerExpression: MethodPointerExpression{
			Expression: expression,
			MethodName: methodName,
		},
	}
}

func (m *MethodReferenceExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitMethodReferenceExpression(m)
}

func (m *MethodReferenceExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	var ret Expression
	mname := transformer.Transform(m.MethodName)
	if m.Expression == nil {
		ret = NewMethodReferenceExpression(THIS_EXPRESSION, mname)
	} else {
		ret = NewMethodReferenceExpression(transformer.Transform(m.Expression), mname)
	}
	ret.SetSourcePosition(m)
	ret.CopyNodeMetaData(m)
	return ret
}

func (m *MethodReferenceExpression) GetText() string {
	expression := m.GetExpression()
	methodName := m.GetMethodName()

	if expression == nil {
		return "::" + methodName.GetText()
	} else {
		return fmt.Sprintf("%s::%s", expression.GetText(), methodName.GetText())
	}
}
