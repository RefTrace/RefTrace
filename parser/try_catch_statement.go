package parser

import (
	"errors"
)

// TryCatchStatement represents a try { ... } catch () finally {} statement in Go
type TryCatchStatement struct {
	TryStatement       Statement
	FinallyStatement   Statement
	CatchStatements    []CatchStatement
	ResourceStatements []ExpressionStatement
}

func NewTryCatchStatement(tryStatement, finallyStatement Statement) *TryCatchStatement {
	return &TryCatchStatement{
		TryStatement:     tryStatement,
		FinallyStatement: finallyStatement,
	}
}

func (t *TryCatchStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitTryCatchFinally(t)
}

func (t *TryCatchStatement) GetTryStatement() Statement {
	return t.TryStatement
}

func (t *TryCatchStatement) GetFinallyStatement() Statement {
	return t.FinallyStatement
}

func (t *TryCatchStatement) GetCatchStatement(idx int) *CatchStatement {
	if idx >= 0 && idx < len(t.CatchStatements) {
		return &t.CatchStatements[idx]
	}
	return nil
}

func (t *TryCatchStatement) GetCatchStatements() []CatchStatement {
	return t.CatchStatements
}

func (t *TryCatchStatement) GetResourceStatement(idx int) *ExpressionStatement {
	if idx >= 0 && idx < len(t.ResourceStatements) {
		return &t.ResourceStatements[idx]
	}
	return nil
}

func (t *TryCatchStatement) GetResourceStatements() []ExpressionStatement {
	return t.ResourceStatements
}

func IsResource(expression Expression) bool {
	isResource, ok := expression.GetNodeMetaData("_IS_RESOURCE").(bool)
	return ok && isResource
}

func (t *TryCatchStatement) SetTryStatement(tryStatement Statement) {
	t.TryStatement = tryStatement
}

func (t *TryCatchStatement) SetFinallyStatement(finallyStatement Statement) {
	t.FinallyStatement = finallyStatement
}

func (t *TryCatchStatement) SetCatchStatement(idx int, catchStatement CatchStatement) {
	t.CatchStatements[idx] = catchStatement
}

func (t *TryCatchStatement) AddCatch(catchStatement CatchStatement) *TryCatchStatement {
	t.CatchStatements = append(t.CatchStatements, catchStatement)
	return t
}

func (t *TryCatchStatement) AddResource(resourceStatement ExpressionStatement) (*TryCatchStatement, error) {
	resourceExpression := resourceStatement.GetExpression()
	if _, ok := resourceExpression.(*DeclarationExpression); !ok {
		if _, ok := resourceExpression.(*VariableExpression); !ok {
			return nil, errors.New("resourceStatement should be a variable declaration statement or a variable")
		}
	}
	resourceExpression.SetNodeMetaData("_IS_RESOURCE", true)
	t.ResourceStatements = append(t.ResourceStatements, resourceStatement)
	return t, nil
}
