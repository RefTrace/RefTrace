package parser

// SwitchStatement represents a switch (object) { case value: ... case [1, 2, 3]: ...  default: ... } statement in Go.
type SwitchStatement struct {
	*BaseStatement
	expression       Expression
	caseStatements   []*CaseStatement
	defaultStatement Statement
}

func NewSwitchStatement(expression Expression) *SwitchStatement {
	return NewSwitchStatementWithDefault(expression, NewEmptyStatement())
}

func NewSwitchStatementWithDefault(expression Expression, defaultStatement Statement) *SwitchStatement {
	return &SwitchStatement{
		BaseStatement:    NewBaseStatement(),
		expression:       expression,
		defaultStatement: defaultStatement,
	}
}

func NewSwitchStatementFull(expression Expression, caseStatements []*CaseStatement, defaultStatement Statement) *SwitchStatement {
	return &SwitchStatement{
		BaseStatement:    NewBaseStatement(),
		expression:       expression,
		caseStatements:   caseStatements,
		defaultStatement: defaultStatement,
	}
}

func (s *SwitchStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitSwitch(s)
}

func (s *SwitchStatement) GetCaseStatements() []*CaseStatement {
	return s.caseStatements
}

func (s *SwitchStatement) GetExpression() Expression {
	return s.expression
}

func (s *SwitchStatement) SetExpression(e Expression) {
	s.expression = e
}

func (s *SwitchStatement) GetDefaultStatement() Statement {
	return s.defaultStatement
}

func (s *SwitchStatement) SetDefaultStatement(defaultStatement Statement) {
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
