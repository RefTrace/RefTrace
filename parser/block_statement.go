package parser

import (
	"strings"
)

// BlockStatement represents a list of statements and a scope.
type BlockStatement struct {
	Statement
	statements []Statement
	scope      *VariableScope
}

// NewBlockStatement creates a new BlockStatement with an empty list of statements and a new VariableScope.
func NewBlockStatement() *BlockStatement {
	return &BlockStatement{
		statements: make([]Statement, 0),
		scope:      NewVariableScope(),
	}
}

// NewBlockStatementWithStatementsAndScope creates a BlockStatement with a scope and children statements.
func NewBlockStatementWithStatementsAndScope(statements []Statement, scope *VariableScope) *BlockStatement {
	return &BlockStatement{
		statements: append([]Statement(nil), statements...),
		scope:      scope,
	}
}

// Visit implements the GroovyCodeVisitor interface.
func (bs *BlockStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitBlockStatement(bs)
}

// GetStatements returns the list of statements.
func (bs *BlockStatement) GetStatements() []Statement {
	return bs.statements
}

// AddStatement adds a statement to the list of statements.
func (bs *BlockStatement) AddStatement(statement Statement) {
	bs.statements = append(bs.statements, statement)
}

// AddStatements adds a list of statements to the existing list of statements.
func (bs *BlockStatement) AddStatements(listOfStatements []Statement) {
	bs.statements = append(bs.statements, listOfStatements...)
}

// GetText returns a string representation of the block statement.
func (bs *BlockStatement) GetText() string {
	var texts []string
	for _, statement := range bs.statements {
		texts = append(texts, statement.GetText())
	}
	return "{ " + strings.Join(texts, "; ") + " }"
}

// String returns a string representation of the BlockStatement.
func (bs *BlockStatement) String() string {
	return bs.Statement.String() + strings.Join(statementsToStrings(bs.statements), ", ")
}

// IsEmpty returns true if the block statement has no statements.
func (bs *BlockStatement) IsEmpty() bool {
	return len(bs.statements) == 0
}

// GetVariableScope returns the variable scope of the block statement.
func (bs *BlockStatement) GetVariableScope() *VariableScope {
	return bs.scope
}

// SetVariableScope sets the variable scope of the block statement.
func (bs *BlockStatement) SetVariableScope(scope *VariableScope) {
	bs.scope = scope
}

// Helper function to convert statements to strings
func statementsToStrings(statements []Statement) []string {
	var strs []string
	for _, stmt := range statements {
		strs = append(strs, stmt.String())
	}
	return strs
}
