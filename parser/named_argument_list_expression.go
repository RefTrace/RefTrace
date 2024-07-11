package parser

// NamedArgumentListExpression represents one or more arguments being passed into a method by name
type NamedArgumentListExpression struct {
	MapExpression
}

func NewNamedArgumentListExpression() *NamedArgumentListExpression {
	return &NamedArgumentListExpression{}
}

func NewNamedArgumentListExpressionWithEntries(mapEntryExpressions []MapEntryExpression) *NamedArgumentListExpression {
	return &NamedArgumentListExpression{
		MapExpression: MapExpression{MapEntryExpressions: mapEntryExpressions},
	}
}

func (n *NamedArgumentListExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewNamedArgumentListExpressionWithEntries(
		TransformExpressions(n.GetMapEntryExpressions(), transformer, func() Expression { return &MapEntryExpression{} }),
	)
	ret.SetSourcePosition(n)
	ret.CopyNodeMetaData(n)
	return ret
}
