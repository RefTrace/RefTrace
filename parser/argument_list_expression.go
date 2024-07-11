package parser

// ArgumentListExpression represents one or more arguments being passed into a method
type ArgumentListExpression struct {
	TupleExpression
}

var (
	EmptyArray     = []interface{}{}
	EmptyArguments = NewArgumentListExpression()
)

func NewArgumentListExpression() *ArgumentListExpression {
	return &ArgumentListExpression{}
}

func NewArgumentListExpressionFromList(expressions []Expression) *ArgumentListExpression {
	return &ArgumentListExpression{TupleExpression{Expressions: expressions}}
}

func NewArgumentListExpressionFromSlice(expressions ...Expression) *ArgumentListExpression {
	return &ArgumentListExpression{TupleExpression{Expressions: expressions}}
}

func NewArgumentListExpressionFromParameters(parameters []*Parameter) *ArgumentListExpression {
	ale := NewArgumentListExpression()
	for _, param := range parameters {
		ale.AddExpression(NewVariableExpression(param))
	}
	return ale
}

func (a *ArgumentListExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewArgumentListExpressionFromList(a.TransformExpressions(a.Expressions, transformer))
	ret.SetSourcePosition(a)
	ret.CopyNodeMetaData(a)
	return ret
}

func (a *ArgumentListExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitArgumentlistExpression(a)
}
