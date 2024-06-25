package parser

import "fmt"

type StatementNode struct {
	BaseASTNode
	statementLabels []string
}

func NewStatementNode() *StatementNode {
	return &StatementNode{}
}

func (s *StatementNode) GetStatementLabels() []string {
	if s.statementLabels == nil {
		return nil
	}
	// Return a copy to prevent external modifications
	labels := make([]string, len(s.statementLabels))
	copy(labels, s.statementLabels)
	return labels
}

func (s *StatementNode) GetStatementLabel() string {
	if len(s.statementLabels) == 0 {
		return ""
	}
	return s.statementLabels[0]
}

func (s *StatementNode) SetStatementLabel(label string) {
	if label != "" {
		s.AddStatementLabel(label)
	}
}

func (s *StatementNode) AddStatementLabel(label string) {
	if s.statementLabels == nil {
		s.statementLabels = make([]string, 0)
	}
	s.statementLabels = append(s.statementLabels, label)
}

func (s *StatementNode) CopyStatementLabels(other *StatementNode) {
	otherLabels := other.GetStatementLabels()
	if otherLabels != nil {
		for _, label := range otherLabels {
			s.AddStatementLabel(label)
		}
	}
}

func (s *StatementNode) IsEmpty() bool {
	return false
}

func (s *StatementNode) String() string {
	return fmt.Sprintf("StatementNode[labels:%v]", s.statementLabels)
}

type DeclarationListStatement struct {
	StatementNode
	declarationStatements []*ExpressionStatement
}

func NewDeclarationListStatement(declarations ...*DeclarationExpression) *DeclarationListStatement {
	return NewDeclarationListStatementFromSlice(declarations)
}

func NewDeclarationListStatementFromSlice(declarations []*DeclarationExpression) *DeclarationListStatement {
	dls := &DeclarationListStatement{}
	dls.declarationStatements = make([]*ExpressionStatement, len(declarations))
	for i, decl := range declarations {
		dls.declarationStatements[i] = ConfigureAST(NewExpressionStatement(decl), decl)
	}
	return dls
}

func (d *DeclarationListStatement) GetDeclarationStatements() []*ExpressionStatement {
	declarationListStatementLabels := d.GetStatementLabels()

	for _, e := range d.declarationStatements {
		if declarationListStatementLabels != nil {
			// Clear existing statement labels before setting labels
			e.statementLabels = nil
			for _, label := range declarationListStatementLabels {
				e.AddStatementLabel(label)
			}
		}
	}

	return d.declarationStatements
}

func (d *DeclarationListStatement) GetDeclarationExpressions() []*DeclarationExpression {
	result := make([]*DeclarationExpression, len(d.declarationStatements))
	for i, e := range d.declarationStatements {
		result[i] = e.GetExpression().(*DeclarationExpression)
	}
	return result
}

// Helper functions (you might want to move these to a separate file)
func NewExpressionStatement(expr Expression) *ExpressionStatement {
	return &ExpressionStatement{expression: expr}
}

func ConfigureAST(stmt *ExpressionStatement, expr ASTNode) *ExpressionStatement {
	stmt.SetSourcePosition(expr)
	return stmt
}

// Additional types (you might want to define these in separate files)
type Expression interface {
	ASTNode
	// Add any additional methods specific to expressions
}

type DeclarationExpression struct {
	Expression
	// Add fields and methods specific to DeclarationExpression
}

type ExpressionStatement struct {
	StatementNode
	expression Expression
}

func (e *ExpressionStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitExpressionStatement(e)
}

func (e *ExpressionStatement) GetExpression() Expression {
	return e.expression
}

func (e *ExpressionStatement) SetExpression(expression Expression) {
	e.expression = expression
}

func (e *ExpressionStatement) GetText() string {
	return e.expression.GetText()
}

func (e *ExpressionStatement) String() string {
	return fmt.Sprintf("%s[expression:%v]", e.String(), e.expression)
}
