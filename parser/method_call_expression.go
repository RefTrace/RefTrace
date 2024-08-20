package parser

import (
	"strings"
)

// MethodCallExpression represents a method call on an object or class.
type MethodCallExpression struct {
	*BaseExpression
	ObjectExpression Expression
	Method           Expression
	Arguments        Expression
	ImplicitThis     bool
	SpreadSafe       bool
	Safe             bool
	GenericsTypes    []*GenericsType
	Target           *MethodNode
}

var NoArguments = &TupleExpression{
	expressions: make([]Expression, 0),
}

func NewMethodCallExpression(objectExpression Expression, method interface{}, arguments Expression) *MethodCallExpression {
	mce := &MethodCallExpression{
		BaseExpression: NewBaseExpression(),
		ImplicitThis:   true,
	}
	mce.SetObjectExpression(objectExpression)

	switch m := method.(type) {
	case string:
		mce.SetMethod(NewConstantExpression(m))
	case Expression:
		mce.SetMethod(m)
	}

	mce.SetArguments(arguments)
	return mce
}

func (mce *MethodCallExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitMethodCallExpression(mce)
}

func (mce *MethodCallExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	answer := NewMethodCallExpression(
		transformer.Transform(mce.ObjectExpression),
		transformer.Transform(mce.Method),
		transformer.Transform(mce.Arguments),
	)
	answer.Safe = mce.Safe
	answer.SpreadSafe = mce.SpreadSafe
	answer.ImplicitThis = mce.ImplicitThis
	answer.GenericsTypes = mce.GenericsTypes
	answer.SetSourcePosition(mce)
	answer.SetMethodTarget(mce.Target)
	answer.CopyNodeMetaData(mce)
	return answer
}

func (mce *MethodCallExpression) GetArguments() Expression {
	return mce.Arguments
}

func (mce *MethodCallExpression) SetArguments(arguments Expression) {
	// make the args a ITupleExpression if they're not already
	if _, ok := arguments.(ITupleExpression); !ok {
		mce.Arguments = NewTupleExpressionWithExpressions(arguments)
		mce.Arguments.SetSourcePosition(arguments)
	} else {
		mce.Arguments = arguments
	}
}

func (mce *MethodCallExpression) GetMethod() Expression {
	return mce.Method
}

func (mce *MethodCallExpression) SetMethod(method Expression) {
	mce.Method = method
}

func (mce *MethodCallExpression) GetMethodAsString() string {
	if ce, ok := mce.Method.(*ConstantExpression); ok {
		return ce.GetText()
	}
	return ""
}

func (mce *MethodCallExpression) GetObjectExpression() Expression {
	return mce.ObjectExpression
}

func (mce *MethodCallExpression) SetObjectExpression(objectExpression Expression) {
	mce.ObjectExpression = objectExpression
}

func (mce *MethodCallExpression) GetReceiver() ASTNode {
	return mce.GetObjectExpression()
}

func (mce *MethodCallExpression) GetText() string {
	var builder strings.Builder
	builder.WriteString(mce.GetObjectExpression().GetText())
	if mce.IsSpreadSafe() {
		builder.WriteRune('*')
	}
	if mce.IsSafe() {
		builder.WriteRune('?')
	}
	builder.WriteRune('.')

	if mce.IsUsingGenerics() {
		builder.WriteRune('<')
		for i, t := range mce.GetGenericsTypes() {
			if i > 0 {
				builder.WriteString(", ")
			}
			builder.WriteString(t.String())
		}
		builder.WriteRune('>')
	}

	method := mce.GetMethod()
	switch m := method.(type) {
	case *GStringExpression:
		builder.WriteString(`"`)
		builder.WriteString(m.GetText())
		builder.WriteString(`"`)
	case *ConstantExpression:
		value := m.GetValue()
		if s, ok := value.(string); ok && IsJavaIdentifier(s) {
			builder.WriteString(s)
		} else {
			builder.WriteString("'")
			builder.WriteString(value.(string))
			builder.WriteString("'")
		}
	default:
		builder.WriteString("(")
		builder.WriteString(m.GetText())
		builder.WriteString(")")
	}

	builder.WriteString(mce.GetArguments().GetText())

	return builder.String()
}

func (mce *MethodCallExpression) IsSafe() bool {
	return mce.Safe
}

func (mce *MethodCallExpression) SetSafe(safe bool) {
	mce.Safe = safe
}

func (mce *MethodCallExpression) IsSpreadSafe() bool {
	return mce.SpreadSafe
}

func (mce *MethodCallExpression) SetSpreadSafe(value bool) {
	mce.SpreadSafe = value
}

func (mce *MethodCallExpression) IsImplicitThis() bool {
	return mce.ImplicitThis
}

func (mce *MethodCallExpression) SetImplicitThis(implicitThis bool) {
	mce.ImplicitThis = implicitThis
}

func (mce *MethodCallExpression) GetGenericsTypes() []*GenericsType {
	return mce.GenericsTypes
}

func (mce *MethodCallExpression) SetGenericsTypes(genericsTypes []*GenericsType) {
	mce.GenericsTypes = genericsTypes
}

func (mce *MethodCallExpression) IsUsingGenerics() bool {
	return mce.GenericsTypes != nil && len(mce.GenericsTypes) > 0
}

func (mce *MethodCallExpression) GetMethodTarget() *MethodNode {
	return mce.Target
}

func (mce *MethodCallExpression) SetMethodTarget(mn *MethodNode) {
	mce.Target = mn
	if mn != nil {
		mce.SetType(mn.returnType)
	} else {
		mce.SetType(OBJECT_TYPE)
	}
}

type MethodCall interface {
	GetReceiver() ASTNode
	GetMethodAsString() string
	GetArguments() Expression
	GetText() string
}

func (mce *MethodCallExpression) SetSourcePosition(node ASTNode) {
	mce.BaseExpression.SetSourcePosition(node)

	switch n := node.(type) {
	case MethodCall:
		if mce, ok := n.(*MethodCallExpression); ok {
			mce.Method.SetSourcePosition(mce.GetMethod())
		} else if node.GetLineNumber() > 0 {
			mce.Method.SetLineNumber(node.GetLineNumber())
			mce.Method.SetColumnNumber(node.GetColumnNumber())
			mce.Method.SetLastLineNumber(node.GetLineNumber())
			mce.Method.SetLastColumnNumber(node.GetColumnNumber() + len(mce.GetMethodAsString()))
		}
		if mce.Arguments != nil {
			mce.Arguments.SetSourcePosition(n.GetArguments())
		}
	case *PropertyExpression:
		mce.Method.SetSourcePosition(n.GetProperty())
	}
}

func (mce *MethodCallExpression) String() string {
	return mce.BaseExpression.GetText() + "[object: " + mce.ObjectExpression.GetText() +
		" method: " + mce.Method.GetText() + " arguments: " + mce.Arguments.GetText() + "]"
}
