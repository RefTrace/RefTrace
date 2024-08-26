package directives

import (
	"errors"
	"fmt"
	"hash/fnv"
	"reft-go/parser"

	"go.starlark.net/starlark"
)

var _ Directive = (*MaxErrorsDirective)(nil)
var _ starlark.Value = (*MaxErrorsDirective)(nil)
var _ starlark.HasAttrs = (*MaxErrorsDirective)(nil)

func (m *MaxErrorsDirective) Attr(name string) (starlark.Value, error) {
	switch name {
	case "num":
		return starlark.MakeInt(m.Num), nil
	default:
		return nil, starlark.NoSuchAttrError(fmt.Sprintf("max errors directive has no attribute %q", name))
	}
}

func (m *MaxErrorsDirective) AttrNames() []string {
	return []string{"num"}
}

type MaxErrorsDirective struct {
	Num int
}

func (m *MaxErrorsDirective) String() string {
	return fmt.Sprintf("MaxErrorsDirective(Num: %d)", m.Num)
}

func (m *MaxErrorsDirective) Type() string {
	return "max_errors_directive"
}

func (m *MaxErrorsDirective) Freeze() {
	// No mutable fields, so no action needed
}

func (m *MaxErrorsDirective) Truth() starlark.Bool {
	return starlark.Bool(m.Num > 0)
}

func (m *MaxErrorsDirective) Hash() (uint32, error) {
	h := fnv.New32()
	h.Write([]byte(fmt.Sprintf("%d", m.Num)))
	return h.Sum32(), nil
}

func MakeMaxErrorsDirective(mce *parser.MethodCallExpression) (Directive, error) {
	if args, ok := mce.GetArguments().(*parser.ArgumentListExpression); ok {
		exprs := args.GetExpressions()
		if len(exprs) != 1 {
			return nil, errors.New("invalid max errors directive")
		}
		expr := exprs[0]
		if constantExpr, ok := expr.(*parser.ConstantExpression); ok {
			if intValue, ok := constantExpr.GetValue().(int); ok {
				return &MaxErrorsDirective{Num: intValue}, nil
			}
		}
	}
	return nil, errors.New("invalid max errors directive")
}
