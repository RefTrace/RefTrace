package parser

import (
	"fmt"
	"strings"
)

// GStringExpression represents a String expression which contains embedded values inside
// it such as "hello there ${user} how are you" which is expanded lazily
type GStringExpression struct {
	*BaseExpression
	verbatimText string
	strings      []*ConstantExpression
	values       []Expression
}

func NewGStringExpression(verbatimText string) *GStringExpression {
	expr := &GStringExpression{
		BaseExpression: NewBaseExpression(),
		verbatimText:   verbatimText,
		strings:        make([]*ConstantExpression, 0),
		values:         make([]Expression, 0),
	}
	expr.SetType(GSTRING_TYPE)
	return expr
}

func NewGStringExpressionWithValues(verbatimText string, strings []*ConstantExpression, values []Expression) *GStringExpression {
	expr := &GStringExpression{
		BaseExpression: NewBaseExpression(),
		verbatimText:   verbatimText,
		strings:        strings,
		values:         values,
	}
	expr.SetType(GSTRING_TYPE)
	return expr
}

func (g *GStringExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitGStringExpression(g)
}

func (g *GStringExpression) String() string {
	return fmt.Sprintf("%s[strings: %v values: %v]", g.BaseExpression.GetText(), g.strings, g.values)
}

func (g *GStringExpression) GetText() string {
	return g.verbatimText
}

func (g *GStringExpression) GetStrings() []*ConstantExpression {
	return g.strings
}

func (g *GStringExpression) GetValues() []Expression {
	return g.values
}

func (g *GStringExpression) AddString(text *ConstantExpression) {
	if text == nil {
		panic("Cannot add a null text expression")
	}
	g.strings = append(g.strings, text)
}

func (g *GStringExpression) AddValue(value Expression) {
	// If the first thing is a value, then we need a dummy empty string in front of it so that when we
	// toString it they come out in the correct order.
	if len(g.strings) == 0 {
		g.strings = append(g.strings, EMPTY_STRING)
	}
	g.values = append(g.values, value)
}

func (g *GStringExpression) GetValue(idx int) Expression {
	return g.values[idx]
}

func (g *GStringExpression) IsConstantString() bool {
	return len(g.values) == 0
}

func (g *GStringExpression) AsConstantString() Expression {
	var buffer strings.Builder
	for _, expression := range g.strings {
		value := expression.GetValue()
		if value != nil {
			buffer.WriteString(fmt.Sprintf("%v", value))
		}
	}
	return NewConstantExpression(buffer.String())
}
