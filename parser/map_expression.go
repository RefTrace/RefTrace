package parser

import (
	"fmt"
	"strings"
)

// MapExpression represents a map expression [1 : 2, "a" : "b", x : y] which creates a mutable Map
type MapExpression struct {
	Expression
	mapEntryExpressions []*MapEntryExpression
}

func NewMapExpression() *MapExpression {
	return &MapExpression{
		mapEntryExpressions: make([]*MapEntryExpression, 0),
	}
}

func NewMapExpressionWithEntries(mapEntryExpressions []*MapEntryExpression) *MapExpression {
	me := &MapExpression{
		mapEntryExpressions: mapEntryExpressions,
	}
	// TODO: get the types of the expressions to specify the
	// map type to Map<X> if possible.
	me.SetType(MAP_TYPE)
	return me
}

func (me *MapExpression) AddMapEntryExpression(expression *MapEntryExpression) {
	me.mapEntryExpressions = append(me.mapEntryExpressions, expression)
}

func (me *MapExpression) GetMapEntryExpressions() []*MapEntryExpression {
	return me.mapEntryExpressions
}

func (me *MapExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitMapExpression(me)
}

func (me *MapExpression) IsDynamic() bool {
	return false
}

func (me *MapExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewMapExpressionWithEntries(TransformExpressions(me.GetMapEntryExpressions(), transformer))
	ret.SetSourcePosition(me)
	ret.CopyNodeMetaData(me)
	return ret
}

func (me *MapExpression) String() string {
	return fmt.Sprintf("%s%v", me.Expression.String(), me.mapEntryExpressions)
}

func (me *MapExpression) GetText() string {
	var sb strings.Builder
	sb.WriteString("[")
	size := len(me.mapEntryExpressions)
	if size > 0 {
		for i, mapEntryExpression := range me.mapEntryExpressions {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(mapEntryExpression.GetKeyExpression().GetText())
			sb.WriteString(":")
			sb.WriteString(mapEntryExpression.GetValueExpression().GetText())
			if sb.Len() > 120 && i < size-1 {
				sb.WriteString(", ... ")
				break
			}
		}
	} else {
		sb.WriteString(":")
	}
	sb.WriteString("]")
	return sb.String()
}

func (me *MapExpression) AddMapEntryExpressionWithExpressions(keyExpression, valueExpression Expression) {
	me.AddMapEntryExpression(NewMapEntryExpression(keyExpression, valueExpression))
}
