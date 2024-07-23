package parser

import (
	"fmt"
)

// BinaryExpression represents two expressions and an operation
type BinaryExpression struct {
	*BaseExpression
	leftExpression  Expression
	rightExpression Expression
	operation       *Token
	safe            bool
}

// NewBinaryExpression creates a new BinaryExpression
func NewBinaryExpression(leftExpression Expression, operation *Token, rightExpression Expression) *BinaryExpression {
	return &BinaryExpression{
		BaseExpression:  NewBaseExpression(),
		leftExpression:  leftExpression,
		rightExpression: rightExpression,
		operation:       operation,
		safe:            false,
	}
}

// NewBinaryExpressionWithSafe creates a new BinaryExpression with safe flag
func NewBinaryExpressionWithSafe(leftExpression Expression, operation *Token, rightExpression Expression, safe bool) *BinaryExpression {
	be := NewBinaryExpression(leftExpression, operation, rightExpression)
	be.safe = safe
	return be
}

func (be *BinaryExpression) String() string {
	return fmt.Sprintf("%s[%v%v%v]", be.BaseExpression.GetText(), be.leftExpression, be.operation, be.rightExpression)
}

func (be *BinaryExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitBinaryExpression(be)
}

func (be *BinaryExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewBinaryExpressionWithSafe(
		transformer.Transform(be.leftExpression),
		be.operation,
		transformer.Transform(be.rightExpression),
		be.safe,
	)
	ret.SetSourcePosition(be)
	ret.CopyNodeMetaData(be)
	return ret
}

func (be *BinaryExpression) GetLeftExpression() Expression {
	return be.leftExpression
}

func (be *BinaryExpression) SetLeftExpression(leftExpression Expression) {
	be.leftExpression = leftExpression
}

func (be *BinaryExpression) SetRightExpression(rightExpression Expression) {
	be.rightExpression = rightExpression
}

func (be *BinaryExpression) GetOperation() *Token {
	return be.operation
}

func (be *BinaryExpression) GetRightExpression() Expression {
	return be.rightExpression
}

func (be *BinaryExpression) GetText() string {
	if be.operation.GetType() == LEFT_SQUARE_BRACKET {
		safeOp := ""
		if be.safe {
			safeOp = "?"
		}
		return fmt.Sprintf("%s%s[%s]", be.leftExpression.GetText(), safeOp, be.rightExpression.GetText())
	}
	return fmt.Sprintf("(%s %s %s)", be.leftExpression.GetText(), be.operation.GetText(), be.rightExpression.GetText())
}

func (be *BinaryExpression) IsSafe() bool {
	return be.safe
}

func (be *BinaryExpression) SetSafe(safe bool) {
	be.safe = safe
}

// NewAssignmentExpression creates an assignment expression
func NewAssignmentExpression(variable Variable, rhs Expression) *BinaryExpression {
	lhs := NewVariableExpressionWithVariable(variable)
	operator := NewToken(ASSIGN, "=", rhs.GetLineNumber(), rhs.GetColumnNumber())

	return NewBinaryExpression(lhs, operator, rhs)
}

// NewInitializationExpression creates a variable initialization expression
func NewInitializationExpression(variable string, typ *ClassNode, rhs Expression) *BinaryExpression {
	lhs := NewVariableExpressionWithString(variable)

	if typ != nil {
		lhs.SetType(typ)
	}

	operator := NewToken(ASSIGN, "=", typ.GetLineNumber(), typ.GetColumnNumber())

	return NewBinaryExpression(lhs, operator, rhs)
}
