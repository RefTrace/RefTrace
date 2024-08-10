package parser

// ArgumentListExpression represents one or more arguments being passed into a method
type ArgumentListExpression struct {
	*TupleExpression
}

var (
	EmptyArray     = []interface{}{}
	EmptyArguments = NewArgumentListExpression()
)

func NewArgumentListExpression() *ArgumentListExpression {
	return &ArgumentListExpression{NewTupleExpression()}
}

func NewArgumentListExpressionFromList(expressions []Expression) *ArgumentListExpression {
	return &ArgumentListExpression{NewTupleExpressionWithExpressions(expressions...)}
}

func NewArgumentListExpressionFromSlice(expressions ...Expression) *ArgumentListExpression {
	return &ArgumentListExpression{NewTupleExpressionWithExpressions(expressions...)}
}

func NewArgumentListExpressionFromParameters(parameters []*Parameter) *ArgumentListExpression {
	ale := NewArgumentListExpression()
	for _, param := range parameters {
		ale.AddExpression(NewVariableExpressionWithVariable(param))
	}
	return ale
}

func (a *ArgumentListExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewArgumentListExpressionFromList(TransformExpressions(a.expressions, transformer))
	ret.SetSourcePosition(a)
	ret.CopyNodeMetaData(a)
	return ret
}

func (a *ArgumentListExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitArgumentlistExpression(a)
}
