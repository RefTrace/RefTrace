package parser

import (
	"fmt"
)

// PropertyNode represents a property (member variable, a getter and setter)
type PropertyNode struct {
	*AnnotatedNode
	field       *FieldNode
	modifiers   int
	getterBlock Statement
	setterBlock Statement
	getterName  string
	setterName  string
}

// NewPropertyNode creates a new PropertyNode with the given field, modifiers, and getter/setter blocks
func NewPropertyNode(field *FieldNode, modifiers int, getterBlock, setterBlock Statement) *PropertyNode {
	return &PropertyNode{
		AnnotatedNode: &AnnotatedNode{},
		field:         field,
		modifiers:     modifiers,
		getterBlock:   getterBlock,
		setterBlock:   setterBlock,
	}
}

// GetField returns the FieldNode of the property
func (p *PropertyNode) GetField() *FieldNode {
	return p.field
}

// SetField sets the FieldNode of the property
func (p *PropertyNode) SetField(field *FieldNode) {
	p.field = field
}

// GetModifiers returns the modifiers of the property
func (p *PropertyNode) GetModifiers() int {
	return p.modifiers
}

// SetModifiers sets the modifiers of the property
func (p *PropertyNode) SetModifiers(modifiers int) {
	p.modifiers = modifiers
}

// GetGetterBlock returns the getter block of the property
func (p *PropertyNode) GetGetterBlock() Statement {
	return p.getterBlock
}

// SetGetterBlock sets the getter block of the property
func (p *PropertyNode) SetGetterBlock(getterBlock Statement) {
	p.getterBlock = getterBlock
}

// GetSetterBlock returns the setter block of the property
func (p *PropertyNode) GetSetterBlock() Statement {
	return p.setterBlock
}

// SetSetterBlock sets the setter block of the property
func (p *PropertyNode) SetSetterBlock(setterBlock Statement) {
	p.setterBlock = setterBlock
}

// GetGetterName returns the getter name of the property
func (p *PropertyNode) GetGetterName() string {
	return p.getterName
}

// GetGetterNameOrDefault returns the getter name or a default name if not set
func (p *PropertyNode) GetGetterNameOrDefault() string {
	if p.getterName != "" {
		return p.getterName
	}
	defaultName := "get" + capitalize(p.GetName())
	// Note: The boolean type check and "is" prefix logic is omitted as it depends on other package structures
	return defaultName
}

// SetGetterName sets the getter name of the property
func (p *PropertyNode) SetGetterName(getterName string) {
	if getterName == "" {
		panic("A non-empty getter name is required")
	}
	p.getterName = getterName
}

// GetSetterName returns the setter name of the property
func (p *PropertyNode) GetSetterName() string {
	return p.setterName
}

// GetSetterNameOrDefault returns the setter name or a default name if not set
func (p *PropertyNode) GetSetterNameOrDefault() string {
	if p.setterName != "" {
		return p.setterName
	}
	return "set" + capitalize(p.GetName())
}

// SetSetterName sets the setter name of the property
func (p *PropertyNode) SetSetterName(setterName string) {
	if setterName == "" {
		panic("A non-empty setter name is required")
	}
	p.setterName = setterName
}

// GetName returns the name of the property
func (p *PropertyNode) GetName() string {
	return p.field.GetName()
}

// GetType returns the type of the property
func (p *PropertyNode) GetType() *ClassNode {
	return p.field.GetType()
}

// GetOriginType returns the origin type of the property
func (p *PropertyNode) GetOriginType() *ClassNode {
	return p.field.GetOriginType()
}

// SetType sets the type of the property
func (p *PropertyNode) SetType(t *ClassNode) {
	p.field.SetType(t)
}

// GetInitialExpression returns the initial expression of the property
func (p *PropertyNode) GetInitialExpression() Expression {
	return p.field.GetInitialExpression()
}

// HasInitialExpression checks if the property has an initial expression
func (p *PropertyNode) HasInitialExpression() bool {
	return p.field.HasInitialExpression()
}

// IsInStaticContext checks if the property is in a static context
func (p *PropertyNode) IsInStaticContext() bool {
	return p.field.IsInStaticContext()
}

// IsDynamicTyped checks if the property is dynamically typed
func (p *PropertyNode) IsDynamicTyped() bool {
	return p.field.IsDynamicTyped()
}

// capitalize is a helper function to capitalize the first letter of a string
func capitalize(s string) string {
	if s == "" {
		return s
	}
	return fmt.Sprintf("%c%s", s[0]-32, s[1:])
}
