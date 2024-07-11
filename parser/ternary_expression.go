package parser

import (
	"fmt"
)

// TernaryExpression represents a ternary expression (booleanExpression) ? expression : expression
type TernaryExpression struct {
	Expression
	booleanExpression BooleanExpression
	truthExpression   Expression
	falseExpression   Expression
}

func NewTernaryExpression(booleanExpression BooleanExpression, truthExpression, falseExpression Expression) *TernaryExpression {
	return &TernaryExpression{
		booleanExpression: booleanExpression,
		truthExpression:   truthExpression,
		falseExpression:   falseExpression,
	}
}

func (t *TernaryExpression) GetBooleanExpression() BooleanExpression {
	return t.booleanExpression
}

func (t *TernaryExpression) GetTrueExpression() Expression {
	return t.truthExpression
}

func (t *TernaryExpression) GetFalseExpression() Expression {
	return t.falseExpression
}

func (t *TernaryExpression) GetText() string {
	return fmt.Sprintf("(%s) ? %s : %s", t.booleanExpression.GetText(), t.truthExpression.GetText(), t.falseExpression.GetText())
}

func (t *TernaryExpression) GetType() *ClassNode {
	w := WideningCategories{}
	return w.lowestUpperBound(t.truthExpression.GetType(), t.falseExpression.GetType())
}

func (t *TernaryExpression) String() string {
	return fmt.Sprintf("%s[%s ? %s : %s]", t.Expression.GetText(), t.booleanExpression, t.truthExpression, t.falseExpression)
}

func (t *TernaryExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitTernaryExpression(t)
}
