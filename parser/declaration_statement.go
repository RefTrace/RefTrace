// declaration_statement.go
package parser

// DeclarationStatement represents a declaration statement.
type DeclarationStatement struct {
	*BaseStatement
	Declaration *DeclarationExpression
}

// NewDeclarationStatement creates a new DeclarationStatement.
func NewDeclarationStatement(declaration *DeclarationExpression) *DeclarationStatement {
	return &DeclarationStatement{
		BaseStatement: NewBaseStatement(),
		Declaration:   declaration,
	}
}

// GetText implements the Statement interface.
func (d *DeclarationStatement) GetText() string {
	return d.Declaration.GetText()
}
