package parser

import (
	"fmt"
)

// Expression is the base interface for any expression
type Expression interface {
	ASTNode
	GetType() *ClassNode
	SetType(*ClassNode)
	TransformExpression(transformer ExpressionTransformer) Expression
}

// BaseExpression provides a base implementation for Expression
type BaseExpression struct {
	BaseASTNode
	expressionType *ClassNode
}

var (
	EmptyExpressionArray = []Expression{}
	NullType             = MakeFromString("null")
)

func (e *BaseExpression) GetType() *ClassNode {
	if e.expressionType == NullType {
		e.expressionType = dynamicType()
	}
	return e.expressionType
}

func (e *BaseExpression) SetType(t *ClassNode) {
	e.expressionType = t
}

// TransformExpressions transforms a list of expressions
func TransformExpressions(expressions []Expression, transformer ExpressionTransformer) []Expression {
	list := make([]Expression, 0, len(expressions))
	for _, expression := range expressions {
		expression = transformer.Transform(expression)
		list = append(list, expression)
	}
	return list
}

// TransformExpressionsTyped transforms a list of expressions and checks that all transformed expressions have the given type
func TransformExpressionsTyped[T Expression](expressions []Expression, transformer ExpressionTransformer) ([]T, error) {
	list := make([]T, 0, len(expressions))
	for _, expression := range expressions {
		transformed := transformer.Transform(expression)
		typedExpr, ok := transformed.(T)
		if !ok {
			return nil, fmt.Errorf("transformed expression should have type %T but has type %T", *new(T), transformed)
		}
		list = append(list, typedExpr)
	}
	return list, nil
}

// ExpressionTransformer is an interface for transforming expressions
type ExpressionTransformer interface {
	Transform(Expression) Expression
}
