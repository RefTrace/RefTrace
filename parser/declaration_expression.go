package parser

import (
	"fmt"
	"strings"
)

var _ Expression = (*DeclarationExpression)(nil)

// DeclarationExpression represents one or more local variables.
type DeclarationExpression struct {
	*BinaryExpression
}

// NewDeclarationExpression creates a new DeclarationExpression.
func NewDeclarationExpression(left Expression, operation *Token, right Expression) *DeclarationExpression {
	check(left)
	return &DeclarationExpression{
		BinaryExpression: NewBinaryExpression(
			left,
			NewToken(ASSIGN, "=", operation.startLine, operation.startColumn),
			right,
		),
	}
}

func check(left Expression) {
	switch left.(type) {
	case *VariableExpression:
		// all good
	case *TupleExpression:
		tuple := left.(*TupleExpression)
		if len(tuple.expressions) == 0 {
			panic("GroovyBugError: one element required for left side")
		}
	default:
		panic(fmt.Sprintf("GroovyBugError: illegal left expression for declaration: %v", left))
	}
}

func (d *DeclarationExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitDeclarationExpression(d)
}

func (d *DeclarationExpression) GetVariableExpression() *VariableExpression {
	if v, ok := d.leftExpression.(*VariableExpression); ok {
		return v
	}
	return nil
}

func (d *DeclarationExpression) GetTupleExpression() *TupleExpression {
	if t, ok := d.leftExpression.(*TupleExpression); ok {
		return t
	}
	return nil
}

func (d *DeclarationExpression) GetText() string {
	var text strings.Builder

	if !d.IsMultipleAssignmentDeclaration() {
		v := d.GetVariableExpression()
		if v.IsDynamicTyped() {
			text.WriteString("def")
		} else {
			text.WriteString(FormatTypeName(v.GetType()))
		}
		text.WriteString(" ")
		text.WriteString(v.GetText())
	} else {
		t := d.GetTupleExpression()
		text.WriteString("def (")
		for i, e := range t.expressions {
			if v, ok := e.(*VariableExpression); ok {
				if !v.IsDynamicTyped() {
					text.WriteString(FormatTypeName(v.GetType()))
					text.WriteString(" ")
				}
			}
			text.WriteString(e.GetText())
			if i < len(t.expressions)-1 {
				text.WriteString(", ")
			}
		}
		text.WriteString(")")
	}
	text.WriteString(" ")
	text.WriteString(d.operation.GetText())
	text.WriteString(" ")
	text.WriteString(d.rightExpression.GetText())

	return text.String()
}

func (d *DeclarationExpression) SetLeftExpression(leftExpression Expression) {
	check(leftExpression)
	d.leftExpression = leftExpression
}

func (d *DeclarationExpression) SetRightExpression(rightExpression Expression) {
	d.rightExpression = rightExpression
}

func (d *DeclarationExpression) IsMultipleAssignmentDeclaration() bool {
	_, ok := d.leftExpression.(*TupleExpression)
	return ok
}
