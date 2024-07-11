package parser

// MapEntryExpression represents an entry inside a map expression such as 1 : 2 or 'foo' : 'bar'.
type MapEntryExpression struct {
	Expression
	keyExpression   Expression
	valueExpression Expression
}

func NewMapEntryExpression(keyExpression, valueExpression Expression) *MapEntryExpression {
	return &MapEntryExpression{
		keyExpression:   keyExpression,
		valueExpression: valueExpression,
	}
}

func (m *MapEntryExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitMapEntryExpression(m)
}

func (m *MapEntryExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewMapEntryExpression(
		transformer.Transform(m.keyExpression),
		transformer.Transform(m.valueExpression),
	)
	ret.SetSourcePosition(m)
	ret.CopyNodeMetaData(m)
	return ret
}

func (m *MapEntryExpression) String() string {
	return m.Expression.String() + "(key: " + m.keyExpression.String() + ", value: " + m.valueExpression.String() + ")"
}

func (m *MapEntryExpression) GetKeyExpression() Expression {
	return m.keyExpression
}

func (m *MapEntryExpression) GetValueExpression() Expression {
	return m.valueExpression
}

func (m *MapEntryExpression) SetKeyExpression(keyExpression Expression) {
	m.keyExpression = keyExpression
}

func (m *MapEntryExpression) SetValueExpression(valueExpression Expression) {
	m.valueExpression = valueExpression
}
