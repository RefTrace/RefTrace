package parser

var _ Statement = (*WhileStatement)(nil)

// WhileStatement represents a while (condition) { ... } loop in Go
type WhileStatement struct {
	*BaseStatement
	BooleanExpression *BooleanExpression
	LoopBlock         Statement
}

// NewWhileStatement creates a new WhileStatement
func NewWhileStatement(booleanExpression *BooleanExpression, loopBlock Statement) *WhileStatement {
	return &WhileStatement{
		BaseStatement:     NewBaseStatement(),
		BooleanExpression: booleanExpression,
		LoopBlock:         loopBlock,
	}
}

// Visit implements the Statement interface
func (w *WhileStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitWhileLoop(w)
}

// GetBooleanExpression returns the boolean expression of the while statement
func (w *WhileStatement) GetBooleanExpression() *BooleanExpression {
	return w.BooleanExpression
}

// GetLoopBlock returns the loop block of the while statement
func (w *WhileStatement) GetLoopBlock() Statement {
	return w.LoopBlock
}

// SetBooleanExpression sets the boolean expression of the while statement
func (w *WhileStatement) SetBooleanExpression(booleanExpression *BooleanExpression) {
	w.BooleanExpression = booleanExpression
}

// SetLoopBlock sets the loop block of the while statement
func (w *WhileStatement) SetLoopBlock(loopBlock Statement) {
	w.LoopBlock = loopBlock
}
