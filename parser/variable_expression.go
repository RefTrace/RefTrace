package parser

import (
	"fmt"
)

// VariableExpression represents a local variable name, the simplest form of expression. e.g. "foo".
type VariableExpression struct {
	Expression
	Variable
	variable         string
	modifiers        int
	inStaticContext  bool
	isDynamicTyped   bool
	accessedVariable Variable
	closureShare     bool
	useRef           bool
	originType       *ClassNode
}

var (
	THIS_EXPRESSION  = NewVariableExpression("this", DynamicType())
	SUPER_EXPRESSION = NewVariableExpression("super", DynamicType())
)

func NewVariableExpression(name string, typ *ClassNode) *VariableExpression {
	ve := &VariableExpression{
		variable:   name,
		originType: typ,
	}
	if IsPrimitiveType(typ) {
		ve.SetType(GetWrapper(typ))
	} else {
		ve.SetType(typ)
	}
	return ve
}

func NewVariableExpressionWithVariable(variable Variable) *VariableExpression {
	ve := NewVariableExpression(variable.GetName(), variable.GetOriginType())
	ve.SetAccessedVariable(variable)
	ve.SetModifiers(variable.GetModifiers())
	return ve
}

func (ve *VariableExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitVariableExpression(ve)
}

func (ve *VariableExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	return ve
}

func (ve *VariableExpression) GetText() string {
	return ve.variable
}

func (ve *VariableExpression) GetName() string {
	return ve.variable
}

func (ve *VariableExpression) String() string {
	typeStr := ""
	if !ve.IsDynamicTyped() {
		typeStr = fmt.Sprintf(" type: %v", ve.GetType())
	}
	return fmt.Sprintf("%v[variable: %s%s]", ve.Expression.String(), ve.variable, typeStr)
}

func (ve *VariableExpression) GetInitialExpression() Expression {
	return nil
}

func (ve *VariableExpression) HasInitialExpression() bool {
	return false
}

func (ve *VariableExpression) IsInStaticContext() bool {
	if ve.accessedVariable != nil && ve.accessedVariable != ve {
		return ve.accessedVariable.IsInStaticContext()
	}
	return ve.inStaticContext
}

func (ve *VariableExpression) SetInStaticContext(inStaticContext bool) {
	ve.inStaticContext = inStaticContext
}

func (ve *VariableExpression) SetType(cn *ClassNode) {
	ve.Expression.SetType(cn)
	ve.isDynamicTyped = ve.isDynamicTyped || IsDynamicTyped(cn)
}

func (ve *VariableExpression) IsDynamicTyped() bool {
	if ve.accessedVariable != nil && ve.accessedVariable != ve {
		return ve.accessedVariable.IsDynamicTyped()
	}
	return ve.isDynamicTyped
}

func (ve *VariableExpression) IsClosureSharedVariable() bool {
	if ve.accessedVariable != nil && ve.accessedVariable != ve {
		return ve.accessedVariable.IsClosureSharedVariable()
	}
	return ve.closureShare
}

func (ve *VariableExpression) SetClosureSharedVariable(inClosure bool) {
	ve.closureShare = inClosure
}

func (ve *VariableExpression) GetModifiers() int {
	return ve.modifiers
}

func (ve *VariableExpression) SetUseReferenceDirectly(useRef bool) {
	ve.useRef = useRef
}

func (ve *VariableExpression) IsUseReferenceDirectly() bool {
	return ve.useRef
}

func (ve *VariableExpression) GetType() *ClassNode {
	if ve.accessedVariable != nil && ve.accessedVariable != ve {
		return ve.accessedVariable.GetType()
	}
	return ve.Expression.GetType()
}

func (ve *VariableExpression) GetOriginType() *ClassNode {
	if ve.accessedVariable != nil && ve.accessedVariable != ve {
		return ve.accessedVariable.GetOriginType()
	}
	return ve.originType
}

func (ve *VariableExpression) IsThisExpression() bool {
	return ve.variable == "this"
}

func (ve *VariableExpression) IsSuperExpression() bool {
	return ve.variable == "super"
}

func (ve *VariableExpression) SetModifiers(modifiers int) {
	ve.modifiers = modifiers
}

func (ve *VariableExpression) GetAccessedVariable() Variable {
	return ve.accessedVariable
}

func (ve *VariableExpression) SetAccessedVariable(origin Variable) {
	ve.accessedVariable = origin
}
