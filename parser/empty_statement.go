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

func (e *EmptyStatement) throwUnsupportedOperationException() error {
	return errors.New("EmptyStatement.INSTANCE is immutable")
}

// ASTNode overrides

func (e *EmptyStatement) SetColumnNumber(n int) {
	panic(e.throwUnsupportedOperationException())
}

func (e *EmptyStatement) SetLastColumnNumber(n int) {
	panic(e.throwUnsupportedOperationException())
}

func (e *EmptyStatement) SetLastLineNumber(n int) {
	panic(e.throwUnsupportedOperationException())
}

func (e *EmptyStatement) SetLineNumber(n int) {
	panic(e.throwUnsupportedOperationException())
}

func (e *EmptyStatement) SetMetaDataMap(meta map[interface{}]interface{}) {
	panic(e.throwUnsupportedOperationException())
}

func (e *EmptyStatement) SetSourcePosition(node ASTNode) {
	panic(e.throwUnsupportedOperationException())
}

// Statement overrides

// TODO: check errors
func (e *EmptyStatement) AddStatementLabel(label string) {
	return
}

func (e *EmptyStatement) SetStatementLabel(label string) {
	return
}
