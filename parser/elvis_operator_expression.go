package parser

// ElvisOperatorExpression represents a short ternary expression x ?: y
type ElvisOperatorExpression struct {
	*TernaryExpression
}

// NewElvisOperatorExpression creates a new ElvisOperatorExpression
func NewElvisOperatorExpression(base, falseValue Expression) *ElvisOperatorExpression {
	return &ElvisOperatorExpression{
		TernaryExpression: &TernaryExpression{
			booleanExpression: asBooleanExpression(base),
			truthExpression:   base,
			falseExpression:   falseValue,
		},
	}
}

func asBooleanExpression(base Expression) *BooleanExpression {
	baseInterface := base.(interface{})
	if be, ok := baseInterface.(*BooleanExpression); ok {
		return be
	}
	be := &BooleanExpression{BaseExpression: NewBaseExpression()}
	be.SetSourcePosition(base)
	return be
}

// TransformExpression transforms the expression
func (e *ElvisOperatorExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewElvisOperatorExpression(
		transformer.Transform(e.GetTrueExpression()),
		transformer.Transform(e.GetFalseExpression()),
	)
	ret.SetSourcePosition(e)
	ret.CopyNodeMetaData(e)
	return ret
}

// Visit implements the Visitable interface
func (e *ElvisOperatorExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitShortTernaryExpression(e)
}
