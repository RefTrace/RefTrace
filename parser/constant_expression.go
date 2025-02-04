package parser

import (
	"fmt"
	"reflect"
)

// Compile-time assertion that ConstantExpression implements the Expression interface
var _ Expression = (*ConstantExpression)(nil)

// ConstantExpression represents a constant expression such as null, true, false.
type ConstantExpression struct {
	*BaseExpression
	value        interface{}
	constantName string
}

// Predefined constant expressions
var (
	NULL             = NewConstantExpression(nil)
	TRUE             = NewConstantExpression(true)
	FALSE            = NewConstantExpression(false)
	EMPTY_STRING     = NewConstantExpression("")
	PRIM_TRUE        = NewConstantExpressionPrimitive(true)
	PRIM_FALSE       = NewConstantExpressionPrimitive(false)
	VOID             = NewConstantExpression(reflect.TypeOf((*interface{})(nil)).Elem())
	EMPTY_EXPRESSION = NewEmptyExpression()
)

// NewConstantExpression creates a new ConstantExpression
func NewConstantExpression(value interface{}) *ConstantExpression {
	return NewConstantExpressionPrimitive(value)
}

// NewConstantExpressionPrimitive creates a new ConstantExpression with primitive type preservation
func NewConstantExpressionPrimitive(value interface{}) *ConstantExpression {
	expr := &ConstantExpression{BaseExpression: NewBaseExpression(), value: value}
	if value != nil {
		expr.setTypeFromValue(value, true)
	}
	return expr
}

func (c *ConstantExpression) setTypeFromValue(value interface{}, keepPrimitive bool) {
	if keepPrimitive {
		switch value.(type) {
		case int:
			c.SetType(INT_TYPE)
		case int64:
			c.SetType(LONG_TYPE)
		case bool:
			c.SetType(BOOLEAN_TYPE)
		case float64:
			c.SetType(DOUBLE_TYPE)
		case float32:
			c.SetType(FLOAT_TYPE)
		case rune:
			c.SetType(CHAR_TYPE)
		default:
			c.SetType(Make(reflect.TypeOf(value)))
		}
	} else {
		c.SetType(Make(reflect.TypeOf(value)))
	}
}

func (c *ConstantExpression) String() string {
	return fmt.Sprintf("%s[%v]", c.BaseExpression.GetText(), c.value)
}

func (c *ConstantExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitConstantExpression(c)
}

// TODO: implement
/*
func (c *ConstantExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	return c
}
*/

func (c *ConstantExpression) GetValue() interface{} {
	return c.value
}

func (c *ConstantExpression) GetText() string {
	if c.value == nil {
		return "null"
	}
	return fmt.Sprintf("%v", c.value)
}

func (c *ConstantExpression) GetConstantName() string {
	return c.constantName
}

func (c *ConstantExpression) SetConstantName(constantName string) {
	c.constantName = constantName
}

func (c *ConstantExpression) IsNullExpression() bool {
	return c.value == nil
}

func (c *ConstantExpression) IsTrueExpression() bool {
	return c.value == true
}

func (c *ConstantExpression) IsFalseExpression() bool {
	return c.value == false
}

func (c *ConstantExpression) IsEmptyStringExpression() bool {
	return c.value == ""
}
