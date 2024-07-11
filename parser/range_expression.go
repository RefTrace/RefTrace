package parser

// RangeExpression represents a range expression such as for iterating.
// E.g.: for i := range 0..10 {...}
type RangeExpression struct {
	Expression
	from           Expression
	to             Expression
	exclusiveLeft  bool
	exclusiveRight bool
}

func NewRangeExpression(from, to Expression, inclusive bool) *RangeExpression {
	return NewRangeExpressionWithExclusive(from, to, false, !inclusive)
}

func NewRangeExpressionWithExclusive(from, to Expression, exclusiveLeft, exclusiveRight bool) *RangeExpression {
	r := &RangeExpression{
		from:           from,
		to:             to,
		exclusiveLeft:  exclusiveLeft,
		exclusiveRight: exclusiveRight,
	}
	r.SetType(ClassHelper.RANGE_TYPE.GetPlainNodeReference())
	return r
}

func (r *RangeExpression) Visit(visitor GroovyCodeVisitor) {
	visitor.VisitRangeExpression(r)
}

func (r *RangeExpression) TransformExpression(transformer ExpressionTransformer) Expression {
	ret := NewRangeExpressionWithExclusive(
		transformer.Transform(r.GetFrom()),
		transformer.Transform(r.GetTo()),
		r.IsExclusiveLeft(),
		r.IsExclusiveRight(),
	)
	ret.SetSourcePosition(r)
	ret.CopyNodeMetaData(r)
	return ret
}

func (r *RangeExpression) GetFrom() Expression {
	return r.from
}

func (r *RangeExpression) GetTo() Expression {
	return r.to
}

func (r *RangeExpression) IsInclusive() bool {
	return !r.IsExclusiveRight()
}

func (r *RangeExpression) IsExclusiveLeft() bool {
	return r.exclusiveLeft
}

func (r *RangeExpression) IsExclusiveRight() bool {
	return r.exclusiveRight
}

func (r *RangeExpression) GetText() string {
	left := ""
	if r.IsExclusiveLeft() {
		left = "<"
	}
	right := ""
	if r.IsExclusiveRight() {
		right = "<"
	}
	return "(" + r.GetFrom().GetText() + left + ".." + right + r.GetTo().GetText() + ")"
}
