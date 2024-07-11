package parser

// ElvisOperatorExpression represents a short ternary expression x ?: y
type ElvisOperatorExpression struct {
	TernaryExpression
}

// NewElvisOperatorExpression creates a new ElvisOperatorExpression
func NewElvisOperatorExpression(base, falseValue Expression) *ElvisOperatorExpression {
	return &ElvisOperatorExpression{
		TernaryExpression: TernaryExpression{
			BooleanExpression: asBooleanExpression(base),
			TrueExpression:    base,
			FalseExpression:   falseValue,
		},
	}
}

func asBooleanExpression(base Expression) BooleanExpression {
	if be, ok := base.(BooleanExpression); ok {
		return be
	}
	be := BooleanExpression{Expression: base}
	be.SetSourcePosition(base.GetSourcePosition())
	return be
}

// TransformExpression transforms the expression
func (e *ElvisOperatorExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewElvisOperatorExpression(
		transformer.Transform(e.GetTrueExpression()),
		transformer.Transform(e.GetFalseExpression()),
	)
	ret.SetSourcePosition(e.GetSourcePosition())
	ret.CopyNodeMetaData(e)
	return ret
}

// Visit implements the Visitable interface
func (e *ElvisOperatorExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitShortTernaryExpression(e)
}
