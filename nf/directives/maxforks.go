package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*MaxForksDirective)(nil)
var _ starlark.Value = (*MaxForksDirective)(nil)
var _ starlark.HasAttrs = (*MaxForksDirective)(nil)

func (m *MaxForksDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "num":
		return starlark.MakeInt(m.Num), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("max forks directive has no attribute %q", name))
	}
}

func (m *MaxForksDirective) AttrNames() []string {
	return []string{"num"}
}

type MaxForksDirective struct {
	Num int
}

func (m *MaxForksDirective) String() string {
	return fmt.Sprintf("MaxForksDirective(Num: %d)", m.Num)
}

func (m *MaxForksDirective) Type() string {
	return "max_forks_directive"
}

func (m *MaxForksDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *MaxForksDirective) Truth() starlark.Bool {
	return starlark.Bool(m.Num > 0)
}

func (m *MaxForksDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%d", m.Num)))
	return h.Sum32(), nil
}

func MakeMaxForksDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid max forks directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &MaxForksDirective{Num: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid max forks directive")
}
