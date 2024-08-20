package parser

// DoWhileStatement represents a do { ... } while (condition) loop in Go
type DoWhileStatement struct {
	*BaseStatement
	booleanExpression *BooleanExpression
	loopBlock         Statement
}

// NewDoWhileStatement creates a new DoWhileStatement
func NewDoWhileStatement(booleanExpression *BooleanExpression, loopBlock Statement) *DoWhileStatement {
	return &DoWhileStatement{
		BaseStatement:     NewBaseStatement(),
		booleanExpression: booleanExpression,
		loopBlock:         loopBlock,
	}
}

// Visit implements the Statement interface
func (d *DoWhileStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitDoWhileLoop(d)
}

// GetBooleanExpression returns the boolean expression of the do-while loop
func (d *DoWhileStatement) GetBooleanExpression() *BooleanExpression {
	return d.booleanExpression
}

// GetLoopBlock returns the loop block of the do-while loop
func (d *DoWhileStatement) GetLoopBlock() Statement {
	return d.loopBlock
}

// SetBooleanExpression sets the boolean expression of the do-while loop
func (d *DoWhileStatement) SetBooleanExpression(booleanExpression *BooleanExpression) {
	d.booleanExpression = booleanExpression
}

// SetLoopBlock sets the loop block of the do-while loop
func (d *DoWhileStatement) SetLoopBlock(loopBlock Statement) {
	d.loopBlock = loopBlock
}
