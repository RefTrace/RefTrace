package parser

import (
	"fmt"
	"strings"
)

// DeclarationExpression represents one or more local variables.
type DeclarationExpression struct {
	BinaryExpression
}

// NewDeclarationExpression creates a new DeclarationExpression.
func NewDeclarationExpression(left Expression, operation *Token, right Expression) *DeclarationExpression {
	check(left)
	return &DeclarationExpression{
		BinaryExpression: BinaryExpression{
			Left:      left,
			Operation: NewToken(ASSIGN, "=", operation.StartLine, operation.StartColumn),
			Right:     right,
		},
	}
}

func check(left Expression) {
	switch left.(type) {
	case *VariableExpression:
		// all good
	case *TupleExpression:
		tuple := left.(*TupleExpression)
		if len(tuple.Expressions) == 0 {
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
	if v, ok := d.Left.(*VariableExpression); ok {
		return v
	}
	return nil
}

func (d *DeclarationExpression) GetTupleExpression() *TupleExpression {
	if t, ok := d.Left.(*TupleExpression); ok {
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
			text.WriteString(FormatTypeName(v.Type))
		}
		text.WriteString(" ")
		text.WriteString(v.GetText())
	} else {
		t := d.GetTupleExpression()
		text.WriteString("def (")
		for i, e := range t.Expressions {
			if v, ok := e.(*VariableExpression); ok {
				if !v.IsDynamicTyped() {
					text.WriteString(FormatTypeName(v.Type))
					text.WriteString(" ")
				}
			}
			text.WriteString(e.GetText())
			if i < len(t.Expressions)-1 {
				text.WriteString(", ")
			}
		}
		text.WriteString(")")
	}
	text.WriteString(" ")
	text.WriteString(d.Operation.Text)
	text.WriteString(" ")
	text.WriteString(d.Right.GetText())

	return text.String()
}

func (d *DeclarationExpression) SetLeftExpression(leftExpression Expression) {
	check(leftExpression)
	d.Left = leftExpression
}

func (d *DeclarationExpression) SetRightExpression(rightExpression Expression) {
	d.Right = rightExpression
}

func (d *DeclarationExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewDeclarationExpression(
		transformer.Transform(d.Left),
		d.Operation,
		transformer.Transform(d.Right),
	)
	ret.SetSourcePosition(d)
	ret.AddAnnotations(d.GetAnnotations())
	ret.SetDeclaringClass(d.GetDeclaringClass())
	ret.CopyNodeMetaData(d)
	return ret
}

func (d *DeclarationExpression) IsMultipleAssignmentDeclaration() bool {
	_, ok := d.Left.(*TupleExpression)
	return ok
}
