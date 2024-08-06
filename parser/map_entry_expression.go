package parser

// MapEntryExpression represents an entry inside a map expression such as 1 : 2 or 'foo' : 'bar'.
type MapEntryExpression struct {
	*BaseExpression
	keyExpression   Expression
	valueExpression Expression
}

func NewMapEntryExpression(keyExpression, valueExpression Expression) *MapEntryExpression {
	return &MapEntryExpression{
		BaseExpression:  NewBaseExpression(),
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
	return m.BaseExpression.GetText() + "(key: " + m.keyExpression.GetText() + ", value: " + m.valueExpression.GetText() + ")"
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
