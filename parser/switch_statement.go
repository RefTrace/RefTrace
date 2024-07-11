package parser

import "ast"

// SwitchStatement represents a switch (object) { case value: ... case [1, 2, 3]: ...  default: ... } statement in Go.
type SwitchStatement struct {
	ast.Statement
	expression       ast.Expression
	caseStatements   []*CaseStatement
	defaultStatement ast.Statement
}

func NewSwitchStatement(expression ast.Expression) *SwitchStatement {
	return NewSwitchStatementWithDefault(expression, ast.EmptyStatement)
}

func NewSwitchStatementWithDefault(expression ast.Expression, defaultStatement ast.Statement) *SwitchStatement {
	return &SwitchStatement{
		expression:       expression,
		defaultStatement: defaultStatement,
	}
}

func NewSwitchStatementFull(expression ast.Expression, caseStatements []*CaseStatement, defaultStatement ast.Statement) *SwitchStatement {
	return &SwitchStatement{
		expression:       expression,
		caseStatements:   caseStatements,
		defaultStatement: defaultStatement,
	}
}

func (s *SwitchStatement) Visit(visitor ast.GroovyCodeVisitor) {
	visitor.VisitSwitch(s)
}

func (s *SwitchStatement) GetCaseStatements() []*CaseStatement {
	return s.caseStatements
}

func (s *SwitchStatement) GetExpression() ast.Expression {
	return s.expression
}

func (s *SwitchStatement) SetExpression(e ast.Expression) {
	s.expression = e
}

func (s *SwitchStatement) GetDefaultStatement() ast.Statement {
	return s.defaultStatement
}

func (s *SwitchStatement) SetDefaultStatement(defaultStatement ast.Statement) {
	s.defaultStatement = defaultStatement
}

func (s *SwitchStatement) AddCase(caseStatement *CaseStatement) {
	s.caseStatements = append(s.caseStatements, caseStatement)
}

// GetCaseStatement returns the case statement of the given index or nil
func (s *SwitchStatement) GetCaseStatement(idx int) *CaseStatement {
	if idx >= 0 && idx < len(s.caseStatements) {
		return s.caseStatements[idx]
	}
	return nil
}
