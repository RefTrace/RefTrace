package parser

import (
	"fmt"
)

// FieldExpression represents a field access such as the expression "this.foo".
type FieldExpression struct {
	*BaseExpression
	field  *FieldNode
	useRef bool
}

// NewFieldExpression creates a new FieldExpression
func NewFieldExpression(field *FieldNode) *FieldExpression {
	if field == nil {
		panic("field cannot be nil")
	}
	return &FieldExpression{BaseExpression: NewBaseExpression(), field: field}
}

// Visit implements the GroovyCodeVisitor interface
func (f *FieldExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitFieldExpression(f)
}

// TransformExpression implements the Expression interface
func (f *FieldExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	return f
}

// GetField returns the FieldNode
func (f *FieldExpression) GetField() *FieldNode {
	return f.field
}

// GetFieldName returns the name of the field
func (f *FieldExpression) GetFieldName() string {
	return f.field.GetName()
}

// GetText returns the string representation of the field access
func (f *FieldExpression) GetText() string {
	return "this." + f.GetFieldName()
}

// GetType returns the type of the field
func (f *FieldExpression) GetType() *ClassNode {
	return f.field.GetType()
}

// SetType sets the type of the field
func (f *FieldExpression) SetType(t *ClassNode) {
	f.field.SetType(t)
}

// IsDynamicTyped returns whether the field is dynamically typed
func (f *FieldExpression) IsDynamicTyped() bool {
	return f.field.IsDynamicTyped()
}

// IsUseReferenceDirectly returns the useRef flag
func (f *FieldExpression) IsUseReferenceDirectly() bool {
	return f.useRef
}

// SetUseReferenceDirectly sets the useRef flag
func (f *FieldExpression) SetUseReferenceDirectly(useRef bool) {
	f.useRef = useRef
}

// String returns a string representation of the FieldExpression
func (f *FieldExpression) String() string {
	return fmt.Sprintf("field(%s %s)", f.GetType().GetText(), f.GetFieldName())
}
