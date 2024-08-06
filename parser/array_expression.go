package parser

import (
	"fmt"
	"strings"
)

// ArrayExpression represents an array object construction.
type ArrayExpression struct {
	*BaseExpression
	initExpressions []Expression
	sizeExpressions []Expression
	elementType     *ClassNode
}

func makeArray(base *ClassNode, sizeExpressions []Expression) *ClassNode {
	ret := base.MakeArray()
	if sizeExpressions == nil {
		return ret
	}
	size := len(sizeExpressions)
	for i := 1; i < size; i++ {
		ret = ret.MakeArray()
	}
	return ret
}

func NewArrayExpression(elementType *ClassNode, initExpressions, sizeExpressions []Expression) *ArrayExpression {
	ae := &ArrayExpression{
		BaseExpression:  NewBaseExpression(),
		elementType:     elementType,
		initExpressions: initExpressions,
		sizeExpressions: sizeExpressions,
	}
	ae.SetType(makeArray(elementType, sizeExpressions))

	// Handle initExpressions
	if initExpressions != nil {
		ae.initExpressions = initExpressions
	}

	// Check for valid input
	if initExpressions == nil && (sizeExpressions == nil || len(sizeExpressions) == 0) {
		panic("Either an initializer or defined size must be given")
	}

	if len(ae.initExpressions) > 0 && sizeExpressions != nil && len(sizeExpressions) > 0 {
		panic(fmt.Sprintf("Both an initializer (%s) and a defined size (%s) cannot be given",
			ae.formatInitExpressions(), ae.formatSizeExpressions()))
	}

	// Type check initExpressions
	for _, item := range ae.initExpressions {
		if item != nil && !isExpression(item) {
			panic(fmt.Sprintf("Item: %v is not an Expression", item))
		}
	}

	// Type check sizeExpressions if no initializer
	if len(ae.initExpressions) == 0 {
		for _, item := range sizeExpressions {
			if !isExpression(item) {
				panic(fmt.Sprintf("Item: %v is not an Expression", item))
			}
		}
	}

	return ae
}

// Helper function to check if an interface{} is an Expression
func isExpression(item interface{}) bool {
	_, ok := item.(Expression)
	return ok
}

func (ae *ArrayExpression) AddExpression(initExpression Expression) {
	ae.initExpressions = append(ae.initExpressions, initExpression)
}

func (ae *ArrayExpression) GetExpressions() []Expression {
	return ae.initExpressions
}

func (ae *ArrayExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitArrayExpression(ae)
}

func (ae *ArrayExpression) IsDynamic() bool {
	return false
}

func (ae *ArrayExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	exprList := TransformExpressions(ae.initExpressions, transformer)
	var sizes []Expression
	if !ae.HasInitializer() {
		sizes = TransformExpressions(ae.sizeExpressions, transformer)
	}
	ret := NewArrayExpression(ae.elementType, exprList, sizes)
	ret.SetSourcePosition(ae)
	ret.CopyNodeMetaData(ae)
	return ret
}

func (ae *ArrayExpression) GetExpression(i int) Expression {
	return ae.initExpressions[i]
}

func (ae *ArrayExpression) GetElementType() *ClassNode {
	return ae.elementType
}

func (ae *ArrayExpression) GetText() string {
	return "[" + ae.formatInitExpressions() + "]"
}

func (ae *ArrayExpression) formatInitExpressions() string {
	texts := make([]string, len(ae.initExpressions))
	for i, expr := range ae.initExpressions {
		texts[i] = expr.GetText()
	}
	return "{" + strings.Join(texts, ", ") + "}"
}

func (ae *ArrayExpression) formatSizeExpressions() string {
	texts := make([]string, len(ae.sizeExpressions))
	for i, expr := range ae.sizeExpressions {
		texts[i] = "[" + expr.GetText() + "]"
	}
	return strings.Join(texts, "")
}

func (ae *ArrayExpression) HasInitializer() bool {
	return ae.sizeExpressions == nil
}

func (ae *ArrayExpression) GetSizeExpression() []Expression {
	return ae.sizeExpressions
}

func (ae *ArrayExpression) String() string {
	if ae.HasInitializer() {
		return fmt.Sprintf("%s[elementType: %v, init: {%s}]", ae.BaseExpression.GetText(), ae.GetElementType(), ae.formatInitExpressions())
	}
	return fmt.Sprintf("%s[elementType: %v, size: %s]", ae.BaseExpression.GetText(), ae.GetElementType(), ae.formatSizeExpressions())
}
