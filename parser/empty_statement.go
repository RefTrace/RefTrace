package parser

import (
	"errors"
)

type EmptyStatement struct {
	Statement
}

func NewEmptyStatement() *EmptyStatement {
	return &EmptyStatement{}
}

func (e *EmptyStatement) IsEmpty() bool {
	return true
}

func (e *EmptyStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitEmptyStatement(e)
}

// INSTANCE is an immutable singleton recommended for use when source range
// or any other occurrence-specific metadata is not needed.
var INSTANCE = &EmptyStatement{
	Statement: Statement{
		ASTNode: ASTNode{},
	},
}

func (e *EmptyStatement) throwUnsupportedOperationException() error {
	return errors.New("EmptyStatement.INSTANCE is immutable")
}

// ASTNode overrides

func (e *EmptyStatement) SetColumnNumber(n int) error {
	return e.throwUnsupportedOperationException()
}

func (e *EmptyStatement) SetLastColumnNumber(n int) error {
	return e.throwUnsupportedOperationException()
}

func (e *EmptyStatement) SetLastLineNumber(n int) error {
	return e.throwUnsupportedOperationException()
}

func (e *EmptyStatement) SetLineNumber(n int) error {
	return e.throwUnsupportedOperationException()
}

func (e *EmptyStatement) SetMetaDataMap(meta map[interface{}]interface{}) error {
	return e.throwUnsupportedOperationException()
}

func (e *EmptyStatement) SetSourcePosition(node ASTNode) error {
	return e.throwUnsupportedOperationException()
}

// Statement overrides

func (e *EmptyStatement) AddStatementLabel(label string) error {
	return e.throwUnsupportedOperationException()
}

func (e *EmptyStatement) SetStatementLabel(label string) error {
	return e.throwUnsupportedOperationException()
}
