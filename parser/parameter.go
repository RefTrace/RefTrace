package parser

import (
	"fmt"
)

type Parameter struct {
	AnnotatedNode
	Variable
	DefaultVariable
	paramType       *ClassNode
	name            string
	originType      *ClassNode
	dynamicTyped    bool
	closureShare    bool
	defaultValue    Expression
	hasDefaultValue bool
	inStaticContext bool
	modifiers       int
}

func NewParameter(paramType *ClassNode, name string) *Parameter {
	p := &Parameter{
		name:       name,
		paramType:  paramType,
		originType: paramType,
	}
	p.setType(paramType)
	return p
}

func NewParameterWithDefault(paramType *ClassNode, name string, defaultValue Expression) *Parameter {
	p := NewParameter(paramType, name)
	p.SetInitialExpression(defaultValue)
	return p
}

func (p *Parameter) String() string {
	typeStr := ""
	if p.paramType != nil {
		typeStr = fmt.Sprintf(", type: %s", p.paramType.GetText())
	}
	return fmt.Sprintf("%s[name: %s%s, hasDefaultValue: %t]", p.AnnotatedNode.GetText(), p.name, typeStr, p.HasInitialExpression())
}

func (p Parameter) Name() string {
	return p.name
}

func (p Parameter) Type() *ClassNode {
	return p.paramType
}

func (p Parameter) Modifiers() int {
	return p.modifiers
}

func (p *Parameter) setType(paramType *ClassNode) {
	p.paramType = paramType
	p.dynamicTyped = p.dynamicTyped || isDynamicTyped(paramType)
}

func (p Parameter) GetDefaultValue() Expression {
	return p.defaultValue
}

func (p Parameter) HasInitialExpression() bool {
	return p.hasDefaultValue
}

func (p Parameter) GetInitialExpression() Expression {
	return p.defaultValue
}

func (p *Parameter) SetInitialExpression(init Expression) {
	p.defaultValue = init
	p.hasDefaultValue = init != nil
}

func (p Parameter) IsInStaticContext() bool {
	return p.inStaticContext
}

func (p *Parameter) SetInStaticContext(inStaticContext bool) {
	p.inStaticContext = inStaticContext
}

func (p Parameter) IsDynamicTyped() bool {
	return p.dynamicTyped
}

func (p Parameter) IsClosureSharedVariable() bool {
	return p.closureShare
}

func (p *Parameter) SetClosureSharedVariable(inClosure bool) {
	p.closureShare = inClosure
}

func (p Parameter) GetModifiers() int {
	return p.modifiers
}

func (p *Parameter) SetModifiers(modifiers int) {
	p.modifiers = modifiers
}

func (p Parameter) OriginType() *ClassNode {
	return p.originType
}

func (p *Parameter) SetOriginType(cn *ClassNode) {
	p.originType = cn
}

func (p Parameter) IsImplicit() bool {
	return (p.GetModifiers() & ACC_MANDATED) != 0
}

func (p Parameter) IsReceiver() bool {
	return p.name == "this"
}

// Helper functions

func isDynamicTyped(cn *ClassNode) bool {
	// Implement this function based on your ClassHelper.isDynamicTyped logic
	return false
}

// Add this method to explicitly implement the IsFinal method for Parameter
func (p Parameter) IsFinal() bool {
	return (p.Modifiers() & ACC_FINAL) != 0
}

func (p Parameter) IsPrivate() bool {
	return (p.Modifiers() & ACC_PRIVATE) != 0
}

func (p Parameter) IsProtected() bool {
	return (p.Modifiers() & ACC_PROTECTED) != 0
}

func (p Parameter) IsPublic() bool {
	return (p.Modifiers() & ACC_PUBLIC) != 0
}

func (p Parameter) IsStatic() bool {
	return (p.Modifiers() & ACC_STATIC) != 0
}

func (p Parameter) IsVolatile() bool {
	return (p.Modifiers() & ACC_VOLATILE) != 0
}
