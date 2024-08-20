package parser

// NamedArgumentListExpression represents one or more arguments being passed into a method by name
type NamedArgumentListExpression struct {
	*MapExpression
}

func NewNamedArgumentListExpression() *NamedArgumentListExpression {
	return &NamedArgumentListExpression{}
}

func NewNamedArgumentListExpressionWithEntries(mapEntryExpressions []*MapEntryExpression) *NamedArgumentListExpression {
	return &NamedArgumentListExpression{
		MapExpression: NewMapExpressionWithEntries(mapEntryExpressions),
	}
}
