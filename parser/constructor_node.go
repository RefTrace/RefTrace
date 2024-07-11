package parser

// ConstructorNode represents a constructor declaration
type ConstructorNode struct {
	MethodNode
}

// NewConstructorNode creates a new ConstructorNode with the given modifiers and code
func NewConstructorNode(modifiers int, code Statement) *ConstructorNode {
	return NewConstructorNodeWithParams(modifiers, []*Parameter{}, []*ClassNode{}, code)
}

// NewConstructorNodeWithParams creates a new ConstructorNode with the given modifiers, parameters, exceptions, and code
func NewConstructorNodeWithParams(modifiers int, parameters []*Parameter, exceptions []*ClassNode, code Statement) *ConstructorNode {
	return &ConstructorNode{
		MethodNode: *NewMethodNode("<init>", modifiers, VOID_TYPE, parameters, exceptions, code),
	}
}

// FirstStatementIsSpecialConstructorCall checks if the first statement is a special constructor call
func (c *ConstructorNode) FirstStatementIsSpecialConstructorCall() bool {
	code := c.GetFirstStatement()
	exprStmt, ok := code.(*ExpressionStatement)
	if !ok {
		return false
	}

	expression := exprStmt.GetExpression()
	cce, ok := expression.(*ConstructorCallExpression)
	if !ok {
		return false
	}

	return cce.IsSpecialCall()
}
