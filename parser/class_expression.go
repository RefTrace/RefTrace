package parser

// ClassExpression represents access to a Java/Groovy class in an expression,
// such as when invoking a static method or accessing a static type
type ClassExpression struct {
	*BaseExpression
}

func NewClassExpression(typ IClassNode) *ClassExpression {
	ce := &ClassExpression{BaseExpression: NewBaseExpression()}
	ce.SetType(typ)
	return ce
}

func (ce *ClassExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitClassExpression(ce)
}

func (ce *ClassExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	return ce
}

func (ce *ClassExpression) GetText() string {
	return ce.GetType().GetName()
}

func (ce *ClassExpression) String() string {
	return ce.BaseExpression.GetText() + "[type: " + ce.GetType().GetName() + "]"
}
