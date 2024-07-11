package parser

import (
	"strings"
)

// IfStatement represents an if (condition) { ... } else { ... } statement in Go.
type IfStatement struct {
	BooleanExpression BooleanExpression
	IfBlock           Statement
	ElseBlock         Statement
}

// NewIfStatement creates a new IfStatement with the given condition and blocks.
func NewIfStatement(booleanExpression BooleanExpression, ifBlock Statement, elseBlock Statement) *IfStatement {
	return &IfStatement{
		BooleanExpression: booleanExpression,
		IfBlock:           ifBlock,
		ElseBlock:         elseBlock,
	}
}

// SetBooleanExpression sets the boolean expression for the if statement.
func (i *IfStatement) SetBooleanExpression(booleanExpression BooleanExpression) {
	i.BooleanExpression = booleanExpression
}

// SetIfBlock sets the if block for the if statement.
func (i *IfStatement) SetIfBlock(statement Statement) {
	i.IfBlock = statement
}

// SetElseBlock sets the else block for the if statement.
func (i *IfStatement) SetElseBlock(statement Statement) {
	i.ElseBlock = statement
}

// Visit implements the Statement interface.
func (i *IfStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitIfElse(i)
}

// GetBooleanExpression returns the boolean expression of the if statement.
func (i *IfStatement) GetBooleanExpression() BooleanExpression {
	return i.BooleanExpression
}

// GetIfBlock returns the if block of the if statement.
func (i *IfStatement) GetIfBlock() Statement {
	return i.IfBlock
}

// GetElseBlock returns the else block of the if statement.
func (i *IfStatement) GetElseBlock() Statement {
	return i.ElseBlock
}

// GetText returns a string representation of the if statement.
func (i *IfStatement) GetText() string {
	var text strings.Builder
	text.WriteString("if (")
	text.WriteString(i.BooleanExpression.GetText())
	text.WriteString(") ")
	text.WriteString(i.IfBlock.GetText())

	if i.ElseBlock != nil && !i.ElseBlock.IsEmpty() {
		if _, ok := i.IfBlock.(*BlockStatement); !ok {
			text.WriteString(";")
		}
		text.WriteString(" else ")
		text.WriteString(i.ElseBlock.GetText())
	}

	return text.String()
}
