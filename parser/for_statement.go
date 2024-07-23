package parser

// ForStatement represents a standard for loop in Go
type ForStatement struct {
	*BaseStatement
	Variable             *Parameter
	CollectionExpression Expression
	LoopBlock            Statement
	Scope                *VariableScope
}

var ForLoopDummy = &Parameter{paramType: OBJECT_TYPE, name: "forLoopDummyParameter"}

func NewForStatement(variable *Parameter, collectionExpression Expression, loopBlock Statement) *ForStatement {
	return &ForStatement{
		BaseStatement:        NewBaseStatement(),
		Variable:             variable,
		CollectionExpression: collectionExpression,
		LoopBlock:            loopBlock,
	}
}

func (f *ForStatement) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitForLoop(f)
}

func (f *ForStatement) GetCollectionExpression() Expression {
	return f.CollectionExpression
}

func (f *ForStatement) GetLoopBlock() Statement {
	return f.LoopBlock
}

func (f *ForStatement) GetVariable() *Parameter {
	return f.Variable
}

func (f *ForStatement) GetVariableType() *ClassNode {
	return f.Variable.GetType()
}

func (f *ForStatement) SetCollectionExpression(collectionExpression Expression) {
	f.CollectionExpression = collectionExpression
}

func (f *ForStatement) SetVariableScope(variableScope *VariableScope) {
	f.Scope = variableScope
}

func (f *ForStatement) GetVariableScope() *VariableScope {
	return f.Scope
}

func (f *ForStatement) SetLoopBlock(loopBlock Statement) {
	f.LoopBlock = loopBlock
}
