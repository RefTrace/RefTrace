package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*MaxRetriesDirective)(nil)
var _ starlark.Value = (*MaxRetriesDirective)(nil)
var _ starlark.HasAttrs = (*MaxRetriesDirective)(nil)

func (m *MaxRetriesDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "num":
		return starlark.MakeInt(m.Num), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("max retries directive has no attribute %q", name))
	}
}

func (m *MaxRetriesDirective) AttrNames() []string {
	return []string{"num"}
}

type MaxRetriesDirective struct {
	Num int
}

func (m *MaxRetriesDirective) String() string {
	return fmt.Sprintf("MaxRetriesDirective(Num: %d)", m.Num)
}

func (m *MaxRetriesDirective) Type() string {
	return "max_retries_directive"
}

func (m *MaxRetriesDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *MaxRetriesDirective) Truth() starlark.Bool {
	return starlark.Bool(m.Num > 0)
}

func (m *MaxRetriesDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%d", m.Num)))
	return h.Sum32(), nil
}

func MakeMaxRetriesDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid max retries directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &MaxRetriesDirective{Num: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid max retries directive")
}
