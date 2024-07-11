package parser

// LoopingStatement is an interface that provides some sort of looping mechanism.
// Typically in the form of a block that will be executed repeatedly.
// DoWhileStatements, WhileStatements, and ForStatements are all examples of LoopingStatements.
type LoopingStatement interface {
	// GetLoopBlock gets the loop block.
	GetLoopBlock() Statement

	// SetLoopBlock sets the loop block.
	SetLoopBlock(loopBlock Statement)
}
