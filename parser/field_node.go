package parser

import (
	"errors"
	"fmt"
	"reflect"
)

type FieldNode struct {
	*AnnotatedNode
	Variable
	name                   string
	modifiers              int
	fieldType              *ClassNode
	owner                  *ClassNode
	initialValueExpression Expression
	dynamicTyped           bool
	holder                 bool
	originType             *ClassNode
}

func NewStatic(theClass reflect.Type, name string) (*FieldNode, error) {
	field, found := theClass.FieldByName(name)
	if !found {
		return nil, errors.New(fmt.Sprintf("%s field not found", name))
	}
	fldType := Make(field.Type)
	return &FieldNode{
		name:      name,
		modifiers: ACC_PUBLIC | ACC_STATIC,
		fieldType: fldType,
		owner:     Make(theClass),
	}, nil
}

func NewFieldNode(name string, modifiers int, fieldType, owner *ClassNode, initialValueExpression Expression) *FieldNode {
	f := &FieldNode{
		name:                   name,
		modifiers:              modifiers,
		owner:                  owner,
		initialValueExpression: initialValueExpression,
	}
	f.SetType(fieldType)
	return f
}

func (f *FieldNode) GetInitialExpression() Expression {
	return f.initialValueExpression
}

func (f *FieldNode) GetName() string {
	return f.name
}

func (f *FieldNode) GetType() *ClassNode {
	return f.fieldType
}

func (f *FieldNode) SetType(fieldType *ClassNode) {
	f.fieldType = fieldType
	f.originType = fieldType
	f.dynamicTyped = f.dynamicTyped || IsDynamicTyped(fieldType)
}

func (f *FieldNode) GetOwner() *ClassNode {
	return f.owner
}

func (f *FieldNode) IsHolder() bool {
	return f.holder
}

func (f *FieldNode) SetHolder(holder bool) {
	f.holder = holder
}

func (f *FieldNode) IsDynamicTyped() bool {
	return f.dynamicTyped
}

func (f *FieldNode) GetModifiers() int {
	return f.modifiers
}

func (f *FieldNode) SetModifiers(modifiers int) {
	f.modifiers = modifiers
}

func (f *FieldNode) IsEnum() bool {
	return (f.GetModifiers() & ACC_ENUM) != 0
}

func (f *FieldNode) SetOwner(owner *ClassNode) {
	f.owner = owner
}

func (f *FieldNode) HasInitialExpression() bool {
	return f.initialValueExpression != nil
}

func (f *FieldNode) IsInStaticContext() bool {
	return f.IsStatic()
}

func (f *FieldNode) GetInitialValueExpression() Expression {
	return f.initialValueExpression
}

func (f *FieldNode) SetInitialValueExpression(initialValueExpression Expression) {
	f.initialValueExpression = initialValueExpression
}

func (f *FieldNode) Equals(obj interface{}) bool {
	if obj != nil && reflect.TypeOf(obj).String() == "org.codehaus.groovy.ast.decompiled.LazyFieldNode" {
		return obj.(interface{ Equals(*FieldNode) bool }).Equals(f)
	}
	return f.AnnotatedNode.declaringClass.Equals(obj)
}

func (f *FieldNode) GetOriginType() *ClassNode {
	return f.originType
}

func (f *FieldNode) SetOriginType(cn *ClassNode) {
	f.originType = cn
}

func (f *FieldNode) Rename(name string) {
	f.GetDeclaringClass().RenameField(f.name, name)
	f.name = name
}
